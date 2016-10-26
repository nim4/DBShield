package dbshield

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
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
func initModel() {
	var err error
	path := path.Join(config.Config.DBDir, config.Config.TargetIP+"_"+config.Config.DBType)
	logger.Infof("Internal DB: %s", path)
	if training.DBConLearning == nil {
		training.DBConLearning, err = bolt.Open(path+"_learning.db", 0600, nil)
		if err != nil {
			panic(err)
		}
	}

	if training.DBConProtect == nil {
		training.DBConProtect, err = bolt.Open(path+"_abnormal.db", 0600, nil)
		if err != nil {
			panic(err)
		}
	}

	if config.Config.SyncInterval != 0 {
		training.DBConLearning.NoSync = true
		training.DBConProtect.NoSync = true
		ticker := time.NewTicker(time.Second * time.Duration(config.Config.SyncInterval))
		go func() {
			for range ticker.C {
				training.DBConLearning.Sync()
				training.DBConProtect.Sync()
			}
		}()
	}
}

func closeHandlers() {
	if training.DBConLearning != nil {
		training.DBConLearning.Sync()
		training.DBConLearning.Close()
	}
	if training.DBConProtect != nil {
		training.DBConProtect.Sync()
		training.DBConProtect.Close()
	}
	if logger.Output != nil {
		logger.Output.Close()
	}
	// pprof.StopCPUProfile()
}

//catching Interrupts
func initSignal() {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt)
	go func() {
		<-term
		logger.Info("Shutting down...")
		//Closing open handler politely
		closeHandlers()
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
	db := config.Config.DB
	logger.Infof("Connected from: %s", listenConn.RemoteAddr())
	serverConn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		logger.Warning(err)
		listenConn.Close()
		return err
	}
	logger.Infof("Connected to: %s", serverConn.RemoteAddr())
	defer db.Close()
	db.SetSockets(listenConn, serverConn)
	db.SetReader(dbms.ReadPacket)
	err = db.Handler()
	if err != nil {
		logger.Warning(err)
	}
	return err
}
