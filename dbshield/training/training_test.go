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
	training.DBConLearning, err = bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}
	training.DBConProtect, err = bolt.Open(path+"2", 0600, nil)
	if err != nil {
		panic(err)
	}
}

func TestAddToTrainingSet(t *testing.T) {
	var err error
	c := sql.QueryContext{
		Query:    []byte("select * from test;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
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
		Query:    []byte("select * from test;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	c2 := sql.QueryContext{
		Query:    []byte("select * from user;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	training.AddToTrainingSet(c1)
	if !training.CheckQuery(c1) {
		t.Error("Expected false")
	}
	if training.CheckQuery(c2) {
		t.Error("Expected true")
	}

	tmpCon := training.DBConLearning
	defer func() {
		training.DBConLearning = tmpCon
	}()
	tmpfile, err := ioutil.TempFile("", "testdb")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()
	path := tmpfile.Name()
	training.DBConLearning, err = bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}
	if training.CheckQuery(c1) {
		t.Error("Expected false")
	}
}
