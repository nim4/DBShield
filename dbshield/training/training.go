package training

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
)

//DBCon holds a pointer to our local database connection
var DBCon *bolt.DB

//AddToTrainingSet records query context in local database
func AddToTrainingSet(context sql.QueryContext) {

	pattern := sql.Pattern(context.Query)
	//logger.Debug("Pattern", pattern)

	if err := DBCon.Update(func(tx *bolt.Tx) error {
		var contextArray []sql.QueryContext
		b := tx.Bucket([]byte("queries"))
		if b == nil {
			panic(errors.New("Bucket not found"))
		}
		v := b.Get(pattern)
		if v != nil {
			//logger.Debug("Key found: ", string(v))
			if err := json.Unmarshal(v, &contextArray); err != nil {
				return err
			}
		}
		contextArray = append(contextArray, context)
		//logger.Debug("Context Array:", contextArray)
		encoded, err := json.Marshal(contextArray)
		//logger.Debug("JSON:", string(encoded))
		if err != nil {
			return err
		}
		if err := b.Put(pattern, encoded); err != nil {
			return err
		}

		return nil
	}); err != nil {
		logger.Warning(err)
	}
}

//CheckQuery pattern, returns true if it finds the pattern
//We should keep it as fast as possible
func CheckQuery(context sql.QueryContext) bool {
	pattern := sql.Pattern(context.Query)
	if err := DBCon.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("queries"))
		if b == nil {
			panic(errors.New("Bucket not found"))
		}
		v := b.Get(pattern)
		if v != nil {
			var contextArray []sql.QueryContext
			if err := json.Unmarshal(v, &contextArray); err != nil {
				return err
			}
			if config.Config.CheckUser || config.Config.CheckSource {
				validUser := false
				validAddr := false
				for _, con := range contextArray {
					if config.Config.CheckUser && con.User == context.User {
						validUser = true
					}
					if config.Config.CheckSource && con.Client == context.Client {
						validAddr = true
					}
					if validUser && validAddr {
						return nil
					}
				}
				if !validUser {
					return fmt.Errorf("User '%v' is not valid for this query", context.User)
				}
				return fmt.Errorf("Client '%v' is not valid for this query", context.Client)
			}
			return nil
		}
		return fmt.Errorf("Pattern not found: %v", pattern)
	}); err != nil {
		logger.Warning(err, "("+context.Query+")")
		return false
	}
	return true
}
