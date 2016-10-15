package dbshield

import "testing"

func TestDbNameToStruct(t *testing.T) {
	_, err := dbNameToStruct("mysql")
	if err != nil {
		t.Error("Expected struct, got ", err)
		return
	}
	_, err = dbNameToStruct("invalid")
	if err == nil {
		t.Error("Expected error")
		return
	}
}
