// +build !windows

package dbshield

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/sql"
	"github.com/nim4/DBShield/dbshield/training"
)

func TestMain(m *testing.M) {
	os.Chdir("../")
	m.Run()
}

func TestSetConfigFile(t *testing.T) {
	err := SetConfigFile("Invalid.yml")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestShowConfig(t *testing.T) {
	SetConfigFile("conf/dbshield.yml")
	err := ShowConfig()
	if err != nil {
		t.Error("Got error", err)
	}
}

func TestPurge(t *testing.T) {
	SetConfigFile("conf/dbshield.yml")
	err := Purge()
	if err == nil {
		t.Error("Expected error")
	}
}

func TestPostConfig(t *testing.T) {
	SetConfigFile("conf/dbshield.yml")
	config.Config.DBType = "Invalid"
	err := postConfig()
	if err == nil {
		t.Error("Expected error")
	}

	config.Config.ListenPort = 0
	config.Config.DBType = "mysql"
	err = postConfig()
	if err != nil {
		t.Error("Expected nil got ", err)
	}
}

func TestEveryThing(t *testing.T) {
	closeHandlers()
	SetConfigFile("conf/dbshield.yml")
	//It should fail if port is already open
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Config.ListenIP, config.Config.ListenPort))
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	err = mainListner()
	if err == nil {
		t.Error("Expected error")
	}

	go func() {
		timer := time.NewTimer(time.Second * 2)
		<-timer.C
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	err = Start()
	if err != nil {
		t.Error("Got error", err)
	}
	file, _ := ioutil.TempFile(os.TempDir(), "tempDB")
	defer os.Remove(file.Name())
	training.DBCon, _ = bolt.Open(file.Name(), 0600, nil)
	training.DBCon.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("pattern"))
		tx.CreateBucket([]byte("abnormal"))
		tx.CreateBucket([]byte("state"))
		return nil
	})
	c := sql.QueryContext{
		Query:    []byte("select * from test;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	training.CheckQuery(c)
	err = training.AddToTrainingSet(c)
	if err != nil {
		t.Error("Got error", err)
	}

	Patterns()
	Abnormals()

	//Test without bucket
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
	Patterns()
	Abnormals()
}
