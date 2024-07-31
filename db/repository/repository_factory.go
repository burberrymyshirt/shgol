package repository

import (
	"log"
	"os"
	"strings"
)

func NewUrlRepository() UrlRepository {
	switch dbType := strings.ToUpper(os.Getenv("DB_TYPE")); dbType {
	case "MYSQL", "MARIADB", "SQLITE":
		return NewGormUrlRepository()
	// case "SQLITE":
	// 	return NewSqliteUrlRepository()
	default:
		log.Fatalf("unsupported RDBMS provided")
	}
	return nil
}
