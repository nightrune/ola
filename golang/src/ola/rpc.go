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
	"github.com/golang/protobuf/proto"
	"ola/ola_rpc"
	"sync/atomic"
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

type OutstandingResponse struct {
	ret chan ResponseData
	id  int
}

type RpcChannel struct {
	sock                  net.Conn
	outstanding_responses map[int]OutstandingResponse
	closer                chan bool
	running               bool
	sequence_number       uint32
}

func NewRpcChannel(sock net.Conn) *RpcChannel {
	// Start a go closure to read and send the
	rpc_channel := new(RpcChannel)
	rpc_channel.sock = sock
	rpc_channel.closer = make(chan bool, 1)
	rpc_channel.running = false
	return rpc_channel
}

func (m *RpcChannel) Run() {
	m.running = true
	go m._read_forever()
}

func (m *RpcChannel) PendingRPCs() bool {
	return false
}

func (m *RpcChannel) CallMethod(method *MethodDescriptor,
	request_data []byte, c chan *ResponseData) {
	is_streaming := false
	// Actual
	message := new(ola_rpc.RpcMessage)
	message.Type = new(ola_rpc.Type)
	if method.OutputType() == STREAMING_NO_RESPONSE {
		if c == nil {
			logger.Fatal(fmt.Sprintf(
				"Calling streaming method %s with a channel that is non-nul",
				method.String()))
			return
		}
		is_streaming = true
	}

	if is_streaming {
		*message.Type = ola_rpc.Type_STREAM_REQUEST
	} else {
		*message.Type = ola_rpc.Type_REQUEST
	}

	id := atomic.AddUint32(&m.sequence_number, 1)
	message.Id = &id
	message.Buffer = request_data
	data, err := proto.Marhall(message)
	if err != nil {
		ola.Fatal("Failed to marshal data in rpc...")
		data := ResponseData{data: nil, err: errors.New(
			"Failed to marshal the message")}
		c <- &data
		return
	}

	err = m._SendMsg(data)
	if err != nil {
		ola.Warn("Failed to send message...")
		data := ResponseData{data: nil, err: errors.New("Failed to send message")}
		c <- &data
		return
	}
}

func (m *RpcChannel) _SendMsg(data []byte) error {

	return nil
}

func (m *RpcChannel) Close() {
	m.running = false
	m.closer <- true
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
			logger.Info("Read timed out..")
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

func (m *RpcChannel) _read_forever() {
	var header []byte
	var buf []byte
	var msg_bytes_read uint32

	defer m.sock.Close()

	header = make([]byte, 4)
	for {
		bytes_read, err := m.read(header)
		if err != nil {
			if checkConnError(err) {
				m.Close()
				return
			}
		} else if bytes_read != len(header) {
			logger.Fatal("Couldn't get header size..")
			m.Close()
			return
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

			// Handle new message
		}

		select {
		case <-m.closer:
			logger.Info("Close recieved on channel.. closing channel")
			return
		}
	}
}
