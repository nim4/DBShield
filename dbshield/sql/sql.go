package sql

import (
	"strconv"
	"time"

	"github.com/xwb1989/sqlparser"
)

//QueryContext holds information around query
type QueryContext struct {
	Query  string
	User   string
	Client string
	Time   time.Time
}

//Pattern returns pattern of given query
func Pattern(query string) (pattern []byte) {
	tokenizer := sqlparser.NewStringTokenizer(query)
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
