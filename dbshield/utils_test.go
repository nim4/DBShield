// +build !windows

package dbshield

import (
	"net"
	"os"
	"testing"

	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/mock"
)

func TestDbNameToStruct(t *testing.T) {
	_, err := dbNameToStruct("db2")
	if err != nil {
		t.Error("Expected struct, got ", err)
		return
	}
	_, err = dbNameToStruct("mysql")
	if err != nil {
		t.Error("Expected struct, got ", err)
		return
	}
	_, err = dbNameToStruct("oracle")
	if err != nil {
		t.Error("Expected struct, got ", err)
		return
	}
	_, err = dbNameToStruct("postgres")
	if err != nil {
		t.Error("Expected struct, got ", err)
		return
	}
	//Invalid case is tested in postConfig test
}

func TestInitLogging(t *testing.T) {
	config.Config.LogPath = "stdout"
	initLogging()
	config.Config.LogPath = "stderr"
	initLogging()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic")
		}
	}()
	config.Config.LogPath = os.TempDir()
	initLogging()
}

func TestHandleClient(t *testing.T) {
	var s mock.ConnMock
	err := handleClient(s, nil)
	if err == nil {
		t.Error("Expected error got nil")
	}
	ls, _ := net.Listen("tcp4", "localhost:0")
	go func() {
		for {
			conn, _ := ls.Accept()
			conn.Close()
		}
	}()

	ra, err := net.ResolveTCPAddr("tcp4", ls.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	err = handleClient(s, ra)
	if err == nil {
		t.Error("Expected error got nil")
	}
}

func TestCloseHandlers(t *testing.T) {
	logger.Output = os.Stderr
	defer func() {
		if r := recover(); r != nil {
			t.Error("Panic!")
		}
	}()
	closeHandlers()
}
