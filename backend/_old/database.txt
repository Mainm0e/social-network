package config

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const DATABASE_NAME = "test"

/*
Initialise is called in the event that a database needs to be created. It takes the
file path for the desired database, as well as the file path for the sql database
creation file (both as strings), and executes the sql file, piping the queries
directly into the specified database file. An error value is returned, which is
non-nil in the event errors are encountered in opening / creating the database file
or the sql file.
*/
func initialise(databasePathAndName, sqlFilePathAndName string) error {
	// Initialise specified database
	database, err := sql.Open("sqlite3", databasePathAndName)
	if err != nil {
		return err
	}
	defer database.Close()
	// Open sql database creation file
	file, err := os.Open(sqlFilePathAndName)
	if os.IsNotExist(err) {
		return errors.New("database sql creation file ( " +
			sqlFilePathAndName + " ) not found")
	} else {
		// Read the file into a buffer
		buf := make([]byte, 1024)
		var str string
		for {
			n, err := file.Read(buf)
			if err != nil {
				break
			}
			str += string(buf[:n])
		}

		// Execute the SQL
		_, err = database.Exec(str)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
SetupDatabase initializes (if necessary) and sets up a connection to a
SQLite database. returning the created database connection and a non-nil
error if one is encountered.
**NOTE**
For now the function just opens the database bases on the filename passed
to it, but it can be modified to open a database based on the environment
later :)
*/
func SetupDatabase(filename string) (*sql.DB, error) {
	fmt.Println("Setting up database...")
	var err error

	// Check if database file exists
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		// If not, initialise it
		err = initialise(filename, "../db/"+DATABASE_NAME+".sql")
		if err != nil {
			return nil, errors.New("Error initialising database: " + err.Error())
		}
	}

	// Open / connect to database
	DB, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, errors.New("Error opening database connection %s" + err.Error())
	}
	return DB, nil
}
