package ola

import (
	"testing"
)

func TestRpcCallMethod(t *testing.T) {
	messages := make(chan *ResponseData, 2)

	sock := MockConn{}
	rpc_chan := NewRpcChannel(sock)
	rpc_chan.Run()
	defer rpc_chan.Close()
	rpc_chan.CallMethod(NewMethodDescriptor(1, "test"), make([]byte, 0), messages)
	<-messages
}
