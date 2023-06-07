package main

import (
	"backend/database"
	"errors"
)

func main() {
	err := database.Check("./backend/database/database.db", "./backend/database/test.sql")
	if err != nil {
		panic(errors.New("Error starting database: " + err.Error()))
	}
}
