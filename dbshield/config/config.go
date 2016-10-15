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

//Config holds current configurations
var Config struct {
	Learning    bool
	CheckUser   bool
	CheckSource bool

	LogLevel uint
	LogPath  string

	DBType string
	DB     utils.DBMS `json:"-"`

	Threads uint
	DBDir   string

	ListenIP   string
	ListenPort uint

	TargetIP   string
	TargetPort uint

	TLSPrivateKey  string
	TLSCertificate string

	Action     string
	ActionFunc func(net.Conn) error `json:"-"`

	//Key-> database.table.column
	//Masks map[string]mask
}

func strConfig(key string) (ret string, err error) {
	if viper.IsSet(key) {
		ret = viper.GetString(key)
		return
	}
	err = fmt.Errorf("Invalid '%s' cofiguration", key)
	return
}

func strConfigDefualt(key, defaultValue string) (ret string) {
	if viper.IsSet(key) {
		ret = viper.GetString(key)
		return
	}
	logger.Infof("'%s' not configured, assuming: %s", key, defaultValue)
	ret = defaultValue
	return
}

func intConfig(key string, defaultValue, min uint) (ret uint, err error) {
	if viper.IsSet(key) {
		ret = uint(viper.GetInt(key))
		if ret < min {
			err = fmt.Errorf("Invalid '%s' cofiguration: %v\n", key, ret)
			return
		}
		return
	}
	logger.Infof("'%s' not configured, assuming: %s", key, defaultValue)
	ret = defaultValue
	return
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
			Config.Action = viper.GetString("action")
			switch Config.Action {
			case "drop": //Close the connection
				Config.ActionFunc = utils.ActionDrop
			case "pass": //Pass the query to server
				Config.ActionFunc = nil
			default:
				return errors.New("Invalid 'action' cofiguration: " + Config.Action)
			}
		} else {
			logger.Infof("'action' not configured, assuming: drop")
			Config.ActionFunc = utils.ActionDrop
		}
	}

	Config.TargetIP, err = strConfig("targetIP")
	if err != nil {
		return err
	}

	Config.TLSPrivateKey, err = strConfig("tlsPrivateKey")
	if err != nil {
		return err
	}

	Config.TLSCertificate, err = strConfig("tlsCertificate")
	if err != nil {
		return err
	}

	//String values
	Config.DBDir = strConfigDefualt("dbDir", os.TempDir()+"/model")
	os.MkdirAll(Config.DBDir, 0740) //Make dbDir, just in case its not there

	Config.DBType = strConfigDefualt("dbms", "mysql")

	Config.ListenIP = strConfigDefualt("listenIP", "0.0.0.0")

	Config.LogPath = strConfigDefualt("logPath", "stderr")

	//Integer values
	Config.Threads, err = intConfig("threads", 4, 1)
	if err != nil {
		return err
	}

	Config.LogLevel, err = intConfig("logLevel", 2, 0)
	if err != nil {
		return err
	}

	Config.ListenPort, err = intConfig("listenPort", 0, 0)
	if err != nil {
		return err
	}

	Config.TargetPort, err = intConfig("targetPort", 0, 0)
	if err != nil {
		return err
	}
	/* Masking
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
	*/
	return nil
}
