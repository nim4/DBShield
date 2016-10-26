package sql

import (
	"bytes"
	"strconv"
	"time"

	"github.com/xwb1989/sqlparser"
)

//QueryContext holds information around query
type QueryContext struct {
	Query    []byte
	User     []byte
	Client   []byte
	Database []byte
	Time     time.Time
}

//Unmarshal []byte into QueryContext
func (c *QueryContext) Unmarshal(b []byte) (n int) {
	n = bytes.IndexByte(b, 0x00)
	c.Query = b[:n]
	n++

	i := bytes.IndexByte(b[n:], 0x00)
	c.User = b[n : n+i]
	n += i + 1

	i = bytes.IndexByte(b[n:], 0x00)
	c.Client = b[n : n+i]
	n += i + 1

	i = bytes.IndexByte(b[n:], 0x00)
	c.Database = b[n : n+i]
	n += i

	n += bytes.IndexByte(b[n:], 0x00)
	c.Time.UnmarshalBinary(b[n:])
	n += 8
	return
}

//Marshal load []byte into QueryContext
func (c *QueryContext) Marshal() (b []byte) {
	t, _ := c.Time.MarshalBinary()
	b = append(b, c.Query...)
	b = append(b, 0x00)
	b = append(b, c.User...)
	b = append(b, 0x00)
	b = append(b, c.Client...)
	b = append(b, 0x00)
	b = append(b, c.Database...)
	b = append(b, 0x00)
	b = append(b, t...)
	return
}

//Pattern returns pattern of given query
func Pattern(query []byte) (pattern []byte) {
	tokenizer := sqlparser.NewStringTokenizer(string(query))
	for {
		typ, val := tokenizer.Scan()
		switch typ {
		case sqlparser.ID: //table, database, variable & ... names
			pattern = append(pattern, val...)
		case 0: //End of query
			return
		default:
			//because its 4x faster than "enconding/binary" (but 10x uglier)
			pattern = append(pattern, []byte(strconv.Itoa(typ))...)
		}
	}
}
