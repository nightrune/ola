/**
 * This implements the rpc channel for go.
 *
 */
package ola

import "net"

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

func (self *RpcChannel) CallMethod() chan ResponseData {
	return nil
}

func (self *RpcChannel) Close() {

}

func (self *RpcChannel) _read_forever() {
	//
}
