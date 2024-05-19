package database

import (
	"github.com/adaggerboy/genesis-academy-case-app/models/config"
)

var (
	databaseFabric *DatabaseFabric[IDatabase] = NewDatabaseFabric[IDatabase]()
	database       IDatabase
)

func GetDatabaseFabric() *DatabaseFabric[IDatabase] {
	return databaseFabric
}

func InitDatabase(conf config.DatabaseEndpointConfig) error {
	db, err := databaseFabric.NewDatabaseController(conf)
	if err != nil {
		return err
	}
	database = db
	return nil
}

func GetDatabase() (db IDatabase) {
	return database
}
