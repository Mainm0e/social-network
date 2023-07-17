package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

/*
Initialise is called in the event that a database needs to be created. It takes the
file path for the desired database, as well as the file path for the sql database
creation file (both as strings), and executes the sql file, piping the queries
directly into the specified database file. An error value is returned, which is
non-nil in the event errors are encountered in opening / creating the database file
or the sql file.
*/
// func initialise(databasePathAndName, sqlFilePathAndName string) error {
// 	// Initialise specified database
// 	database, err := sql.Open("sqlite3", databasePathAndName)
// 	if err != nil {
// 		return err
// 	}
// 	defer database.Close()
// 	// Open sql database creation file
// 	file, err := os.Open(sqlFilePathAndName)
// 	if os.IsNotExist(err) {
// 		return errors.New("database sql creation file ( " +
// 			sqlFilePathAndName + " ) not found")
// 	} else {
// 		// Read the file into a buffer
// 		buf := make([]byte, 1024)
// 		var str string
// 		for {
// 			n, err := file.Read(buf)
// 			if err != nil {
// 				break
// 			}
// 			str += string(buf[:n])
// 		}

// 		// Execute the SQL
// 		_, err = database.Exec(str)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

/*
executeMigration is a function used to migrate the database to the latest version.
It takes in a database connection and applies all the migrations in the
./backend/db/migrations folder. It returns a non-nil error value if an error is
encountered in applying the migrations.
*/
func executeMigration(DB *sql.DB) error {
	fmt.Println("migrating...")
	migrations := &migrate.FileMigrationSource{
		Dir: "./db/migrations",
	}
	fmt.Println("migration:", migrations)
	n, err := migrate.Exec(DB, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return errors.New("error applying migrations: " + err.Error())
	}
	fmt.Println("applied migrations!\n", n)
	return nil
}

/*
openConnection is a helper function that opens a connection to the specified database.
It takes the database name as a string and returns an error value, which is non-nil
in the event that an error is encountered in opening the connection.
*/
func openConnection(dataSourceName string) error {
	var err error
	// Connect global database variable to specified database
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Printf("error while opening database: %v", err)
		return err
	}
	// Set max open connections to 1, although this may be redundant as the database
	// by design can only have one write operation at a time
	DB.SetMaxOpenConns(1)
	return nil
	/* REMEMBER TO SETUP SHUTDOWN HOOK IN SERVER CODE!!! */
	// Start a goroutine to listen for a shutdown signal
	// go func() {
	// 	// Create a channel to receive OS signals
	// 	signals := make(chan os.Signal, 1)
	// 	// Notify this channel when we receive a SIGINT (Ctrl+C)
	// 	// or SIGTERM (Docker stop, Kubernetes, etc.)
	// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	// 	// Block until we receive a signal
	// 	sig := <-signals
	// 	log.Printf("Received shutdown signal: %s", sig)
	// 	// Perform shutdown tasks here.
	// 	// For example, close database connections.
	// 	if err := database.DB.Close(); err != nil {
	// 		log.Fatalf("Could not close database: %v", err)
	// 	}
	// 	log.Println("Shutdown complete. Exiting.")
	// 	// Exit the application
	// 	os.Exit(0)
	// }()
}

/*
Check is a global function that takes a database file path and name and the sql
database creation file path and name as string inputs. It first checks if the database
file exists and if not, initialises it by calling the local initialise function. It
then calls the local openConnection function to open a connection to the database and
finally calls the local executeMigration function to migrate the database to the latest
version. It returns a non-nil error value if an error is encountered in any of these
steps.
*/
func Check(dbFile, sqlFile string) error {
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		// Initialise database if specified input does not already exist
		fmt.Print(Colour.Yellow + "\ndatabase not found, initialising...\n\n" + Colour.Reset)
		/* err = initialise(dbFile, sqlFile)
		if err != nil {
			return err
		} */
	}

	// Initialize the database connection, which is used throughout the application's lifetime
	err = openConnection(dbFile)
	if err != nil {
		return err
	}

	err = executeMigration(DB)
	if err != nil {
		return (errors.New("error migrating database: " + err.Error()))
	}
	return nil
}

