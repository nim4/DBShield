package dbms

import (
	"crypto/tls"
	"net"
	"testing"

	"github.com/nim4/mock"
)

func TestPascalString(t *testing.T) {
	str, size := pascalString([]byte{0x3, 0x41, 0x41, 0x41})
	if size != 3 {
		t.Error("Expected 3, got ", size)
	}
	if str != "AAA" {
		t.Error("Expected 'AAA', got ", str)
	}
}

func TestRemoteAddrToIP(t *testing.T) {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:80")
	ip := remoteAddrToIP(addr)
	if ip != "127.0.0.1" {
		t.Error("Expected '127.0.0.1', got ", ip)
	}
}

func TestThreeByteBigEndianToInt(t *testing.T) {
	res := threeByteBigEndianToInt([]byte{1, 2, 3})
	if res != 197121 {
		t.Error("Expected 197121, got ", res)
	}
}

func TestHandlePanic(t *testing.T) {
	defer handlePanic()
	panic("")
}

func TestTurnSSL(t *testing.T) {
	cert, err := tls.LoadX509KeyPair("../../cert/server-cert.pem", "../../cert/server-key.pem")
	if err != nil {
		t.Fatal(err)
	}

	mock.ReturnError(true)
	defer mock.ReturnError(false)
	var s mock.ConnMock
	_, _, err = turnSSL(s, s, cert)
	if err == nil {
		t.Error("Expected error")
	}
}
