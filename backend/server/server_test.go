package server

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"
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
