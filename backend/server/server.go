package server

import (
	"log"
	"os"
	"time"
)

const (
	// The HTTP port that the server will listen on
	HTTP_PORT = ":8080"
	// The HTTPS port that the server will listen on
	HTTPS_PORT = ":8443"
	// The path to the TLS certificate
	TLS_CERT_PATH = "" // TODO
	// The path to the TLS key
	TLS_KEY_PATH = "" // TODO
	// Pacth to log files
	LOG_PATH = "./backend/logs/"
)

/*
initiateLogging creates a log file with each instance of server startup, and sets
the output of the log package to the log file. All log messages will be written to
the log file which allows for easier debugging and a less cluttered terminal.
*/
func initiateLogging() {
	// Create a logfile with the name "log_YYYMMDD_HHMMSS.log" in the :/backend/logs directory
	logFile, err := os.OpenFile(LOG_PATH+"log_"+time.Now().Format("20060102_150405")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	// Do not defer close as this will prevent the log file from being written to!

	// Set the output of the log package to the log file
	log.SetOutput(logFile)
}
