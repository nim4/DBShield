package utils_test

import (
	"net"
	"testing"

	"github.com/nim4/DBShield/dbshield/utils"
)

func TestActionDrop(t *testing.T) {
	s := new(net.Conn)
	ret := utils.ActionDrop(*s)
	if ret != nil {
		t.Error("Expected nil got ", ret)
	}
}
