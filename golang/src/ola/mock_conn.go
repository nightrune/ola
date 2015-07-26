/**
 * Mock net.Conn that we can use to in depedency injection and testing
 */
package ola

import (
	"net"
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

type MockConn struct {
	// For Read
	read_data     []byte
	read_error    error
	read_deadline time.Time
	read_time     time.Time

	// For Write
	write_data     []byte
	write_n        int
	write_error    error
	write_deadline time.Time
	write_time     time.Time

	local_addr  net.Addr
	remote_addr net.Addr
}

func (m MockConn) Read(b []byte) (n int, err error) {
	copied := copy(b, m.read_data)
	return copied, m.read_error
}

func (m *MockConn) SetRead(b []byte, err error) {
	m.read_data = b
	m.read_error = err
}

func (m MockConn) Write(b []byte) (n int, err error) {
	return m.write_n, m.write_error
}

func (m *MockConn) ExpectWrite(b []byte, return_n int, return_err error) {
	m.write_data = b
	m.write_n = return_n
	m.write_error = return_err
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
