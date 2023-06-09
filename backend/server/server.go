package server

import (
	"backend/db"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	// TODO: fix "handlers" package, maybe make struct which can be looped over to register handlers?
	mux.HandleFunc("/login", handlers.LoginPage)
	mux.HandleFunc("/register", handlers.RegisterPage)
	mux.HandleFunc("/main", handlers.MainPage)

	// Return the mux
	return mux
}

/*
setupHTTP creates a new http.Server instance with the specified properties, and starts
the server with ListenAndServe. The server instance is sent through the serverCh channel
to be used by the StartServer() function. If an error occurs, it is logged. The function
is blocking and will run until the server is closed.
*/
func setupHTTP(mux *http.ServeMux, serverCh chan<- *http.Server) {
	// Create a new http.Server with properties
	srv := &http.Server{
		Addr:         HTTP_PORT,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		// IdleTimeout:  15 * time.Second,
	}

	// Send the server instance through the channel
	serverCh <- srv

	// Start the HTTP server with ListenAndServe (blocking)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("ListenAndServe() error: %s", err)
	}
}

/*
setupHTTPS creates a new HTTPS http.Server instance with the specified properties, and starts
the server with ListenAndServeTLS. The server instance is sent through the serverCh channel
to be used by the StartServer() function. If an error occurs, it is logged. The function is
blocking and will run until the server is closed.
*/
func setupHTTPS(mux *http.ServeMux, serverCh chan<- *http.Server) {
	// Create a new http.Server with properties
	srv := &http.Server{
		Addr:         HTTPS_PORT,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		// IdleTimeout:  15 * time.Second,
	}

	// Send the server instance through the channel
	serverCh <- srv

	// Start the HTTPS server with ListenAndServeTLS (blocking)
	if err := srv.ListenAndServeTLS(TLS_CERT_PATH, TLS_KEY_PATH); err != http.ErrServerClosed {
		log.Printf("ListenAndServeTLS() error: %s", err)
	}
}

/*
StartServer starts a server instance on a port number using the input protocol specified.
The server package includes predefined constants for the HTTP and HTTPS ports, as well as
the TLS certificate and key paths. The server will initialise a websocket manager, register
websocket event handlers, and start the manager in a separate goroutine, whilst also registering
handlers for the HTTP/S routes. Finally the server will start a relevant HTTP/S server instance
in a separate goroutine, and wait for a signal to shutdown the server. The server will then
gracefully shutdown and close all connections. The function returns an error, which is non-nil
if an error occurs at any point during the server setup.
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

	// Setup channel to receive the server instance, enabling graceful shutdown
	serverCh := make(chan *http.Server)

	// TODO: initialise websocket manager and associated event handlers

	// Initialise routes
	mux := initialiseRoutes()

	// If HTTP is specified, setup HTTP server in a goroutine
	if protocol == "http" {
		go setupHTTP(mux, serverCh)
	}

	// If HTTPS is specified, setup HTTPS server in a goroutine
	if protocol == "https" {
		go setupHTTPS(mux, serverCh)
	}

	// Receive the server instance from the channel (corresponding to code in setupHTTP() and setupHTTPS())
	srv := <-serverCh

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Reading from the stop channel, meaning it is blocking until a signal is received
	<-stop

	// Create a context to allow graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), SHUTDOWN_TIMEOUT)
	defer cancel()

	// Perform the graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		return errors.New("Graceful shutdown of server failed: " + err.Error())
	}

	log.Print("Server exited properly")

	return nil
}
