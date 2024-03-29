package server

import (
	"backend/db"
	"backend/events"
	"backend/handlers"
	"backend/server/sessions"
	"backend/sockets"
	"backend/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
logAndResetRequest() is a middleware helper function that logs the request body
and resets the request body so that it can be read again by later middleware or
handlers. This function takes a pointer to an http.Request and returns a pointer
to an http.Request. It is also used for debugging purposes.
*/
func logAndResetRequest(r *http.Request) *http.Request {
	// If this is a WebSocket handshake request, don't try to read the body
	if r.Header.Get("Upgrade") == "websocket" {
		log.Println("logAndResetRequest() -- WebSocket handshake")
		return r
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return r
	}
	// Always close the request body after reading it to free up resources
	defer r.Body.Close()

	var event events.Event

	// Decode json byte slice
	if err := json.Unmarshal(body, &event); err != nil {
		// If problem decoding, log the error, reset the request body and return
		log.Println("logAndResetRequest() error - Error decoding event:", err)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		return r
	}

	// Print the event type
	log.Printf("logAndResetRequest() -- Event: %v ; Payload: %v", event.Type, string(event.Payload))

	// Reset the request body
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return r
}

/*
initiateLogging creates a log file with each instance of server startup, and sets
the output of the log package to the log file. All log messages will be written to
the log file which allows for easier debugging and a less cluttered terminal.
*/
func initiateLogging(logPath string) error {
	// Check if the log directory exists
	_, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		// If it doesn't exist, create it
		if err := os.Mkdir(logPath, 0755); err != nil {
			return errors.New("Error creating log directory: " + err.Error())
		}
	}

	// Create a logfile with the name "log_YYYMMDD_HHMMSS.log" in the :/backend/logs directory
	LogFile, err = os.OpenFile(logPath+"log_"+time.Now().Format("20060102_150405")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return errors.New("Error opening log file: " + err.Error())
	}
	// Do not defer close as this will prevent the log file from being written to!

	// Set the output of the log package to the log file
	log.SetOutput(LogFile)

	return nil
}

/*
loggerMiddleware is a middleware function which logs the URL path of each request to the
server. It takes an input of an http.Handler and returns an http.Handler. It can be coupled
with various other middleware functions to create a middleware chain. This pattern is used
to allow for various logic steps to be chained prior to the end handler filling the request,
with each step being self-contained and responsible for either handling off the request to
the next link in the chain or finalising the response itself in the event that:
  - It is the final link in the chain
  - It encounters an error or any other condition that prevents further processing.

This pattern facilitates ease of maintenance should additional middleware functions be required.
*/
func loggerMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Print out entire request body for debugging purposes
		r = logAndResetRequest(r)

		handler.ServeHTTP(w, r)
	})
}

/*
corsMiddleware is a function which takes an http.Handler as an input and returns an http.Handler.
It is used to handle CORS (Cross-Origin Resource Sharing) requests. CORS is a mechanism that
uses additional HTTP headers to tell browsers to give a web application running at one origin,
access to selected resources from a different origin. A web application executes a cross-origin
HTTP request when it requests a resource that has a different origin (domain, protocol, or port)
from its own. CORS allows web applications to bypass a browser's same-origin policy, allowing
for the safe use of resources from multiple sources. corsMiddleware() can be coupled with various
other middleware functions to create a middleware chain implemented by the loggerMiddleware()
function.
*/
func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests with the "credentials" header set to "true"
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Allow requests from the specific origin of the frontend application
		w.Header().Set("Access-Control-Allow-Origin", FRONTEND_ORIGIN) // Change this to your frontend origin
		// Allow specific HTTP methods, which provides some protection against CSRF attacks
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		// Allow the Content-Type header, which is required to be sent with POST requests
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, credentials, withCredentials")
		// Set the Access-Control-Max-Age header to cache preflight request (600 seconds = 10 minutes)
		w.Header().Set("Access-Control-Max-Age", "600")

		// Handle preflight requests, which is another way of saying "handle OPTIONS requests"
		// OPTIONS requests are sent by the browser to check if the server will allow a request
		// with the specified method and headers. If the server responds with a 200 OK, the
		// browser will send the actual request. If the server responds with a 403 Forbidden,
		// the browser will not send the actual request.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Begin wrapping, call the next handler in the chain
		handler.ServeHTTP(w, r)
	})
}

/*
authenticationMiddleware is a middleware function which handles authentication logic
for each request to the server. It takes an input of an http.Handler and returns an
http.Handler, calling the sessions.CheckAuthentication() function to check if the user
is authenticated. It can be coupled with various other middleware functions to create a
middleware chain implemented by the loggerMiddleware() function.
*/
func authenticationMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If this is a WebSocket handshake request, don't try to read the body
		if r.Header.Get("Upgrade") == "websocket" {
			// TODO: maybe implement websocket authentication here (?)
			handler.ServeHTTP(w, r)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body", http.StatusBadRequest)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		// Always close the request body after reading it to free up resources
		defer r.Body.Close()

		// Create a new reader with the body for JSON decoding
		reader1 := io.NopCloser(bytes.NewReader(body))
		// Initialise  event struct
		var event events.Event

		err = json.NewDecoder(reader1).Decode(&event)
		if err != nil {
			log.Println("Error decoding JSON", http.StatusBadRequest)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Refer to events catalogue in handlers package (handlers_structs.go)
		if _, ok := handlers.Events[event.Type]; ok {
			// Check that event type is not login or register (these events do not require authentication)
			if event.Type != "login" && event.Type != "register" {
				// Check if the request contains a sessionID cookie
				cookie, err := r.Cookie(sessions.COOKIE_NAME)
				if err != nil {
					// Handle error: No sessionID cookie found.
					log.Printf("authenticationWiddleware() error - No sessionID cookie found: %v", err)
					http.Error(w, "Invalid session", http.StatusUnauthorized)
					return
				}

				// Validate the sessionID cookie
				isValid, err := sessions.CookieCheck(cookie)
				if !isValid || err != nil {
					// Handle error: Invalid sessionID cookie.
					log.Printf("authenticationWiddleware() error - Invalid sessionID cookie: %v", err)
					http.Error(w, "Invalid session", http.StatusUnauthorized)
					return
				}
			}
		}

		// Create a new reader with the body for the next handler
		r.Body = io.NopCloser(bytes.NewReader(body))
		handler.ServeHTTP(w, r)
	})
}

