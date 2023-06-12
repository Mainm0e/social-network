package server

import (
	"bufio"
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestInitiateLogging(t *testing.T) {
	// Make sure the directory exists
	if _, err := os.Stat(TEST_LOG_PATH); os.IsNotExist(err) {
		err = os.MkdirAll(TEST_LOG_PATH, 0755)
		if err != nil {
			t.Fatalf("Failed to create log directory: %v", err)
		}
	}

	// Call the initiateLogging function
	err := initiateLogging(TEST_LOG_PATH)
	if err != nil {
		t.Fatalf("initiateLogging() failed: %v", err)
	}

	// List all files in the LOG_PATH directory
	files, err := os.ReadDir(TEST_LOG_PATH)
	if err != nil {
		t.Fatalf("Error reading log directory: %v", err)
	}

	// Search for the log file created by initiateLogging
	var logFilePath string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "log_") {
			logFilePath = TEST_LOG_PATH + file.Name()
			break
		}
	}

	if logFilePath == "" {
		t.Fatal("Log file not created")
	}

	// Test logging by writing a log message
	logMessage := "Test log message"
	log.Print(logMessage)

	// Read the content of the log file
	logFile, err := os.Open(logFilePath)
	if err != nil {
		t.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	scanner := bufio.NewScanner(logFile)
	var logFileContent string
	for scanner.Scan() {
		logFileContent += scanner.Text()
	}

	// Check if the log message is present in the log file
	if !strings.Contains(logFileContent, logMessage) {
		t.Fatal("Log message not written to log file")
	}
}

func TestSetupHTTP(t *testing.T) {
	// Subtest for server setup
	t.Run("Server Setup", func(t *testing.T) {
		// Specify the address the server should listen on
		testPort := "localhost:8081"

		// Create a new ServeMux
		mux := http.NewServeMux()

		// Register a test handler function
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Test Passed"))
		})

		// Channel to receive server instance
		serverCh := make(chan *http.Server, 1)

		// Start setupHTTP in a separate goroutine
		go setupHTTP(serverCh, testPort)

		// Give the server time to start
		time.Sleep(1 * time.Second)

		// Subtest for making GET request
		t.Run("GET Request", func(t *testing.T) {
			// Making an HTTP GET request to the test endpoint
			response, err := http.Get("http://localhost:8081/test")
			if err != nil {
				t.Fatalf("Error making GET request: %v", err)
			}
			defer response.Body.Close()

			// Reading the response body
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			// Check if the response body matches the expected response
			if string(body) != "Test Passed" {
				t.Fatalf("Expected 'Test Passed', got '%s'", string(body))
			}
		})

		// Subtest for server shutdown
		t.Run("Server Shutdown", func(t *testing.T) {
			// Clean up - Shutdown the server
			srv := <-serverCh
			if err := srv.Shutdown(context.TODO()); err != nil {
				t.Fatalf("Error shutting down server: %v", err)
			}
		})
	})
}

func TestSetupHTTPS(t *testing.T) {
	// Subtest for server setup
	t.Run("Server Setup", func(t *testing.T) {
		// Specify the address the server should listen on
		testPort := "localhost:8444"

		// Create a new ServeMux
		mux := http.NewServeMux()

		// Register a test handler function
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Test Passed"))
		})

		// Channel to receive server instance
		serverCh := make(chan *http.Server, 1)

		// Start setupHTTPS in a separate goroutine
		go setupHTTPS(serverCh, testPort)

		// Give the server time to start
		time.Sleep(1 * time.Second)

		// Subtest for making GET request
		t.Run("GET Request", func(t *testing.T) {
			// Specify to skip verification of certificates (because we're using a self-signed cert)
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}

			// Making an HTTPS GET request to the test endpoint
			response, err := client.Get("https://localhost:8444/test")
			if err != nil {
				t.Fatalf("Error making GET request: %v", err)
			}
			defer response.Body.Close()

			// Reading the response body
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			// Check if the response body matches the expected response
			if string(body) != "Test Passed" {
				t.Fatalf("Expected 'Test Passed', got '%s'", string(body))
			}
		})

		// Subtest for server shutdown
		t.Run("Server Shutdown", func(t *testing.T) {
			// Clean up - Shutdown the server
			srv := <-serverCh
			if err := srv.Shutdown(context.TODO()); err != nil {
				t.Fatalf("Error shutting down server: %v", err)
			}
		})
	})
}
