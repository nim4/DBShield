package config_test

import (
	"testing"

	"github.com/nim4/DBShield/dbshield/config"
)

func TestParseConfig(t *testing.T) {
	err := config.ParseConfig("../../conf/dbshield.yml")
	if err != nil {
		t.Error("Not Expected error", err)
	}
}
