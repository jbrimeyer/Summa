package summa

import (
	"fmt"
	"log"
	"os"
)

type logType string

var logFile *os.File
var infoLog logType = "INFO"
var errLog logType = "ERROR"

// Printf prints a message to the log using the appropriate log type
func (l logType) Printf(format string, args ...interface{}) {
	log.Printf(fmt.Sprintf("%s: %s", l, format), args...)
}

// startLogging opens the Summa log file and sets it as the log output stream
func startLogging(filePath string) error {
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	log.SetOutput(logFile)

	return nil
}

// endLogging closes the Summa log file and sets Stderr as the new log output stream
func endLogging() {
	if logFile != nil {
		logFile.Close()
	}
	log.SetOutput(os.Stderr)
}
