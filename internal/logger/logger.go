package logger

import (
	"fmt"
	"os"
)

var debugMode = false

const (
	infoLabel  = ""
	errorLabel = "[ERROR] "
	debugLabel = "[DEBUG] "
	fatalLabel = "[FATAL] "
	warnLabel  = "[WARNING] "
)

func SetDebugMode(mode bool) {
	debugMode = mode
}

func IsDebugMode() bool {
	return debugMode
}

func log(level string, format string, startColor string, args ...interface{}) {
	fmt.Printf(startColor+level+format+Reset+"\n", args...)
}

func Debug(format string, args ...interface{}) {
	if debugMode {
		log(debugLabel, format, Blue, args...)
	}
}

func Warn(format string, args ...interface{}) {
	log(warnLabel, format, Yellow, args...)
}

func Info(format string, args ...interface{}) {
	log(infoLabel, format, Green, args...)
}

func Error(format string, args ...interface{}) {
	log(errorLabel, format, Red, args...)
}

func Fatal(format string, args ...interface{}) {
	log(fatalLabel, format, Red, args...)
	os.Exit(1)
}
