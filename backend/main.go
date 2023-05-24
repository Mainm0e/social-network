package main

import (
	"backend/config"
)

func main() {
	config.SetupDatabase("db/database.db")
	config.Migrate()
}
