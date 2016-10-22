package logger

import (
	"log"
	"os"
)

//Level of logging
var Level uint

//Output of logging functions
var Output *os.File

func println(title string, msg ...interface{}) {
	arg := append([]interface{}{title}, msg...)
	log.Println(arg...)
}

func printFormat(format string, msg ...interface{}) {
	log.Printf(format, msg...)
}

//Debug level > 1 logging with "[DEBUG]" prefix
func Debug(msg ...interface{}) {
	if Level > 1 {
		println("[DEBUG]", msg...)
	}
}

//Debugf level > 1 format logging with "[DEBUG]" prefix
func Debugf(format string, msg ...interface{}) {
	if Level > 1 {
		printFormat("[DEBUG] "+format, msg...)
	}
}

//Info level > 0 logging with "[INFO]" prefix
func Info(msg ...interface{}) {
	if Level > 0 {
		println("[INFO] ", msg...)
	}
}

//Infof level > 0 format logging with "[INFO]" prefix
func Infof(format string, msg ...interface{}) {
	if Level > 0 {
		printFormat("[INFO]  "+format, msg...)
	}
}

//Warning any level logging with "[WARN]" prefix
func Warning(msg ...interface{}) {
	println("[WARN] ", msg...)
}

//Warningf any level format logging with "[WARN]" prefix
func Warningf(format string, msg ...interface{}) {
	printFormat("[WARN]  "+format, msg...)
}
