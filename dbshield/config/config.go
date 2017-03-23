package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/utils"
	"github.com/spf13/viper"
)

type mask struct {
	MatchExp         *regexp.Regexp
	ReplaceExp       []byte
	PaddingCharacter []byte
}

//Configurations structure to hold user configurations
type Configurations struct {
	Learning    bool
	CheckUser   bool
	CheckSource bool

	LogLevel uint
	LogPath  string

	DBType string
	DB     uint `json:"-"`

	DBDir string

	ListenIP   string
	ListenPort uint

	TargetIP   string
	TargetPort uint

	TLSPrivateKey  string
	TLSCertificate string

	HTTP         bool
	HTTPSSL      bool
	HTTPAddr     string
	HTTPPassword string

	Action     string
	ActionFunc func() error `json:"-"`

	Timeout time.Duration

	SyncInterval time.Duration
	//Key-> database.table.column
	//Masks map[string]mask
}

//Config holds current configs
var Config Configurations

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
		tmp := viper.GetInt(key)
		if tmp < 0 {
			err = fmt.Errorf("Invalid '%s' cofiguration: %v\n", key, tmp)
			return
		}
		ret = uint(tmp)
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

func configGeneral() (err error) {
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

	Config.ListenPort, err = intConfig("listenPort", 0, 0)
	if err != nil {
		return err
	}

	Config.TargetPort, err = intConfig("targetPort", 0, 0)
	if err != nil {
		return err
	}
	Config.TargetIP, err = strConfig("targetIP")
	if err != nil {
		return err
	}

	//String values
	Config.TLSPrivateKey, err = strConfig("tlsPrivateKey")
	if err != nil {
		return err
	}

	Config.TLSCertificate, err = strConfig("tlsCertificate")
	if err != nil {
		return err
	}

	Config.DBDir = strConfigDefualt("dbDir", os.TempDir()+"/model")
	err = os.MkdirAll(Config.DBDir, 0740) //Make dbDir, just in case its not there
	if err != nil {
		return err
	}

	Config.DBType = strConfigDefualt("dbms", "mysql")

	Config.ListenIP = strConfigDefualt("listenIP", "0.0.0.0")

	if timeout := viper.GetString("timeout"); timeout != "" {
		Config.Timeout, err = time.ParseDuration(timeout)
		if err != nil {
			return err
		}
	} else {
		Config.Timeout = 5 * time.Second
	}

	if syn := viper.GetString("syncInterval"); syn != "" {
		Config.SyncInterval, err = time.ParseDuration(syn)
		if err != nil {
			return err
		}
	} else {
		Config.SyncInterval = 5 * time.Second
	}
	return nil
}

func configProtect() error {
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
	return nil
}

func configLog() error {
	var err error
	Config.LogPath = strConfigDefualt("logPath", "stderr")
	Config.LogLevel, err = intConfig("logLevel", 3, 0)
	return err
}

func configHTTP() error {
	Config.HTTP = viper.GetBool("http")
	if Config.HTTP {
		Config.HTTPPassword = viper.GetString("httpPassword")
		httpIP := strConfigDefualt("httpIP", "127.0.0.1")
		httpPort, err := intConfig("httpPort", 8070, 1)
		if err != nil {
			return err
		}
		Config.HTTPSSL = viper.GetBool("httpSSL")
		Config.HTTPAddr = fmt.Sprintf("%s:%d", httpIP, httpPort)
	}
	return nil
}

//ParseConfig and return error if its not valid
func ParseConfig(configFile string) error {
	Config = Configurations{} // Reset configs
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig() // Read the config file
	if err != nil {
		return fmt.Errorf("Fatal error - config file: %s \n", err)
	}
	err = configGeneral()
	if err != nil {
		return err
	}

	if !Config.Learning {
		err = configProtect()
		if err != nil {
			return err
		}
	}

	err = configLog()
	if err != nil {
		return err
	}

	err = configHTTP()
	if err != nil {
		return err
	}
	return nil
}
