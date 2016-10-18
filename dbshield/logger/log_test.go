package logger_test

import (
	"testing"

	"github.com/nim4/DBShield/dbshield/logger"
)

func TestLogger(t *testing.T) {
	format := "%s"
	msg := "Test"
	logger.Level = 2
	logger.Debug(msg)
	logger.Debugf(format, msg)
	logger.Info(msg)
	logger.Infof(format, msg)
	logger.Warning(msg)
	logger.Warningf(format, msg)
}
