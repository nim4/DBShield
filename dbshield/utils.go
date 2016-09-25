package dbshield

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/dbms"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/training"
	"github.com/nim4/DBShield/dbshield/utils"
)

func initModel() {
	var err error
	training.DBCon, err = bolt.Open(fmt.Sprintf("%s/%s_%s.db", config.Config.DBDir, config.Config.TargetIP, config.Config.DBType), 0600, nil)
	if err != nil {
		panic(err)
	}

	if err := training.DBCon.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("queries"))
		return err
	}); err != nil {
		panic(err)
	}
}

func initSignal() {
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt)
	go func() {
		for {
			<-term
			logger.Info("Shutting down...")
			if training.DBCon != nil {
				training.DBCon.Close()
			}
			if logger.Output != nil {
				logger.Output.Close()
			}
			os.Exit(0)
		}
	}()
}

//initLogging redirect log output to file
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
}

func dbNameToStruct(db string) (d utils.DBMS, err error) {
	switch strings.ToLower(db) {
	case "oracle":
		d = &dbms.Oracle{}
	case "mysql":
		d = &dbms.MySQL{}
	default:
		err = fmt.Errorf("Unknown DBMS: %s", db)
	}
	return
}
