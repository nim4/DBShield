package dbms

import (
	"bytes"
	"errors"
	"testing"

	"github.com/nim4/mock"
)

func TestReadWrite(t *testing.T) {
	var s = mock.ConnMock{Error: errors.New("Dummy Error")}
	err := readWrite(s, s, ReadPacket)
	if err == nil {
		t.Error("Expected error")
	}

	s.Error = nil
	err = readWrite(s, s, ReadPacket)
	if err != nil {
		t.Error("Got error", err)
	}
}

func BenchmarkReadPacket(b *testing.B) {
	var buf [1024]byte
	for i := 0; i < b.N; i++ {
		s := bytes.NewReader(buf[:])
		ReadPacket(s)
	}
}
