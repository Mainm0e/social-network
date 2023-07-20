package main

import (
	"backend/server"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	err := server.AaaawwwwwSheeeetttttItsAboutToGoDown(basePath, "http", server.LOG_PATH)
	if err != nil {
		fmt.Printf("Error: %s", err) // Keep this line during development, for debugging using terminal
		log.Printf("Error: %s", err)
	}
}
