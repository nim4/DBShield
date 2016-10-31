/*
Package dbshield implements the database firewall functionality
*/
package dbshield

import (
	"encoding/json"
	"fmt"
	"net"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/httpserver"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
	"github.com/nim4/DBShield/dbshield/training"
)

//Version of the library
var Version = "1.0.0-beta4"

var configFile string

//SetConfigFile of DBShield
func SetConfigFile(cf string) error {
	configFile = cf
	err := config.ParseConfig(configFile)
	if err != nil {
		return err
	}
	return postConfig()
}

//Check config file and writes it to STDUT
func Check() error {
	confJSON, err := json.MarshalIndent(config.Config, "", "    ")
	fmt.Println(string(confJSON))
	return err
}

//Patterns lists the captured patterns
func Patterns() {
	initModel(
		path.Join(config.Config.DBDir,
			config.Config.TargetIP+"_"+config.Config.DBType) + ".db")

	training.DBCon.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pattern"))
		if b != nil {
			return b.ForEach(func(k, v []byte) error {
				if strings.Index(string(k), "_client_") == -1 && strings.Index(string(k), "_user_") == -1 {
					fmt.Printf(
						`-----Pattern: 0x%x
Sample: %s
`,
						k,
						v,
					)
				}
				return nil
			})
		}
		return nil
	})
}

//Abnormals detected querties
func Abnormals() {
	initModel(
		path.Join(config.Config.DBDir,
			config.Config.TargetIP+"_"+config.Config.DBType) + ".db")

	training.DBCon.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("abnormal"))
		if b != nil {
			return b.ForEach(func(k, v []byte) error {
				var c sql.QueryContext
				c.Unmarshal(v)
				fmt.Printf("[%s] [User: %s] [Database: %s] %s\n",
					c.Time.Format(time.RFC1123),
					c.User,
					c.Database,
					c.Query)
				return nil
			})
		}
		return nil
	})
}

func postConfig() (err error) {

	config.Config.DB, err = dbNameToStruct(config.Config.DBType)
	if err != nil {
		return err
	}

	tmpDBMS := generateDBMS()
	if config.Config.ListenPort == 0 {
		config.Config.ListenPort = tmpDBMS.DefaultPort()
	}
	if config.Config.TargetPort == 0 {
		config.Config.TargetPort = tmpDBMS.DefaultPort()
	}
	return
}

func mainListner() error {
	if config.Config.HTTP {
		logger.Infof("Web interface on https://%s/", config.Config.HTTPAddr)
		go httpserver.Serve()
	}
	serverAddr, _ := net.ResolveTCPAddr("tcp", config.Config.TargetIP+":"+strconv.Itoa(int(config.Config.TargetPort)))
	l, err := net.Listen("tcp", config.Config.ListenIP+":"+strconv.Itoa(int(config.Config.ListenPort)))
	if err != nil {
		return err
	}
	// Close the listener when the application closes.
	defer l.Close()

	for {
		// Listen for an incoming connection.
		listenConn, err := l.Accept()
		if err != nil {
			logger.Warningf("Error accepting connection: %v", err)
			continue
		}
		go handleClient(listenConn, serverAddr)
	}
}

//Start the proxy
func Start() (err error) {
	initModel(
		path.Join(config.Config.DBDir,
			config.Config.TargetIP+"_"+config.Config.DBType) + ".db")

	initLogging()
	logger.Infof("Config file: %s", configFile)
	logger.Infof("Listening: %s:%v",
		config.Config.ListenIP,
		config.Config.ListenPort)
	logger.Infof("Backend: %s (%s:%v)",
		config.Config.DBType,
		config.Config.TargetIP,
		config.Config.TargetPort)
	logger.Infof("Protect: %v", !config.Config.Learning)
	go mainListner()
	signalHandler()
	return nil
}
