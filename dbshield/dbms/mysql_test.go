package dbms_test

import (
	"bytes"
	"testing"

	"github.com/nim4/DBShield/dbshield/dbms"
)

func TestMySQLGetUsernameDB(t *testing.T) {
	// GetUsernameDB gets buf[4:] (remove leading 4 bytes from packet dump)

	u, d := dbms.MySQLGetUsernameDB([]byte{
		5, 162, 43, 0, 1, 0, 0, 0, 45, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 119, 112, 0, 20, 220, 242, 53, 124, 75, 14, 62, 51,
		39, 210, 196, 162, 213, 205, 209, 232, 229, 71, 70, 62, 109, 121, 115, 113,
		108, 95, 110, 97, 116, 105, 118, 101, 95, 112, 97, 115, 115, 119, 111, 114,
		100, 0,
	})
	if string(u) != "wp" {
		t.Error("Expected 'wp' username got", string(u))
	}
	if len(d) != 0 {
		t.Error("Expected empty db name got", string(d))
	}

	u, d = dbms.MySQLGetUsernameDB([]byte{
		13, 162, 43, 0, 1, 0, 0, 0, 45, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 119, 112, 0, 20, 0, 205, 124, 7, 125, 243, 111, 168,
		162, 10, 54, 115, 147, 159, 57, 126, 109, 123, 162, 138, 119, 112, 0, 109,
		121, 115, 113, 108, 95, 110, 97, 116, 105, 118, 101, 95, 112, 97, 115, 115,
		119, 111, 114, 100, 0,
	})

	if string(u) != "wp" {
		t.Error("Expected 'wp' username got", string(u))
	}
	if string(u) != "wp" {
		t.Error("Expected 'wp' db name", string(d))
	}

	u, d = dbms.MySQLGetUsernameDB([]byte{
		13, 162, 43, 0, 1, 0, 0, 0, 45, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0,
	})

	if len(u) != 0 || len(d) != 0 {
		t.Error("Expected empty username & db name got", string(u), string(d))
	}
}

func TestMySQLReadPacket(t *testing.T) {
	const maxPayloadLen = 1<<24 - 1
	var buf [maxPayloadLen]byte
	reader := bytes.NewReader(buf[:])
	b, err := dbms.MySQLReadPacket(reader)
	if err != nil {
		t.Error("Got error", err)
	}
	if bytes.Compare(b, buf[:]) != 0 {
		t.Error("Unexpected output")
	}
	var eofReader bytes.Buffer
	_, err = dbms.MySQLReadPacket(&eofReader)
	if err != nil {
		t.Error("Got error", err)
	}
}

func BenchmarkMySQLReadPacket(b *testing.B) {
	var buf [1024]byte
	for i := 0; i < b.N; i++ {
		s := bytes.NewReader(buf[:])
		dbms.MySQLReadPacket(s)
	}
}
