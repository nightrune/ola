package ola

import (
	"errors"
	"net"
	"testing"
	"time"
)

func Check(err error) {
	nerr, ok := err.(net.Error)
	if ok == true {
		if nerr.Timeout() {
			logger.Info("This worked?")
		}
	}
}

func CheckRpcClose(t *testing.T, chn *RpcChannel) {
	if chn.IsClosed() {
		logger.Info("Closing Rpc Channel")
	} else {
		logger.Debug("Failed to close the channel")
		t.Fail()
	}
}

func TestRpcClose(t *testing.T) {
	sock := new(MockConn)
	rpc_chan := NewRpcChannel(sock)
	defer CheckRpcClose(t, rpc_chan)
	read_err := NewMockNetworkError(true, true,
		errors.New("Socket read timed out"))
	sock.SetRead(make([]byte, 4), read_err)
	rpc_chan.Run()
	rpc_chan.Close()
	time.Sleep(time.Duration(1) * time.Second)
}

func TestRpcCallMethod(t *testing.T) {
	messages := make(chan *ResponseData, 2)
	sock := new(MockConn)
	rpc_chan := NewRpcChannel(sock)
	rpc_chan.Run()
	defer rpc_chan.Close()
	read_err := new(MockNetworkError)
	read_err.timeout = true
	read_err.temp = true
	read_err.err = errors.New("Socket read timed out")
	sock.SetRead(make([]byte, 4), read_err)
	sock.ExpectWrite([]byte{}, 0, nil)
	rpc_chan.CallMethod(NewMethodDescriptor(1, "test"), make([]byte, 0), messages)
	response := <-messages
	if response.err != nil {
		logger.Info("Got error from response")
		t.Fail()
	}
}
