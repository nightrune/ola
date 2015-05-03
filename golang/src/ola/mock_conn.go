/**
 * Mock net.Conn that we can use to in depedency injection and testing
 */
package ola

import (
	"net"
	"time"
)

type MockConn struct {
	// For Read
	read_data     []byte
	read_error    error
	read_deadline time.Time
	read_time     time.Time

	// For Write
	write_data     []byte
	write_error    error
	write_deadline time.Time
	write_time     time.Time

	local_addr  net.Addr
	remote_addr net.Addr
}

func (m MockConn) Read(b []byte) (n int, err error) {
	b = m.read_data
	return len(m.read_data), m.read_error
}

func (m MockConn) SetRead(b []byte, err error) {
	m.read_data = b
	m.read_error = err
}

func (m MockConn) Write(b []byte) (n int, err error) {
	return len(b), nil
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
