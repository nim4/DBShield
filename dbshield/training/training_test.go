package training_test

import (
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/sql"
	"github.com/nim4/DBShield/dbshield/training"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard) // Avoid log outputs
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
	training.DBCon.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("pattern"))
		tx.CreateBucket([]byte("abnormal"))
		tx.CreateBucket([]byte("state"))
		return nil
	})
	m.Run()
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
	if err != nil {
		panic(err)
	}
	err = training.AddToTrainingSet(c)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestCheckQuery(t *testing.T) {
	config.Config.CheckUser = true
	config.Config.CheckSource = true
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
	if err != nil {
		panic(err)
	}
	training.DBCon.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("pattern"))
		return err
	})
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic")
		}
	}()
	training.CheckQuery(c1)
}

func BenchmarkAddToTrainingSet(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	c := sql.QueryContext{
		Query:    []byte("select * from test;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	for i := 0; i < b.N; i++ {
		training.AddToTrainingSet(c)
	}
}
