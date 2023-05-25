package config

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//for now it just opens the database base on filename pass to it,
// but it can be modified to open database based on environment later :)
/*
SetupDatabase initializes and sets up a connection to a SQLite database.
returns the created database connection and an error if any occurs.
*/
func SetupDatabase(filename string) (*sql.DB, error) {
	fmt.Println("Setting up database...")
	var err error
	DB, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, errors.New("Error opening database connection %s" + err.Error())
	}
	return DB, nil
}
