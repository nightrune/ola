/**
 * This implements the rpc channel for go.
 *
 */
package ola

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	//"ola/ola_proto"
	"errors"
	"github.com/golang/protobuf/proto"
	"ola/ola_rpc"
	"time"
)

var _ = ola_rpc.Type_name // Here while we build the rpc channel, delete later

const VERSION_MASK uint32 = 0xFF000000
const SIZE_MASK uint32 = 0x00FFFFFF
const PROTOCOL_VERSION = 1

// todo(Sean) Find a better place to put this since there
// is no static method provided by protobufs to get this
// string....
const STREAMING_NO_RESPONSE = "STREAMIN_NO_RESPONSE"

type NotImplemented struct {
	s string
}

func NewNotImplemented(s string) *NotImplemented {
	ni := new(NotImplemented)
	ni.s = s
	return ni
}

func (m *NotImplemented) Error() string {
	return m.s
}

type ResponseData struct {
	data []byte
	err  error
}

type OutstandingRequest struct {
	ret chan *ResponseData
	req *ola_rpc.RpcMessage
}

type RpcChannel struct {
	sock            net.Conn
	closer          chan bool
	running         bool
	sequence_number uint32
	new_request     chan *OutstandingRequest
	new_response    chan []byte
	write           chan []byte
}

func NewRpcChannel(sock net.Conn) *RpcChannel {
	// Start a go closure to read and send the
	rpc_channel := new(RpcChannel)
	rpc_channel.sock = sock
	rpc_channel.running = false
	return rpc_channel
}

func (m *RpcChannel) Run() {
	if m.sock == nil {
		logger.Warn("Can't Rpc Channel when it has a nil socket")
		return
	}
	m.new_request = make(chan *OutstandingRequest)
	m.new_response = make(chan []byte)
	m.write = make(chan []byte)
	m.closer = make(chan bool, 1)
	m.running = true
	go m._read_forever()
	go m._write_forever()
	go m._demux()
}

func (m *RpcChannel) Close() {
	m.running = false
	close(m.closer)
}

func (m *RpcChannel) PendingRPCs() bool {
	return false
}

func (m *RpcChannel) CallMethod(method *MethodDescriptor,
	request_data []byte, c chan *ResponseData) {
	is_streaming := false
	message := new(ola_rpc.RpcMessage)
	message.Type = new(ola_rpc.Type)
	message.Id = new(uint32)
	if method.OutputType() == STREAMING_NO_RESPONSE {
		if c != nil {
			logger.Warn(fmt.Sprintf(
				"Calling streaming method %s with a channel that is non-nul",
				method.String()))
			return
		}
		is_streaming = true
		c = nil
	}

	if is_streaming {
		*message.Type = ola_rpc.Type_STREAM_REQUEST
	} else {
		*message.Type = ola_rpc.Type_REQUEST
	}

	message.Buffer = request_data
	request := new(OutstandingRequest)
	request.req = message
	request.ret = c
	m.new_request <- request
}

func (m *RpcChannel) IsClosed() bool {
	return !m.running
}

func parseHeader(header []byte) (size uint32, version uint8) {
	binary.Read(bytes.NewReader(header), binary.BigEndian, size)
	version = uint8(size & VERSION_MASK >> 28)
	size = size & SIZE_MASK
	return size, version
}

func checkConnError(err error) bool {
	if err == nil {
		logger.Debug("Got a nil error in check ConnError")
		return false
	}

	nerr, ok := err.(net.Error)
	if ok == true {
		if nerr.Timeout() {
			logger.Debug("Read timed out..")
			return false
		} else {
			logger.Info("Connection was closed during read..\n")
			return true
		}
	}
	logger.Fatal(
		"Unknown error checking connection read error.. Stopping reading")
	return true
}

func (m *RpcChannel) read(b []byte) (n int, err error) {
	m.sock.SetReadDeadline(time.Now().Add(1 * time.Second))
	return m.sock.Read(b)
}

func (m *RpcChannel) _handle_request_timeout(request *OutstandingRequest) {
	logger.Warn("Request timed out")
	response := ResponseData{data: nil, err: errors.New("Request Timedout")}
	request.ret <- &response
}

func (m *RpcChannel) _demux() {
	outstanding_requests := make(map[uint32]*OutstandingRequest)
	var sequence_number uint32 = 0
	select {
	case <-m.new_response:
		// Deal with incoming messages
	case request := <-m.new_request:
		// Deal with new outgoing messagses
		var id uint32
		id = sequence_number
		*request.req.Id = id
		sequence_number = sequence_number + 1
		data, err := proto.Marshal(request.req)
		if err != nil {
			logger.Fatal("Failed to marshal data in rpc...")
			data := ResponseData{data: nil, err: errors.New(
				"Failed to marshal the message")}
			request.ret <- &data
			return
		}

		v, ok := outstanding_requests[id]
		if ok {
			m._handle_request_timeout(v)
			delete(outstanding_requests, id)
		}
		outstanding_requests[id] = request
		m.write <- data
	case <-time.After(10 * time.Millisecond):
		// Timeout code, check on time outs
	case <-m.closer:
		logger.Debug("Timing out all other requests")
		// Close out all other requests
		for _, v := range outstanding_requests {
			m._handle_request_timeout(v)
		}
		// force garbage collection
		outstanding_requests = nil
		return
	}
}

func (m *RpcChannel) _write_forever() {
	var current int
	for {
		select {
		case data := <-m.write:
			current = 0
			for current < len(data) {
				n, err := m.sock.Write(data[current:])
				if err != nil {
					logger.Fatal("Failed to write to rpc socket..")
					m.Close()
					return
				}
				current += n
			}
		case <-m.closer:
			logger.Debug("Stopping write on socket")
			return
		}
	}
}

func (m *RpcChannel) _read_forever() {
	var header []byte
	var buf []byte
	var msg_bytes_read uint32

	defer m.sock.Close()

	header = make([]byte, 4)
	logger.Debug("Starting read forever")
	for {
		bytes_read, err := m.read(header)
		if err != nil {
			if checkConnError(err) {
				m.Close()
			}
		} else if bytes_read != len(header) {
			logger.Fatal("Couldn't read enough bytes to get rpc header..")
			m.Close()
		} else {
			expected_size, protocol_version := parseHeader(header)
			buf = make([]byte, expected_size, expected_size)

			if protocol_version != PROTOCOL_VERSION {
				logger.Warn(fmt.Sprintf("protocol mismatch %d != %d", protocol_version,
					PROTOCOL_VERSION))
			}

			for msg_bytes_read < expected_size {
				bytes_read, err = m.read(buf[msg_bytes_read:])
				if err != nil {
					if checkConnError(err) {
						m.Close()
						return
					}
				} else {
					msg_bytes_read = msg_bytes_read + uint32(bytes_read)
				}
			}
			m.new_response <- buf
		}

		select {
		case <-m.closer:
			logger.Info("Close recieved on channel.. closing channel")
			return
		}
	}
}
