package config

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

/*
SetupDatabase initializes and sets up a connection to a SQLite database.
It takes a filename as a parameter representing the path to the SQLite database file.
The function opens the database connection using the "sqlite3" driver and the provided filename.
The connection is assigned to the package-level DB variable.
If an error occurs during the connection setup, it is printed to the console.
Finally, the function returns the created database connection.
*/
func SetupDatabase(filename string) *sql.DB {
	fmt.Println("Setting up database...")
	var err error
	DB, err = sql.Open("sqlite3", filename) // Assign to package-level DB variable
	if err != nil {
		fmt.Println("Status:", err)
	}
	return DB
}
