package db

import (
	"backend/config"
	"errors"
)

func StartDB() error {
	DB, err := config.SetupDatabase("db/database.db")
	if err != nil {
		return (errors.New("Error setting up database: " + err.Error()))
	}
	defer DB.Close()
	err = config.Migrate(DB)
	if err != nil {
		return (errors.New("Error migrating database: " + err.Error()))
	}
	return nil
}
