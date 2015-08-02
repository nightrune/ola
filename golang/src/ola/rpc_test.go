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
		logger.Info("Rpc Channel is closed!")
	} else {
		logger.Debug("Failed to close the channel")
		t.Fail()
	}
}

func TestRpcClose(t *testing.T) {
	logger.Info("TestRpcClose")
	sock := NewMockConn(t, 2)
	rpc_chan := NewRpcChannel(sock)
	defer CheckRpcClose(t, rpc_chan)
	read_err := NewMockNetworkError(true, true,
		errors.New("Socket read timed out"))
	sock.SetDefaultRead(make([]byte, 4), read_err)
	rpc_chan.Run()
	rpc_chan.Close()
	time.Sleep(time.Duration(1) * time.Second)
}

func TestRpcCallMethod(t *testing.T) {
	messages := make(chan *ResponseData, 2)
	sock := NewMockConn(t, 2)
	rpc_chan := NewRpcChannel(sock)
	rpc_chan.Run()
	defer rpc_chan.Close()
	read_err := NewMockNetworkError(true, true,
		errors.New("Socket read timed out"))
	sock.SetDefaultRead(nil, read_err)
	sock.ExpectWrite([]byte{}, 0, nil)
	rpc_chan.CallMethod(NewMethodDescriptor(1, "test", "", ""), make([]byte, 0),
		messages)
	var response *ResponseData
	select {
	case response = <-messages:
	case <-time.After(time.Second):
		logger.Info("Response Timed Out")
		t.Fail()
		return
	}

	if response.err != nil {
		logger.Info("Got error from response")
		logger.Debug(response.err.Error())
		t.Fail()
	}
}
