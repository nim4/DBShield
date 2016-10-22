package sql_test

import (
	"testing"

	"github.com/nim4/DBShield/dbshield/sql"
)

func TestPattern(t *testing.T) {
	p := sql.Pattern("select * from X;")
	if len(p) < 4 {
		t.Error("Unexpected Pattern")
	}
}
