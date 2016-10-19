package dbms_test

import (
	"testing"

	"github.com/nim4/DBShield/dbshield/dbms"
	"github.com/nim4/mock"
)

func TestReadPacket(t *testing.T) {
	var s mock.ConnMock
	buf, err := dbms.ReadPacket(s)
	if err != nil {
		t.Fatal(err)
	}
	if len(buf) != 0 {
		t.Errorf("Expected empty buff, got %s", buf)
	}

	mock.ReturnError(true)
	defer mock.ReturnError(false)

	_, err = dbms.ReadPacket(s)
	if err == nil {
		t.Errorf("Expected error")
	}
}
