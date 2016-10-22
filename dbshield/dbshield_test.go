package dbshield

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/training"
)

func TestMain(t *testing.T) {
	os.Chdir("../")
}

func TestSetConfigFile(t *testing.T) {
	err := SetConfigFile("Invalid.yml")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestPatterns(t *testing.T) {
	if training.DBCon != nil {
		training.DBCon.Close()
	}
	err := SetConfigFile("conf/dbshield.yml")
	if err != nil {
		t.Error("Got error", err)
	}
	err = Patterns()
	if err != nil {
		t.Error("Got error", err)
	}
}

func TestCheck(t *testing.T) {
	SetConfigFile("conf/dbshield.yml")
	err := Check()
	if err != nil {
		t.Error("Got error", err)
	}
}

func TestPostConfig(t *testing.T) {
	SetConfigFile("conf/dbshield.yml")
	config.Config.DBType = "Invalid"
	err := postConfig()
	if err == nil {
		t.Error("Expected error")
	}

	config.Config.ListenPort = 0
	config.Config.DBType = "mysql"
	err = postConfig()
	if err != nil {
		t.Error("Expected nil got ", err)
	}

	config.Config.TLSCertificate = os.TempDir()
	err = postConfig()
	if err == nil {
		t.Error("Expected error")
	}
}

func TestStart(t *testing.T) {
	if training.DBCon != nil {
		training.DBCon.Close()
	}
	SetConfigFile("conf/dbshield.yml")
	//It should fail if port is already open
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Config.ListenIP, config.Config.ListenPort))
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	err = Start()
	if err == nil {
		t.Error("Expected error")
	}

	config.Config.TargetIP = "in valid"
	err = Start()
	if err == nil {
		t.Error("Expected error")
	}
}
