package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/utils"
	"github.com/spf13/viper"
)

type mask struct {
	MatchExp         *regexp.Regexp
	ReplaceExp       []byte
	PaddingCharacter []byte
}

//Config hold parsed configurations
var Config = struct {
	Learning    bool
	CheckUser   bool
	CheckSource bool

	LogLevel uint
	LogPath  string

	DBType string
	DB     utils.DBMS

	Threads uint
	DBDir   string

	ListenIP   string
	ListenPort uint

	TargetIP   string
	TargetPort uint

	TLSPrivateKey  string
	TLSCertificate string

	Action func(net.Conn) error

	//Key-> database.table.column
	Masks map[string]mask
}{}

//ParseConfig and return error if its not valid
func ParseConfig(configFile string) error {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return fmt.Errorf("Fatal error config file: %s \n", err)
	}
	Config.DBDir = viper.GetString("dbDir")
	os.Mkdir(Config.DBDir, 0740)
	//set default values
	Config.DBType = viper.GetString("dbms")
	switch viper.GetString("mode") {
	case "protect":
		Config.Learning = false
	case "learning":
		Config.Learning = true
	default:
		return errors.New("Invalid cofiguration: mode")
	}

	for _, check := range strings.Split(viper.GetString("additionalChecks"), ",") {
		switch check {
		case "user":
			Config.CheckUser = true
		case "source":
			Config.CheckSource = true
		default:
			return errors.New("Invalid cofiguration: additionalChecks")
		}
	}

	switch viper.GetString("action") {
	case "drop": //Close the connection
		Config.Action = utils.ActionDrop
	case "pass": //Pass the query to server
		Config.Action = nil
	default:
		return errors.New("Invalid cofiguration: action")
	}

	Config.Threads = uint(viper.GetInt("threads"))
	if Config.Threads < 1 {
		return errors.New("Invalid cofiguration: threads")
	}
	Config.LogLevel = uint(viper.GetInt("logLevel"))
	logger.Level = Config.LogLevel

	Config.LogPath = viper.GetString("logPath")

	Config.ListenIP = viper.GetString("listenIP")
	Config.TargetIP = viper.GetString("targetIP")

	Config.ListenPort = uint(viper.GetInt("listenPort"))
	Config.TargetPort = uint(viper.GetInt("targetPort"))

	Config.TLSPrivateKey = viper.GetString("tlsPrivateKey")
	Config.TLSCertificate = viper.GetString("tlsCertificate")

	Config.Masks = make(map[string]mask)
	tmpMasks := viper.Get("masks").([]interface{})
	for _, tm := range tmpMasks {
		m := tm.(map[interface{}]interface{})
		key := m["database"].(string) + "." + m["table"].(string) + "." + m["column"].(string)
		padding := []byte(m["paddingCharacter"].(string))
		Config.Masks[key] = mask{
			MatchExp:         regexp.MustCompile(m["matchRegEx"].(string)),
			ReplaceExp:       []byte(m["replaceRegEx"].(string)),
			PaddingCharacter: []byte{padding[0]},
		}
	}
	return nil
}
