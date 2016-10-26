package training

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
)

//DBConLearning holds a pointer to our local database connection which has valid requests
var DBConLearning *bolt.DB

//DBConProtect holds a pointer to our local database connection which has rejected requests
var DBConProtect *bolt.DB

var dummyErrorToExitIter = errors.New("Dummy Error")

//AddToTrainingSet records query context in local database
func AddToTrainingSet(context sql.QueryContext) error {

	pattern := sql.Pattern(context.Query)
	//logger.Debug("Pattern", pattern)

	if err := DBConLearning.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(pattern)
		if err != nil {
			return err
		}
		encoded := context.Marshal()
		if err = b.Put(uniqID(b), encoded); err != nil {
			return err
		}
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
	pattern := sql.Pattern(context.Query)
	if err := DBConLearning.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pattern))
		if b == nil {
			return errors.New("Pattern not found")
		}
		if !config.Config.CheckUser && !config.Config.CheckSource {
			return nil
		}
		var tmpContext sql.QueryContext
		validUser := !config.Config.CheckUser
		validSource := !config.Config.CheckSource
		err := b.ForEach(func(k []byte, v []byte) error {
			if validUser && validSource {
				return dummyErrorToExitIter
			}
			tmpContext.Unmarshal(v)
			logger.Info(string(tmpContext.User), string(context.User))
			if !validUser && bytes.Compare(tmpContext.User, context.User) == 0 {
				validUser = true
			}
			logger.Info(string(tmpContext.Client), string(context.Client))
			if !validSource && bytes.Compare(tmpContext.Client, context.Client) == 0 {
				validSource = true
			}
			return nil
		})
		if err != nil {
			return err
		}

		if !validUser {
			return fmt.Errorf("User '%v' is not valid for this query", context.User)
		}

		if !validSource {
			return fmt.Errorf("Client '%v' is not valid for this query", context.Client)
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
	return DBConProtect.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists(pattern)
		id, _ := b.NextSequence()
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, id)
		return b.Put(key, context.Marshal())
	})
}

func uniqID(b *bolt.Bucket) []byte {
	buf := make([]byte, 8)
	id, _ := b.NextSequence()
	binary.BigEndian.PutUint64(buf, id)
	return buf
}
