package dbms_test

import (
	"net"
	"time"
)

type mockConn struct {
}

func (m mockConn) Read(b []byte) (n int, err error) {
	return
}

func (m mockConn) Write(b []byte) (n int, err error) {
	return
}

func (m mockConn) Close() error {
	return nil
}

func (m mockConn) LocalAddr() net.Addr {
	addr, _ := net.ResolveTCPAddr("tcp", ":80")
	return addr
}

func (m mockConn) RemoteAddr() net.Addr {
	return m.LocalAddr()
}

func (m mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}
