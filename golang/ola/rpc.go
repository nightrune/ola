/**
 * This implements the rpc channel for go. 
 *
 */
package ola.rpc

import "net"

type ResponseData struct {
	data []byte
	err error
}

type OutstandingResponse struct {
	ret chan ResponseData
	id int
}

type Channel struct {
	net.Conn sock
	outstanding_responses map[int]OutstandingResponse
	closer chan bool
}

func NewChannel(sock net.Conn) *Channel {
	// Start a go closure to read and send the
  rpc_channel := new(Channel)
	rpc_channel.closer = make(chan bool)
  rpc_channel.outstanding_responses = make(map[int]OutstandingResponse)
	return &rpc_channel
}

func (self *Channel) PendingRPCs() bool {
	return false
}

func (self *Channel) CallMethod() chan ResponseData {
	return nil
}

func (self *Channel) Close() {

}

func (self *Channel) _read_forever() {
  for self.
}
