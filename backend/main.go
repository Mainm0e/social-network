package main

import (
	"backend/server"
	"fmt"
	"log"
)

func main() {
	err := server.AaaawwwwwSheeeetttttItsAboutToGoDown("http", server.LOG_PATH)
	if err != nil {
		fmt.Printf("Error: %s", err) // Keep this line during development, for debugging using terminal
		log.Printf("Error: %s", err)
	}
}
