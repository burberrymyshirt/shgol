package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/burberrymyshirt/shurl/config"
	"github.com/burberrymyshirt/shurl/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var con *gorm.DB

// DBConnection fetches the gorm DB from anywhere, regardless of which database is being used
func DBConnection() *gorm.DB {
	return con
}

// DatabaseConnection is done for ease of changing DB. This works with MariaDB and SQL Server.
// This  can however be extended manually, with any of the databases supported by Gorm (see https://gorm.io/docs/connecting_to_the_database.html)
type DatabaseConnection interface {
	DatabaseInit()
	DatabaseMigrator()
}

// DatabaseConnectionFactory returns a database that implements the DatabaseConnection interface from the RDBMS written in the "DB_TYPE" env
func DatabaseConnectionFactory() DatabaseConnection {
	switch dbType := strings.ToUpper(os.Getenv("DB_TYPE")); dbType {
	case "MYSQL", "MARIADB":
		return mariadbConnection{}
	case "SQLITE":
		return sqliteConnection{}
	// case "SQLSERVER":
	// 	return sqlServerConnection{}
	default:
		log.Fatalf("unsupported RDBMS provided")
	}
	return nil // the default case will call os.Exit(1), so this is technically unreachable, but the compiler compolains anyway.
}

// mariadbConnection implements the DatabaseConnection interface
type mariadbConnection struct{}

func (mariadbConnection) DatabaseInit() {
	dbCredentials := config.GetDBConfig()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true&interpolateParams=true",
		dbCredentials.Username,
		dbCredentials.Password,
		dbCredentials.Host,
		dbCredentials.Port,
		dbCredentials.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %s", err.Error())
	}
	con = db
}

func (mariadbConnection) DatabaseMigrator() {
	con.AutoMigrate(model.Url{})
}

type sqliteConnection struct{}

func (sqliteConnection) DatabaseInit() {
	dbCredentials := config.GetDBConfig()

	db, err := gorm.Open(sqlite.Open(dbCredentials.Database), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %s", err.Error())
	}
	con = db
}

func (sqliteConnection) DatabaseMigrator() {
}
