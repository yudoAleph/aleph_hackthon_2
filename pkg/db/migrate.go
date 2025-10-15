package db

import (
	"log"
	"user-service/internal/app/migrations"

	"gorm.io/gorm"
)

// RunMigrations performs database migrations using the migration system
func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Get the underlying SQL DB from GORM
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Create migration runner
	runner := migrations.NewRunner(sqlDB)

	// Run migrations
	if err := runner.MigrateUp(); err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
