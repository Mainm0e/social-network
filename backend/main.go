package main

import "social-network/backend/config"

func main() {
	config.DB = config.SetupDatabase()
	config.Migrate()
}
