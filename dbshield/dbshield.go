/*
Package dbshield implements the database firewall functionality
*/
package dbshield

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/httpserver"
	"github.com/nim4/DBShield/dbshield/logger"
)

//Version of the library
var Version = "1.0-beta"

//Check config file and writes it to STDUT
func Check(configFile string) error {
	err := config.ParseConfig(configFile)
	if err != nil {
		return err
	}
	confJSON, err := json.MarshalIndent(config.Config, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(confJSON))
	return nil
}

func postConfig() (err error) {
	config.Config.DB, err = dbNameToStruct(config.Config.DBType)
	if err != nil {
		return err
	}

	if config.Config.ListenPort == 0 {
		config.Config.ListenPort = config.Config.DB.DefaultPort()
	}
	if config.Config.TargetPort == 0 {
		config.Config.TargetPort = config.Config.DB.DefaultPort()
	}

	err = config.Config.DB.SetCertificate(config.Config.TLSCertificate, config.Config.TLSPrivateKey)
	if err != nil {
		return
	}
	return
}

//Start the proxy
func Start(configFile string) (err error) {
	err = config.ParseConfig(configFile)
	if err != nil {
		return
	}

	err = postConfig()
	if err != nil {
		return err
	}
	initLogging()
	initModel()
	initSignal()

	logger.Infof("Config file: %s", configFile)
	logger.Infof("Listening: %s:%v (Threads: %v)",
		config.Config.ListenIP,
		config.Config.ListenPort,
		config.Config.Threads)
	logger.Infof("Backend: %s (%s:%v)",
		config.Config.DBType,
		config.Config.TargetIP,
		config.Config.TargetPort)
	logger.Infof("Protect: %v", !config.Config.Learning)

	serverAddr, err := net.ResolveTCPAddr("tcp", config.Config.TargetIP+":"+strconv.Itoa(int(config.Config.TargetPort)))
	if err != nil {
		return
	}
	l, err := net.Listen("tcp", config.Config.ListenIP+":"+strconv.Itoa(int(config.Config.ListenPort)))
	if err != nil {
		return
	}
	// Close the listener when the application closes.
	defer l.Close()

	if config.Config.HTTP {
		go httpserver.Serve()
	}
	for {
		var listenConn net.Conn
		// Listen for an incoming connection.
		listenConn, err = l.Accept()
		if err != nil {
			return
		}

		go handleClient(listenConn, serverAddr)
	}
}
