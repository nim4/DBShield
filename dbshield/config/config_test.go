package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestParseConfig(t *testing.T) {
	err := ParseConfig("../../conf/dbshield.yml")
	if err != nil {
		t.Error("Got error", err)
	}

	err = ParseConfig("Invalid.yml")
	if err == nil {
		t.Error("Expected error")
	}

	viper.Set("mode", "Invalid") //make configGeneral fail
	err = ParseConfig("../../conf/dbshield.yml")
	if err == nil {
		t.Error("Expected error")
	}

	viper.Set("mode", "protect") //make configProtect fail
	viper.Set("action", "xyz")
	err = ParseConfig("../../conf/dbshield.yml")
	if err == nil {
		t.Error("Expected error")
	}
	viper.Set("action", "drop")

	viper.Set("logLevel", -1) //make configLog fail
	err = ParseConfig("../../conf/dbshield.yml")
	if err == nil {
		t.Error("Expected error")
	}
	viper.Set("logLevel", 0)

	viper.Set("http", true) //make configHTTP fail
	viper.Set("httpPort", -1)
	err = ParseConfig("../../conf/dbshield.yml")
	if err == nil {
		t.Error("Expected error")
	}
	viper.Set("logLevel", 0)
}

//Check the cases which are not in default config

func TestStrConfigDefualt(t *testing.T) {
	ret := strConfigDefualt("Invalid", "X")
	if ret != "X" {
		t.Error("Expected X got", ret)
	}
}

func TestIntConfig(t *testing.T) {
	ret, err := intConfig("Invalid", 1, 1)
	if err != nil || ret != 1 {
		t.Error("Expected 1 or error got", ret, err)
	}
}

func TestConfigGeneral(t *testing.T) {
	viper.Reset()
	viper.Set("targetIP", "127.0.0.1")
	viper.Set("tlsPrivateKey", "key")
	viper.Set("tlsCertificate", "cert")

	viper.Set("mode", "protect")
	err := configGeneral()
	if err != nil {
		t.Error("Got error", err)
	}

	viper.Set("mode", "Invalid")
	err = configGeneral()
	if err == nil {
		t.Error("Expected error")
	}

	viper.Set("mode", nil)
	err = configGeneral()
	if err != nil {
		t.Error("Got error", err)
	}
	viper.Set("mode", "learning")

	viper.Set("threads", "Invalid")
	err = configGeneral()
	if err == nil {
		t.Error("Expected error")
	}
	viper.Set("threads", 4)

	viper.Set("listenPort", -1)
	err = configGeneral()
	if err == nil {
		t.Error("Expected error")
	}
	viper.Set("listenPort", 0)

	viper.Set("targetPort", -1)
	err = configGeneral()
	if err == nil {
		t.Error("Expected error")
	}
	viper.Set("targetPort", 0)

	viper.Reset()
	err = configGeneral()
	if err == nil {
		t.Error("Expected error")
	}
	viper.Set("targetIP", "127.0.0.1")

	viper.Reset()
	viper.Set("targetIP", "127.0.0.1")
	viper.Set("tlsPrivateKey", "key")
	err = configGeneral()
	if err == nil {
		t.Error("Expected error")
	}

	viper.Reset()
	viper.Set("targetIP", "127.0.0.1")
	viper.Set("tlsCertificate", "cert")
	err = configGeneral()
	if err == nil {
		t.Error("Expected error")
	}
	viper.Set("tlsPrivateKey", "key")

	// Can't make directory named after file.
	fpath := os.TempDir() + "/file"
	f, err := os.Create(fpath)
	if err != nil {
		t.Fatalf("create %q: %s", fpath, err)
	}
	f.Close()

	viper.Set("dbDir", fpath)
	err = configGeneral()
	if err == nil {
		t.Error("Expected error")
	}
}

func TestConfigProtect(t *testing.T) {
	viper.Set("action", "drop")
	err := configProtect()
	if err != nil {
		t.Error("Got error", err)
	}

	viper.Set("action", "pass")
	err = configProtect()
	if err != nil {
		t.Error("Got error", err)
	}

	viper.Set("action", "Invalid")
	err = configProtect()
	if err == nil {
		t.Error("Expected error")
	}

	viper.Set("action", nil)
	err = configProtect()
	if err != nil {
		t.Error("Got error", err)
	}

	viper.Set("additionalChecks", "Invalid")
	err = configProtect()
	if err == nil {
		t.Error("Expected error")
	}
}

func TestConfigHTTP(t *testing.T) {
	viper.Set("http", true)
	viper.Set("httpPort", "Invalid")
	err := configHTTP()
	if err == nil {
		t.Error("Expected error")
	}
}
