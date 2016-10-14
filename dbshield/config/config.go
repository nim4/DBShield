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

func strConfig(dst *string, key, defaultValue string) {
	ret := ""
	if viper.IsSet(key) {
		ret = viper.GetString(key)
	} else {
		logger.Infof("'%s' not configured, assuming: %s", key, defaultValue)
		ret = defaultValue
	}
	dst = &ret
}

func intConfig(dst *uint, key string, defaultValue, min uint) error {
	var ret uint
	if viper.IsSet(key) {
		ret = uint(viper.GetInt(key))
		if ret < min {
			return fmt.Errorf("Invalid '%s' cofiguration: %v\n", key, ret)
		}
	} else {
		logger.Infof("'%s' not configured, assuming: %s", key, defaultValue)
		ret = defaultValue
	}
	dst = &ret
	return nil
}

//ParseConfig and return error if its not valid
func ParseConfig(configFile string) error {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return fmt.Errorf("Fatal error - config file: %s \n", err)
	}

	//set default values
	if viper.IsSet("mode") {
		switch viper.GetString("mode") {
		case "protect":
			Config.Learning = false
		case "learning":
			Config.Learning = true
		default:
			return errors.New("Invalid 'mode' cofiguration: " + viper.GetString("mode"))
		}
	} else {
		logger.Infof("'mode' not configured, assuming: learning")
		Config.Learning = true
	}

	if viper.IsSet("additionalChecks") {
		for _, check := range strings.Split(viper.GetString("additionalChecks"), ",") {
			switch check {
			case "user":
				Config.CheckUser = true
			case "source":
				Config.CheckSource = true
			default:
				return errors.New("Invalid 'additionalChecks' cofiguration: " + check)
			}
		}
	}

	if !Config.Learning {
		if viper.IsSet("action") {
			switch viper.GetString("action") {
			case "drop": //Close the connection
				Config.Action = utils.ActionDrop
			case "pass": //Pass the query to server
				Config.Action = nil
			default:
				return errors.New("Invalid 'action' cofiguration: " + viper.GetString("action"))
			}
		} else {
			logger.Infof("'action' not configured, assuming: drop")
			Config.Action = utils.ActionDrop
		}
	}

	if viper.IsSet("targetIP") {
		Config.TargetIP = viper.GetString("targetIP")
	} else {
		return errors.New("Invalid 'targetIP' cofiguration: " + viper.GetString("targetIP"))
	}

	if viper.IsSet("tlsPrivateKey") {
		Config.TLSPrivateKey = viper.GetString("tlsPrivateKey")
	} else {
		return errors.New("Invalid 'tlsPrivateKey' cofiguration: " + viper.GetString("tlsPrivateKey"))
	}

	if viper.IsSet("tlsCertificate") {
		Config.TargetIP = viper.GetString("tlsCertificate")
	} else {
		return errors.New("Invalid 'tlsCertificate' cofiguration: " + viper.GetString("tlsCertificate"))
	}

	//String values
	strConfig(&Config.DBDir, "dbDir", os.TempDir()+"/model")
	os.Mkdir(Config.DBDir, 0740) //Make dbDir, just in case its not there

	strConfig(&Config.DBType, "dbms", "mysql")

	strConfig(&Config.ListenIP, "listenIP", "0.0.0.0")

	strConfig(&Config.LogPath, "logPath", "stderr")

	//Integer values
	err = intConfig(&Config.Threads, "threads", 4, 1)
	if err != nil {
		return err
	}

	err = intConfig(&Config.LogLevel, "logLevel", 2, 0)
	if err != nil {
		return err
	}

	err = intConfig(&Config.ListenPort, "listenPort", 0, 0)
	if err != nil {
		return err
	}

	err = intConfig(&Config.TargetPort, "targetPort", 0, 0)
	if err != nil {
		return err
	}

	Config.Masks = make(map[string]mask)
	if viper.IsSet("masks") {
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
	}
	return nil
}
