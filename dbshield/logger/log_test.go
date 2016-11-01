package logger_test

import (
	"os"
	"testing"

	"github.com/nim4/DBShield/dbshield/logger"
)

func TestLogger(t *testing.T) {
	logger.Init("stderr", 7)
	format := "%s"
	msg := "Test"
	logger.Debug(msg)
	logger.Debugf(format, msg)
	logger.Info(msg)
	logger.Infof(format, msg)
	logger.Warning(msg)
	logger.Warningf(format, msg)
}

func TestInit(t *testing.T) {

	err := logger.Init(os.TempDir(), 0)
	if err == nil {
		t.Error("Expected error")
	}

	err = logger.Init("stdout", 0)
	if err != nil {
		t.Error("Got error", err)
	}

	err = logger.Init("stderr", 7)
	if err != nil {
		t.Error("Got error", err)
	}
}

func BenchmarkDebugf(b *testing.B) {
	logger.Init("stderr", 7)
	for i := 0; i < b.N; i++ {
		logger.Debugf("%s", "t")
	}
}

func BenchmarkDebug(b *testing.B) {
	logger.Init("stderr", 7)
	for i := 0; i < b.N; i++ {
		logger.Debug("t")
	}
}
