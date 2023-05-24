package config

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

func Migrate() {
	fmt.Println("Migrating...")
	migrations := &migrate.FileMigrationSource{
		Dir: "./db/migrations",
	}
	fmt.Println("migration:", migrations)
	n, err := migrate.Exec(DB, "sqlite3", migrations, migrate.Up)
	if err != nil {
		fmt.Println("Error occcured:", err)
	}

	fmt.Println("Applied migrations!\n", n)
}
