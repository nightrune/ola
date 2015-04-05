/**
 * This implements the rpc channel for go.
 *
 */
package ola

import "net"
import "ola/ola_rpc"

var _ = ola_rpc.Type_name // Here while we build the rpc channel, delete later

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
}

func NewRpcChannel(sock net.Conn) *RpcChannel {
	// Start a go closure to read and send the
	rpc_channel := new(RpcChannel)
	rpc_channel.sock = sock
	return rpc_channel
}

func (m *RpcChannel) PendingRPCs() bool {
	return false
}

func (m *RpcChannel) CallMethod(method *MethodDescriptor,
	request_data []byte, c chan *ResponseData) {
	data := &ResponseData{data: nil, err: NewNotImplemented(
		"This is currently not implemented")}
	c <- data
}

func (m *RpcChannel) Close() {

}

func (m *RpcChannel) _read_forever() {
	//
}
