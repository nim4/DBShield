package dbshield

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/dbms"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/training"
	"github.com/nim4/DBShield/dbshield/utils"
)

//initial boltdb database
func initModel(path string) {
	logger.Infof("Internal DB: %s", path)
	if training.DBCon == nil {
		training.DBCon, _ = bolt.Open(path, 0600, nil)
		training.DBCon.Update(func(tx *bolt.Tx) error {
			tx.CreateBucketIfNotExists([]byte("pattern"))
			tx.CreateBucketIfNotExists([]byte("abnormal"))
			b, _ := tx.CreateBucketIfNotExists([]byte("state"))
			v := b.Get([]byte("QueryCounter"))
			if v != nil {
				training.QueryCounter = binary.BigEndian.Uint64(v)
			}
			v = b.Get([]byte("AbnormalCounter"))
			if v != nil {
				training.AbnormalCounter = binary.BigEndian.Uint64(v)
			}
			return nil
		})
	}

	if config.Config.SyncInterval != 0 {
		training.DBCon.NoSync = true
		ticker := time.NewTicker(time.Second * time.Duration(config.Config.SyncInterval))
		go func() {
			for range ticker.C {
				training.DBCon.Sync()
			}
		}()
	}
}

func closeHandlers() {
	if training.DBCon != nil {
		training.DBCon.Update(func(tx *bolt.Tx) error {
			//Supplied value must remain valid for the life of the transaction
			qCount := make([]byte, 8)
			abCount := make([]byte, 8)

			b := tx.Bucket([]byte("state"))
			binary.BigEndian.PutUint64(qCount, training.QueryCounter)
			b.Put([]byte("QueryCounter"), qCount)

			binary.BigEndian.PutUint64(abCount, training.AbnormalCounter)
			b.Put([]byte("AbnormalCounter"), abCount)

			return nil
		})
		training.DBCon.Sync()
		training.DBCon.Close()
	}
	if logger.Output != nil {
		logger.Output.Close()
	}
}

//catching Interrupts
func signalHandler() {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt)
	<-term
	logger.Info("Shutting down...")
	//Closing open handler politely
	closeHandlers()
}

//initLogging redirect log output to file/stdout/stderr
func initLogging() {
	switch config.Config.LogPath {
	case "stdout":
		logger.Output = os.Stdout
	case "stderr":
		logger.Output = os.Stderr
	default:
		var err error
		logger.Output, err = os.OpenFile(config.Config.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Errorf("Error opening log file: %v", err))
		}
	}
	log.SetOutput(logger.Output)
	logger.Level = config.Config.LogLevel
}

//maps database name to corresponding struct
func dbNameToStruct(db string) (d utils.DBMS, err error) {
	switch strings.ToLower(db) {
	case "db2":
		d = &dbms.DB2{}
	case "mysql", "mariadb":
		d = &dbms.MySQL{}
	case "oracle":
		d = &dbms.Oracle{}
	case "postgres":
		d = &dbms.Postgres{}
	default:
		err = fmt.Errorf("Unknown DBMS: %s", db)
	}
	return
}

func handleClient(listenConn net.Conn, serverAddr *net.TCPAddr) error {
	db := utils.GenerateDBMS(config.Config.DB)
	logger.Debugf("Connected from: %s", listenConn.RemoteAddr())
	serverConn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		logger.Warning(err)
		listenConn.Close()
		return err
	}
	logger.Debugf("Connected to: %s", serverConn.RemoteAddr())
	db.SetSockets(listenConn, serverConn)
	db.SetReader(dbms.ReadPacket)
	err = db.Handler()
	if err != nil {
		logger.Warning(err)
	}
	return err
}
