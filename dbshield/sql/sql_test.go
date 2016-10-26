package sql_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/nim4/DBShield/dbshield/sql"
)

func TestQueryContext(t *testing.T) {
	c := sql.QueryContext{
		Query:    []byte("select * from test;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	r := c
	b := c.Marshal()
	c.Unmarshal(b)

	if bytes.Compare(c.Query, r.Query) != 0 {
		t.Error("Expected", r.Query, "got", c.Query)
	}

	if bytes.Compare(c.Query, r.Query) != 0 {
		t.Error("Expected", r.Database, "got", c.Database)
	}

	if bytes.Compare(c.User, r.User) != 0 {
		t.Error("Expected", r.User, "got", c.User)
	}

	if bytes.Compare(c.Client, r.Client) != 0 {
		t.Error("Expected", r.Client, "got", c.Client)
	}

	if c.Time != r.Time {
		t.Error("Expected", r.Time, "got", c.Time)
	}
}

func TestPattern(t *testing.T) {
	p := sql.Pattern([]byte("select * from X;"))
	if len(p) < 4 {
		t.Error("Unexpected Pattern")
	}
}
