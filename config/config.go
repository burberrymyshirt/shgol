package config

import (
	"os"
)

// DatabaseCredentials is a struct that holds the credentials to connect to a database.
// This is just a convenient way to transfer the credentials to the database connection function.
type DatabaseCredentials struct {
	Username     string
	Password     string
	Host         string
	Port         string
	Database     string
	DatabaseType string
}

func GetDBConfig() DatabaseCredentials {
	dbuser := os.Getenv("DB_USERNAME")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")
	databaseType := os.Getenv("DB_TYPE")

	return DatabaseCredentials{
		Username:     dbuser,
		Password:     dbpassword,
		Host:         dbhost,
		Port:         dbport,
		Database:     database,
		DatabaseType: databaseType,
	}
}
