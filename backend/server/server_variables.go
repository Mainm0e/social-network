package server

import (
	"os"
	"time"
)

const (
	// The HTTP port that the server will listen on
	HTTP_PORT = ":8080" // TODO: change for external hosting
	// The HTTPS port that the server will listen on
	HTTPS_PORT = ":8443" // TODO: change for external hosting
	// The path to the TLS certificate
	TLS_CERT_PATH = "" // TODO
	// The path to the TLS key
	TLS_KEY_PATH = "" // TODO
	// Path to log files (from backend image root directory)
	LOG_PATH = "./logs/"
	// Path to test log files (from server subdirectory)
	TEST_LOG_PATH = "./_test_logs/"
	// Server shutdown timeout
	SHUTDOWN_TIMEOUT = 5 * time.Second
	// Origin of frontend server / app
	FRONTEND_ORIGIN = "http://localhost:3000" // TODO: change for external hosting
)

var (
	// Make logfile a global variable
	LogFile *os.File
)
