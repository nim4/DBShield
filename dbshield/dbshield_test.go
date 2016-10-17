package dbshield_test

import (
	"testing"

	"github.com/nim4/DBShield/dbshield"
)

func TestCheck(t *testing.T) {
	err := dbshield.Check("../conf/dbshield.yml")
	if err != nil {
		t.Error("Got error", err)
	}

	err = dbshield.Check("../conf/XYZ.yml")
	if err == nil {
		t.Error("Expected error")
	}
}
