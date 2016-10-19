package dbshield

import (
	"os"
	"testing"
)

func TestCheck(t *testing.T) {
	err := Check("../conf/dbshield.yml")
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
}

func TestPostConfig(t *testing.T) {
	os.Chdir("../")
	Check("conf/dbshield.yml")
	err := postConfig()
	if err != nil {
		t.Error("Got error", err)
	}
}
