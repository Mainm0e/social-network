package config

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

// probably modify this to use dbconfig.yml later
/*
Migrate function is used to migrate the database to the latest version.
It takes in a database connection and applies all the migrations in the
db/migrations folder. It returns an error if any occurs.
*/
func Migrate(DB *sql.DB) error {
	fmt.Println("Migrating...")
	migrations := &migrate.FileMigrationSource{
		Dir: "./db/migrations",
	}
	fmt.Println("migration:", migrations)
	n, err := migrate.Exec(DB, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return errors.New("Error applying migrations: " + err.Error())
	}
	fmt.Println("Applied migrations!\n", n)
	return nil
}
