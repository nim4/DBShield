package training

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync/atomic"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
)

var (
	//QueryCounter state
	QueryCounter = uint64(0)

	//AbnormalCounter state
	AbnormalCounter = uint64(0)

	//DBCon holds a pointer to our local database connection
	DBCon *bolt.DB

	errInvalidParrent = errors.New("Invalid pattern")
	errInvalidUser    = errors.New("Invalid user")
	errInvalidClient  = errors.New("Invalid client")
)

//AddToTrainingSet records query context in local database
func AddToTrainingSet(context sql.QueryContext) error {
	atomic.AddUint64(&QueryCounter, 1)
	pattern := sql.Pattern(context.Query)
	//logger.Debug("Pattern", pattern)

	if err := DBCon.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pattern"))
		if b == nil {
			return errors.New("Invalid DB")
		}
		if b.Get(pattern) == nil {
			b.Put(pattern, context.Query)
		}

		uKey := bytes.Buffer{}
		uKey.Write(pattern)
		uKey.WriteString("_user_")
		uKey.Write(context.User)
		b.Put(uKey.Bytes(), []byte{0x11})

		cKey := bytes.Buffer{}
		cKey.Write(pattern)
		cKey.WriteString("_client_")
		cKey.Write(context.Client)
		b.Put(cKey.Bytes(), []byte{0x11})
		return nil
	}); err != nil {
		logger.Warning(err)
		return err
	}
	return nil
}

//CheckQuery pattern, returns true if it finds the pattern
//We should keep it as fast as possible
func CheckQuery(context sql.QueryContext) bool {
	atomic.AddUint64(&QueryCounter, 1)
	pattern := sql.Pattern(context.Query)
	if err := DBCon.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pattern"))
		if b == nil {
			panic("Invalid DB")
		}
		if b.Get(pattern) == nil {
			return errInvalidParrent
		}

		key := bytes.Buffer{}
		if config.Config.CheckUser {
			key.Write(pattern)
			key.WriteString("_user_")
			key.Write(context.User)
			if b.Get(key.Bytes()) == nil {
				return errInvalidUser
			}
		}
		if config.Config.CheckSource {
			key.Reset()
			key.Write(pattern)
			key.WriteString("_client_")
			key.Write(context.Client)
			if b.Get(key.Bytes()) == nil {
				return errInvalidClient
			}
		}
		return nil
	}); err != nil {
		logger.Warning(err)
		//Record abnormal
		recordAbnormal(pattern, context)
		return false
	}
	return true
}

func recordAbnormal(pattern []byte, context sql.QueryContext) error {
	atomic.AddUint64(&AbnormalCounter, 1)
	return DBCon.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("abnormal"))
		if b == nil {
			panic("Invalid DB")
		}
		id, _ := b.NextSequence()
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(id))
		return b.Put(buf, context.Marshal())
	})
}
