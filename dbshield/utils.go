package dbshield

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/dbms"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/training"
	"github.com/nim4/DBShield/dbshield/utils"
)

//initial boltdb database
func initModel() {
	var err error
	path := path.Join(config.Config.DBDir, config.Config.TargetIP+"_"+config.Config.DBType+".db")
	logger.Infof("Internal DB: %s", path)
	training.DBCon, err = bolt.Open(path, 0600, &bolt.Options{
		Timeout:    5,
		NoGrowSync: false,
	})
	if err != nil {
		panic(err)
	}

	training.DBCon.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("queries"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("abnormal"))
		return err
	})
}

//catching Interrupts
func initSignal() {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt)
	go func() {
		<-term
		logger.Info("Shutting down...")
		//Closing open handler politely
		if training.DBCon != nil {
			training.DBCon.Close()
		}
		if logger.Output != nil {
			logger.Output.Close()
		}
		os.Exit(0)
	}()
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
	case "mysql":
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
	db := config.Config.DB
	logger.Infof("Connected from: %s", listenConn.RemoteAddr())
	serverConn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		logger.Warning(err)
		listenConn.Close()
		return err
	}
	logger.Infof("Connected to: %s", serverConn.RemoteAddr())
	db.SetSockets(listenConn, serverConn)
	db.SetReader(dbms.ReadPacket)
	err = db.Handler()
	if err != nil {
		logger.Warning(err)
		return err
	}
	db.Close()
	return nil
}
