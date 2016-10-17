package training_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/sql"
	"github.com/nim4/DBShield/dbshield/training"
)

func initTempDB() {
	tmpfile, err := ioutil.TempFile("", "testdb")
	if err != nil {
		panic(err)
	}
	path := tmpfile.Name()
	tmpfile.Close()
	training.DBCon, err = bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}
	if err := training.DBCon.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("queries"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("abnormal"))
		return err
	}); err != nil {
		panic(err)
	}
}

func TestAddToTrainingSet(t *testing.T) {
	var err error
	initTempDB()
	c := sql.QueryContext{
		Query:    "select * from test;",
		Database: "test",
		User:     "test",
		Client:   "127.0.0.1",
		Time:     time.Now().Unix(),
	}
	err = training.AddToTrainingSet(c)
	if err != nil {
		t.Error("Not Expected error", err)
	}
}

func TestCheckQuery(t *testing.T) {
	var err error
	initTempDB()
	c := sql.QueryContext{
		Query:    "select * from test;",
		Database: "test",
		User:     "test",
		Client:   "127.0.0.1",
		Time:     time.Now().Unix(),
	}
	if training.CheckQuery(c) {
		t.Error("Not Expected error", err)
	}
	training.AddToTrainingSet(c)
	if !training.CheckQuery(c) {
		t.Error("Not Expected error", err)
	}
}
