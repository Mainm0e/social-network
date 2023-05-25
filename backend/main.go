package main

import (
	"backend/db"
	"errors"
)

func main() {
	err := db.StartDB()
	if err != nil {
		panic(errors.New("Error starting database: " + err.Error()))
	}
}
