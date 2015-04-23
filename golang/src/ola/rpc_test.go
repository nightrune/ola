package ola

import (
	"fmt"
	"net"
	"testing"
	"time"
)

type MockConn struct {
}

func (m MockConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (m MockConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (m MockConn) Close() error {
	return nil
}

func (m MockConn) LocalAddr() net.Addr {
	return nil
}

func (m MockConn) RemoteAddr() net.Addr {
	return nil
}

func (m MockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m MockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestRpcCallMethod(t *testing.T) {
	messages := make(chan *ResponseData, 2)

	sock := MockConn{}
	rpc_chan := NewRpcChannel(sock)
	rpc_chan.CallMethod(&MethodDescriptor{_index: 1, _name: "test"},
		make([]byte, 0), messages)
	fmt.Printf("%s\n", <-messages)
}
