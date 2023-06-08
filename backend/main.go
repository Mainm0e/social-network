package main

import (
	"backend/db"
	"errors"
)

func main() {
	err := db.Check("./backend/database/database.db", "./backend/database/test.sql")
	if err != nil {
		panic(errors.New("Error starting database: " + err.Error()))
	}
}
