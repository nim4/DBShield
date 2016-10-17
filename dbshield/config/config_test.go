package config_test

import (
	"testing"

	"github.com/nim4/DBShield/dbshield/config"
)

func TestParseConfig(t *testing.T) {
	err := config.ParseConfig("../../conf/dbshield.yml")
	if err != nil {
		t.Error("Got error", err)
	}

	err = config.ParseConfig("../../conf/XYZ.yml")
	if err == nil {
		t.Error("Expected error", err)
	}
}
