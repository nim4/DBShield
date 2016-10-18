package training_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/sql"
	"github.com/nim4/DBShield/dbshield/training"
)

func TestMain(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "testdb")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()
	path := tmpfile.Name()
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
	err = training.AddToTrainingSet(sql.QueryContext{})
	if err == nil {
		t.Error("Expected error")
	}
}

func TestCheckQuery(t *testing.T) {
	c1 := sql.QueryContext{
		Query:    "select * from test;",
		Database: "test",
		User:     "test",
		Client:   "127.0.0.1",
		Time:     time.Now().Unix(),
	}
	c2 := sql.QueryContext{
		Query:    "select * from user;",
		Database: "test",
		User:     "test",
		Client:   "127.0.0.1",
		Time:     time.Now().Unix(),
	}
	training.AddToTrainingSet(c1)
	if !training.CheckQuery(c1) {
		t.Error("Expected false")
	}
	if training.CheckQuery(c2) {
		t.Error("Expected true")
	}

	tmpCon := training.DBCon
	defer func() {
		training.DBCon = tmpCon
	}()
	tmpfile, err := ioutil.TempFile("", "testdb")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()
	path := tmpfile.Name()
	training.DBCon, err = bolt.Open(path, 0600, nil)
	if training.CheckQuery(c1) {
		t.Error("Expected false")
	}
}
