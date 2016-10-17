package dbms_test

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"testing"

	"github.com/nim4/DBShield/dbshield/dbms"
)

func TestReadPacket(t *testing.T) {
	addr := ":" + strconv.Itoa(45000+rand.Intn(10000))
	message := "Message Body"

	go func() {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		if _, err := fmt.Fprintf(conn, message); err != nil {
			t.Fatal(err)
		}
	}()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		return
	}
	defer conn.Close()

	buf, err := dbms.ReadPacket(conn)
	if err != nil {
		t.Fatal(err)
	}
	if res := string(buf); res != message {
		t.Errorf("Expected %s, got %s", message, res)
	}
	buf, err = dbms.ReadPacket(conn)
	if err == nil {
		t.Errorf("Expected error")
	}
}
