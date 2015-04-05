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
	rpc_channel.closer = make(chan bool)
	rpc_channel.outstanding_responses = make(map[int]OutstandingResponse)
	return rpc_channel
}

func (self *RpcChannel) PendingRPCs() bool {
	return false
}

func (self *RpcChannel) CallMethod(method *MethodDescriptor,
	request_data []byte) chan ResponseData {
	c := make(chan ResponseData)
	c <- ResponseData{data: nil, err: NewNotImplemented(
		"This is currently not implemented")}
	return c
}

func (self *RpcChannel) Close() {

}

func (self *RpcChannel) _read_forever() {
	//
}
