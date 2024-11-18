package logger

import (
	"os"

	"github.com/mitchellh/colorstring"
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
	colorstring.Printf("["+startColor+"]"+level+format+"\n", args...)
}

func Debug(format string, args ...interface{}) {
	if debugMode {
		log(debugLabel, format, "blue", args...)
	}
}

func Warn(format string, args ...interface{}) {
	log(warnLabel, format, "yellow", args...)
}

func Info(format string, args ...interface{}) {
	log(infoLabel, format, "green", args...)
}

func Error(format string, args ...interface{}) {
	log(errorLabel, format, "red", args...)
}

func Fatal(format string, args ...interface{}) {
	log(fatalLabel, format, "red", args...)
	os.Exit(1)
}
