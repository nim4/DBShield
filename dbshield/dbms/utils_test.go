package dbms

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/sql"
	"github.com/nim4/DBShield/dbshield/training"
	"github.com/nim4/mock"
)

func TestMain(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "testdb")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()
	path := tmpfile.Name()
	training.DBConLearning, err = bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}
	training.DBConProtect, err = bolt.Open(path+"2", 0600, nil)
	if err != nil {
		panic(err)
	}
}

func TestEbc2asc(t *testing.T) {
	ret := string(ebc2asc([]byte{0xe2, 0xc1}))
	if ret != "SA" {
		t.Error("Expected 'SA', got ", ret)
	}
}

func TestPascalString(t *testing.T) {
	b, size := pascalString([]byte{0x3, 0x41, 0x41, 0x41})
	if size != 3 {
		t.Error("Expected 3, got ", size)
	}
	if bytes.Compare(b, []byte{0x41, 0x41, 0x41}) != 0 {
		t.Error("Expected 'AAA', got ", b)
	}
}

func TestRemoteAddrToIP(t *testing.T) {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:80")
	ip := remoteAddrToIP(addr)
	if bytes.Compare(ip, []byte("127.0.0.1")) != 0 {
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
	panic("Test Panic")
}

func TestProcessContext(t *testing.T) {
	c := sql.QueryContext{
		Query:    []byte("select * from test;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127,0,0,1"),
		Time:     time.Now(),
	}
	config.Config.Learning = true
	processContext(c)
	config.Config.Learning = false
	config.Config.ActionFunc = func() error { return nil }
	processContext(c)
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
