package logger_test

import (
	"testing"

	"github.com/nim4/DBShield/dbshield/logger"
)

func TestLogger(t *testing.T) {
	format := "%s"
	msg := "Test"
	logger.Level = 7
	logger.Debug(msg)
	logger.Debugf(format, msg)
	logger.Info(msg)
	logger.Infof(format, msg)
	logger.Warning(msg)
	logger.Warningf(format, msg)
}

func BenchmarkDebugf(b *testing.B) {
	logger.Level = 7
	for i := 0; i < b.N; i++ {
		logger.Debugf("%s", "t")
	}
}

func BenchmarkDebug(b *testing.B) {
	logger.Level = 7
	for i := 0; i < b.N; i++ {
		logger.Debug("t")
	}
}
