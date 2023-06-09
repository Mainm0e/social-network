package server

import (
	"backend/db"
	"errors"
	"log"
	"net/http"
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
	// Path to log files
	LOG_PATH = "./backend/logs/"
)

/*
initiateLogging creates a log file with each instance of server startup, and sets
the output of the log package to the log file. All log messages will be written to
the log file which allows for easier debugging and a less cluttered terminal.
*/
func initiateLogging() error {
	// Create a logfile with the name "log_YYYMMDD_HHMMSS.log" in the :/backend/logs directory
	logFile, err := os.OpenFile(LOG_PATH+"log_"+time.Now().Format("20060102_150405")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return errors.New("Error opening log file: " + err.Error())
	}
	// Do not defer close as this will prevent the log file from being written to!

	// Set the output of the log package to the log file
	log.SetOutput(logFile)

	return nil
}

/*
initialiseRoutes creates a new ServeMux and registers handler functions for various routes.
The ServeMux is returned and used by the StartServer() function to register the handler
functions for the HTTP/S routes.
*/
func initialiseRoutes() *http.ServeMux {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register handler functions for various routes
	// TODO: fix "handlers" package
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/main", handlers.Main)

	// Return the mux
	return mux
}

func setupHTTP() error {
	// TODO
	return nil
}

func setupHTTPS() error {
	// TODO
	return nil
}

/*
StartServer starts a server instance on a port number using the input protocol specified.
The server package includes predefined constants for the HTTP and HTTPS ports, as well as
the TLS certificate and key paths. The server will initialise a websocket manager, register
websocket event handlers, and start the manager in a separate goroutine, whilst also registering
handlers for the HTTP/S routes. Finally the server will start listening for requests on the
specified port, with routes defined by the initiateRoutes() helper function.
*/
func StartServer(protocol string) error {
	// Initiate logging
	err := initiateLogging()
	if err != nil {
		return errors.New("StartServer() error: " + err.Error())
	}

	// Check input protocol (only HTTP and HTTPS are supported, no quantum entanglement yet)
	if protocol != "http" && protocol != "https" {
		return errors.New("StartServer() error: invalid protocol specified: " + protocol)
	}

	// Check / migrate database
	// TEMP: use first migration file as initial schema for now
	err = db.Check("./backend/db/database.db", "./backend/db/migrations/01_initial_schema.sql")
	if err != nil {
		return errors.New("StartServer() error: " + err.Error())
	}

	// If HTTP is specified, setup HTTP server
	if protocol == "http" {
		err = setupHTTP()
		if err != nil {
			return errors.New("StartServer() error: " + err.Error())
		}
	}
	// If HTTPS is specified, setup HTTPS server
	if protocol == "https" {
		err = setupHTTPS()
		if err != nil {
			return errors.New("StartServer() error: " + err.Error())
		}
	}

	return nil
}
