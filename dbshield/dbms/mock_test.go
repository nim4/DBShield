package dbms_test

import (
	"errors"
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

type mockConnError struct {
}

func (m mockConnError) Read(b []byte) (n int, err error) {
	err = errors.New("Dummy Error")
	return
}

func (m mockConnError) Write(b []byte) (n int, err error) {
	err = errors.New("Dummy Error")
	return
}

func (m mockConnError) Close() error {
	return errors.New("Dummy Error")
}

func (m mockConnError) LocalAddr() net.Addr {
	addr, _ := net.ResolveTCPAddr("tcp", ":80")
	return addr
}

func (m mockConnError) RemoteAddr() net.Addr {
	return m.LocalAddr()
}

func (m mockConnError) SetDeadline(t time.Time) error {
	return errors.New("Dummy Error")
}

func (m mockConnError) SetReadDeadline(t time.Time) error {
	return errors.New("Dummy Error")
}

func (m mockConnError) SetWriteDeadline(t time.Time) error {
	return errors.New("Dummy Error")
}
