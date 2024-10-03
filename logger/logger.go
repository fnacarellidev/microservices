package logger

import (
	"log"
	"os"
)

const flags = log.Lshortfile | log.LstdFlags

func logMessage(level string, output *os.File, msg string, v ...any) {
	l := log.New(output, level, flags)

	_, testing := os.LookupEnv("GO_TESTING")

	if !testing {
		l.Printf(msg, v...)
	}
}

var InfoLog = func(msg string, v ...any) {
	logMessage("INFO: ", os.Stdout, msg, v...)
}

var ErrorLog = func(msg string, v ...any) {
	logMessage("ERROR: ", os.Stderr, msg, v...)
}

var WarningLog = func(msg string, v ...any) {
	logMessage("WARNING: ", os.Stdout, msg, v...)
}
