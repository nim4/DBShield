package config

import "testing"

func TestParseConfig(t *testing.T) {

	err := ParseConfig("../../conf/dbshield.yml")
	if err != nil {
		t.Error("Got error", err)
	}

}

func TestConfigProtect(t *testing.T) {
	err := configProtect()
	if err != nil {
		t.Error("Got error", err)
	}
}
