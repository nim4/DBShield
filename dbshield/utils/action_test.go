package utils_test

import (
	"testing"

	"github.com/nim4/DBShield/dbshield/utils"
)

func TestActionDrop(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		}
	}()
	ret := utils.ActionDrop()
	if ret != nil {
		t.Error("Expected nil got ", ret)
	}
}
