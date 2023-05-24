package config

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

func Migrate() {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations/",
	}

	n, err := migrate.Exec(DB.DB(), "sqlite3", migrations, migrate.Up)
	if err != nil {
		fmt.Println("Error occcured:", err)
	}

	fmt.Printf("Applied %d migrations!\n", n)
}
