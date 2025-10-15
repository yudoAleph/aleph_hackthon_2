package db

import (
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	// Get DB name from environment variable, fallback to default if not set
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "aleph.user_db" // default fallback
	}

	database, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return nil, err
	}
	return database, nil
}
