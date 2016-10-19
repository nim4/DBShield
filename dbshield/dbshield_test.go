package dbshield

import (
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/nim4/DBShield/dbshield/config"
)

func TestMain(t *testing.T) {
	os.Chdir("../")
}

func TestCheck(t *testing.T) {
	err := Check("conf/dbshield.yml")
	if err != nil {
		t.Error("Got error", err)
	}

	err = Check("Invalid.yml")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestStart(t *testing.T) {
	err := Start("Invalid.yml")
	if err == nil {
		t.Error("Expected error")
	}

	//It should fail if port is already open
	l, err := net.Listen("tcp", config.Config.ListenIP+":"+strconv.Itoa(int(config.Config.ListenPort)))
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	err = Start("conf/dbshield.yml")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestPostConfig(t *testing.T) {
	Check("conf/dbshield.yml")
	err := postConfig()
	if err != nil {
		t.Error("Got error", err)
	}
}
