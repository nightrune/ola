/**
 * Mock net.Conn that we can use to in depedency injection and testing
 */
package ola

import (
	"errors"
	"net"
	"testing"
	"time"
)

type MockNetworkError struct {
	timeout bool
	temp    bool
	err     error
}

func (m *MockNetworkError) Error() string {
	return m.err.Error()
}

func (m *MockNetworkError) Timeout() bool {
	return m.timeout
}

func (m *MockNetworkError) Temporary() bool {
	return m.temp
}

func NewMockNetworkError(timeout bool, temp bool, err error) *MockNetworkError {
	mnerr := new(MockNetworkError)
	mnerr.timeout = timeout
	mnerr.temp = temp
	mnerr.err = err
	return mnerr
}

type ReadData struct {
	data []byte
	err  error
}

type WriteData struct {
	data []byte
	n    int
	err  error
}

type MockConn struct {
	// For Read
	default_read  *ReadData
	read_data     []byte
	read_deadline time.Time
	read_chan     chan *ReadData

	// For Write
	write_deadline time.Time
	write_chan     chan *WriteData
	default_write  *WriteData

	local_addr  net.Addr
	remote_addr net.Addr
	test        *testing.T
}

func NewMockConn(t *testing.T, queue_depth int) *MockConn {
	mock := new(MockConn)
	mock.test = t
	read_err := NewMockNetworkError(true, true, errors.New("Read timed out"))
	mock.default_read = &ReadData{data: nil, err: read_err}
	write_err := NewMockNetworkError(true, true, errors.New("Write timed out"))
	mock.default_write = &WriteData{data: nil, n: 0, err: write_err}
	mock.read_chan = make(chan *ReadData, queue_depth)
	mock.write_chan = make(chan *WriteData, queue_depth)
	return mock
}

func (m MockConn) Read(b []byte) (n int, err error) {
	var data *ReadData
	select {
	case <-time.After(time.Duration(m.read_deadline.Sub(time.Now()))):
		data = m.default_read
	case data = <-m.read_chan:
	}

	if data.data == nil {
		return 0, data.err
	}
	copied := copy(b, data.data)
	return copied, data.err
}

/**
 * @notes This function is not thread safe
 */
func (m *MockConn) SetDefaultRead(b []byte, err error) {
	data := new(ReadData)
	data.data = b
	data.err = err
	m.default_read = data
}

func (m *MockConn) SetRead(b []byte, err error) {
	data := new(ReadData)
	data.data = b
	data.err = err
	m.read_chan <- data
}

func (m MockConn) Write(b []byte) (n int, err error) {
	var data *WriteData
	select {
	case data = <-m.write_chan:
		// Verify Data against read data
	case <-time.After(m.write_deadline.Sub(time.Now())):
		data = m.default_write
	}
	return data.n, data.err
}

func (m *MockConn) ExpectWrite(b []byte, return_n int, return_err error) {
	data := new(WriteData)
	data.data = b
	data.n = return_n
	data.err = return_err
	m.write_chan <- data
}

func (m MockConn) Close() error {
	return nil
}

func (m MockConn) LocalAddr() net.Addr {
	return m.local_addr
}

func (m MockConn) RemoteAddr() net.Addr {
	return m.remote_addr
}

func (m MockConn) SetDeadline(t time.Time) error {
	m.read_deadline = t
	m.write_deadline = t
	return nil
}

func (m MockConn) SetReadDeadline(t time.Time) error {
	m.read_deadline = t
	return nil
}

func (m MockConn) SetWriteDeadline(t time.Time) error {
	m.write_deadline = t
	return nil
}
