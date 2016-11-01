package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	flagWarning = 0x1
	flagInfo    = 0x2
	flagDebug   = 0x4
)

var (
	warn  bool
	info  bool
	debug bool
)

//Output of logging functions
var Output *os.File

//Init the log output and logging level
func Init(path string, level uint) error {
	switch path {
	case "stdout":
		Output = os.Stdout
	case "stderr":
		Output = os.Stderr
	default:
		var err error
		Output, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("Error opening log file: %v", err)
		}
	}
	log.SetOutput(Output)
	warn = level&flagWarning == flagWarning
	info = level&flagInfo == flagInfo
	debug = level&flagDebug == flagDebug
	return nil
}

func println(title string, msg ...interface{}) {
	arg := append([]interface{}{title}, msg...)
	log.Println(arg...)
}

func printFormat(format string, msg ...interface{}) {
	log.Printf(format, msg...)
}

//Debug level > 1 logging with "[DEBUG]" prefix
func Debug(msg ...interface{}) {
	if debug {
		println("[DEBUG]", msg...)
	}
}

//Debugf level > 1 format logging with "[DEBUG]" prefix
func Debugf(format string, msg ...interface{}) {
	if debug {
		printFormat("[DEBUG] "+format, msg...)
	}
}

//Info level > 0 logging with "[INFO]" prefix
func Info(msg ...interface{}) {
	if info {
		println("[INFO] ", msg...)
	}
}

//Infof level > 0 format logging with "[INFO]" prefix
func Infof(format string, msg ...interface{}) {
	if info {
		printFormat("[INFO]  "+format, msg...)
	}
}

//Warning any level logging with "[WARN]" prefix
func Warning(msg ...interface{}) {
	if warn {
		println("[WARN] ", msg...)
	}
}

//Warningf any level format logging with "[WARN]" prefix
func Warningf(format string, msg ...interface{}) {
	if warn {
		printFormat("[WARN]  "+format, msg...)
	}
}