/*
initialiseRoutes creates a new http.ServeMux and registers handler functions for various
routes. It then wraps the mux with the CORS, authentication and logging middleware functions
and returns the wrapped mux. The wrapped mux is then passed to the http.Server instance
created by setupHTTP() or setupHTTPS() to be used as the server's handler.
*/
func initialiseRoutes() http.Handler {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a websocket manager
	wsManager := sockets.NewManager()

	// Start the manager's main loop in a separate goroutine.
	go wsManager.Run()

	// Register handler functions for various routes
	mux.HandleFunc("/api", handlers.HTTPEventRouter)
	mux.HandleFunc("/ws", wsManager.ServeWS)

	// Wrap the mux with the CORS middleware and return it
	// Although the return type is an http.Handler, it is actually a wrapped *mux.Router which
	// has been chained with the CORS middleware
	return loggerMiddleware(corsMiddleware(authenticationMiddleware(mux)))
}

/*
setupHTTP creates a new http.Server instance with the specified properties, and starts
the server with ListenAndServe. The server instance is sent through the serverCh channel
to be used by the StartServer() function. If an error occurs, it is logged. The function
is blocking and will run until the server is closed.
*/
func setupHTTP(serverCh chan<- *http.Server, portAddress string) {
	// Initialise routes with middlewares
	handler := initialiseRoutes()

	// Create a new http.Server with properties
	srv := &http.Server{
		Addr:         portAddress,
		Handler:      handler,
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
func setupHTTPS(serverCh chan<- *http.Server, portAddress string) {
	// Initialise routes with middlewares
	handler := initialiseRoutes()

	// Create a new http.Server with properties
	srv := &http.Server{
		Addr:         portAddress,
		Handler:      handler,
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
AaaawwwwwSheeeetttttItsAboutToGoDown is the the function to end all functions, a work of
pure alchemical wizardry. Using 3 pieces of Adi's hair, two spoonfuls of Steve's toejam, a small
piece of Maryams rubber lizard, seventeen of Rick's tears and the inner lining of Salam's bike
tyre, this function mixes it all in a cauldron of nightmares, and turns gold into iron, success
into calamity, and robs all who read its code of at least 3 years of their life. Use with caution.
*/
func AaaawwwwwSheeeetttttItsAboutToGoDown(protocol string, logPath string) error {

	/* 	OLD DESCRIPTION
	StartServer starts a server instance on a port number using the input protocol specified.
	The server package includes predefined constants for the HTTP and HTTPS ports, as well as
	the TLS certificate and key paths. The server will initialise a websocket manager, register
	websocket event handlers, and start the manager in a separate goroutine, whilst also registering
	handlers for the HTTP/S routes. Finally the server will start a relevant HTTP/S server instance
	in a separate goroutine, and wait for a signal to shutdown the server. The server will then
	gracefully shutdown and close all connections. The function returns an error, which is non-nil
	if an error occurs at any point during the server setup.
	*/
	// Initiate logging

	err := initiateLogging(logPath)
	if err != nil {
		return errors.New("StartServer() error: " + err.Error())
	}
	error := utils.InitiateImagesPath()
	if error != nil {
		return errors.New("StartServer() error: " + error.Error())
	}

	// Check input protocol (only HTTP and HTTPS are supported, no quantum entanglement yet)
	if protocol != "http" && protocol != "https" {
		return errors.New("StartServer() error: invalid protocol specified: " + protocol)
	}

	// Check / migrate database
	// TEMP: use first migration file as initial schema for now
	err = db.Check("db/database.db", "db/migrations")
	if err != nil {
		return errors.New("StartServer() error: " + err.Error())
	}

	// Setup channel to receive the server instance, enabling graceful shutdown
	serverCh := make(chan *http.Server)
	fmt.Printf("Server starting on portocol %s...\n", protocol)
	// If HTTP is specified, setup HTTP server in a goroutine
	if protocol == "http" {
		go setupHTTP(serverCh, HTTP_PORT)
	}

	// If HTTPS is specified, setup HTTPS server in a goroutine
	if protocol == "https" {
		go setupHTTPS(serverCh, HTTPS_PORT)
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

	fmt.Println("\nServer shutdown gracefully... like a rabid five-winged swan!") // Keep this during development, for debugging via terminal
	log.Print("Server exited properly")

	// Close log file (package level variable)
	if LogFile != nil {
		LogFile.Close()
	}

	return nil
}

/*********************** DEPRECATED CODE BELOW ************************/

/*
extractCookie is a function which takes an http.Request as an input and returns a pointer
to an http.Cookie and an error. It is used to extract the session cookie from the request
header. It returns an error which is non-nil if a cookie with the sessions.COOKIE_NAME is
not present in the request header.
*/
// func extractCookie(r *http.Request) (*http.Cookie, error) {
// 	cookie, err := r.Cookie(sessions.COOKIE_NAME)
// 	if err != nil {
// 		return nil, errors.New("error in server.extractCookie(): " + err.Error())
// 	}
// 	return cookie, nil
// }
