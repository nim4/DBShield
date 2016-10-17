package dbms_test

import (
	"testing"

	"github.com/nim4/DBShield/dbshield/dbms"
)

func TestReadPacket(t *testing.T) {
	var s mockConn
	buf, err := dbms.ReadPacket(s)
	if err != nil {
		t.Fatal(err)
	}
	if len(buf) != 0 {
		t.Errorf("Expected empty buff, got %s", buf)
	}
}