/*
execInsertQuery is a helper function that executes an insert query and returns the
resulting id. It takes a query string and an optional variadic number of arguments,
and returns the resulting id and an error value, which is non-nil in the event that
an error is encountered in executing the query. It is called by the global
InsertData() function.
*/
func executeInsertionQuery(query string, args ...any) (int64, error) {
	// Prepare statement for inserting data
	statement, err := DB.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer statement.Close()

	// Execute statement
	result, err := statement.Exec(args...)
	if err != nil {
		return -1, err
	}

	// Get id of inserted row
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

/*
CheckDataDoesNotExist checks if the data already exists in the table. It takes the
table name, column name and value as strings, and returns an error value, which is
non-nil in the event that the data already exists in the table. It is called by the
global InsertData() function.
*/
func CheckDataDoesNotExist(table, column string, value any) error {

	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE %s = ?", table, column)

	// QueryRow executes a query that is expected to return at most one row.
	row := DB.QueryRow(query, value)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count != 0 {
		return fmt.Errorf("%v already exists in %s", value, table)
	}

	return nil
}

/*
InsertData is a global function that handles data insertion into the specified
table within the database.

Parameters:
- tableName: The name of the database table where the data will be inserted.
- args: A variadic parameter that represents the data to be inserted. The type
and order of data should align with the database schema for the specific table.

Return Values:
- id: Returns the ID of the last inserted row if the operation was successful.
- error: Returns a non-nil error in case of any issues during the data insertion
process.

The function leverages the InsertRules map to validate the uniqueness of the data
within the specified table (using the CheckDataDoesNotExist function). It ensures
that the data to be inserted does not violate any uniqueness constraints. After
the uniqueness check, the function invokes execInsertQuery to prepare and execute
the SQL insert statement. The function returns an error if the data insertion
failed, otherwise, it provides the ID of the newly inserted data row.

Please note:
  - If the tableName is not defined in the InsertRules map, an error will be returned.
    (i.e. make sure database_variables.go is updated with the new table name!)
  - The args parameter should strictly follow the order and data type expectations of
    the corresponding table schema.
*/
func InsertData(tableName string, args ...any) (int64, error) {
	rule, ok := InsertRules[tableName]
	if !ok {
		return -1, fmt.Errorf("unknown table name: %s", tableName)
	}

	if rule.ExistTable != "" && rule.ExistField != "" {
		err := CheckDataDoesNotExist(rule.ExistTable, rule.ExistField, args[0])
		if err != nil {
			return -1, fmt.Errorf(rule.ExistError)
		}
	}

	for i := range rule.NotExistTables {
		if tableName == "posts" || tableName == "notifications" && rule.NotExistFields[i] == "groupId" && args[i] == 0 {
			continue
		}
		err := CheckDataDoesNotExist(rule.NotExistTables[i], rule.NotExistFields[i], args[i])
		if err == nil {
			return -1, fmt.Errorf(rule.NotExistErrors[i])
		}
	}

	return executeInsertionQuery(rule.Query, args...)
}

/*
DeleteData is a function that deletes a specific row from a specified table
in the database. It takes as input the table name as a string and the value
of the primary key for the row to be deleted. The function first checks if the
provided table name exists in the predefined map TableKeys, and retrieves the
primary key associated with that table. It then checks if the keyValue exists
in the table. If it doesn't exist, it returns an error. If it exists, it prepares
and executes an SQL DELETE query to remove the row. The function returns an
error if any issues are encountered during these operations (for example, if
the table does not exist, or if the keyValue does not exist). If no rows are
affected by the DELETE operation, it returns an "item not found" error. If the
operation is successful, it returns nil.
*/
func DeleteData(tableName string, keyValues ...any) error {
	// Table validity check
	key, ok := TableKeys[tableName]
	if !ok {
		return errors.New("table does not exist")
	}

	// keyValue validity check
	for i, keyValue := range keyValues {
		if err := CheckDataDoesNotExist(tableName, key[i], keyValue); err == nil {
			return errors.New("data does not exist")
		}
	}
	var myQuery string
	if len(keyValues) == 1 {
		myQuery = fmt.Sprintf("DELETE FROM %s WHERE %s = ?", tableName, key[0])
	} else {
		myQuery = fmt.Sprintf("DELETE FROM %s WHERE %s = ? AND %s = ?", tableName, key[0], key[1])
	}
	statement, err := DB.Prepare(myQuery)
	if err != nil {
		return err
	}
	defer statement.Close()

	result, err := statement.Exec(keyValues...)
	if err != nil {
		return err
	}

	// Check if any rows were affected by the DELETE operation
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return errors.New("item not found")
	}
	return nil
}

/*
UpdateData is a global function that updates data in the database.
It takes as input a table name and a variadic number of arguments, which represent
the new data to be updated. The table name should correspond to one of the keys in
the UpdateRules map. The associated value in UpdateRules is a SQL UPDATE query,
which specifies how to update the data in the respective table. The function
constructs the final SQL query by replacing the placeholders in the UPDATE query
with the values from the arguments. The last argument should always be the primary
key of the row to be updated. The function returns an error value, which is non-nil
if an error was encountered during the update process.
*/
func UpdateData(tableName string, args ...any) error {
	myQuery, ok := UpdateRules[tableName]
	if !ok {
		return fmt.Errorf("unknown table name: %s", tableName)
	}

	statement, err := DB.Prepare(myQuery)
	if err != nil {
		return err
	}
	defer statement.Close()

	result, err := statement.Exec(args...)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return errors.New("no rows were updated")
	}

	return nil
}

/*
FetchData fetches data from a specified table in the database according to provided
conditions and arguments. The function operates dynamically, using the FetchRules
map to determine the fields to select and how to scan them.

Parameters:
  - table: The name of the table in the database from which data should be fetched.
  - condition: The WHERE clause in the SQL query to filter the data in the table.
    This is a string which can include placeholders for query parameters.
  - args: Any number of arguments which will be passed to the SQL query. These
    replace the placeholders in the 'condition' parameter.

Returns:
  - A slice of interface{}, where each item in the slice is a struct representing
    a row in the table.
  - An error, which will be non-nil if there was an error executing the query or
    scanning the result set.

Usage:
FetchData("users", "age > ?", 18) will fetch all users older than 18.

Note: This function does not sanitize inputs. Always ensure that input is sanitized
to prevent SQL injection attacks. Additionally, this function does not handle
potential database errors like unavailability, so use it with proper error handling.
*/
func FetchData(table string, condition string, args ...any) ([]any, error) {
	// Check if the table exists in the FetchRules struct (databse_variables.go)
	if _, ok := FetchRules[table]; !ok {
		return nil, fmt.Errorf("unknown table: %s", table)
	}

	// Prepare the SQL query
	var query string
	if condition == "" {
		query = fmt.Sprintf("SELECT %s FROM %s", FetchRules[table].SelectFields, table)
	} else {
		// Count the number of placeholders needed
		placeholders := strings.Count(condition, "?")
		if placeholders != len(args) {
			return nil, fmt.Errorf("the number of args does not match the number of placeholders in the condition")
		}
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s", FetchRules[table].SelectFields, table, condition)
	}

	// Execute the query
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Fetch rows
	result := make([]any, 0)
	for rows.Next() {
		item, err := FetchRules[table].ScanFields(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// DEPRECATED CODE BELOW
