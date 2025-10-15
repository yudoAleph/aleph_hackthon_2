package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

// Runner handles running database migrations
type Runner struct {
	db *sql.DB
}

// NewRunner creates a new migration runner
func NewRunner(db *sql.DB) *Runner {
	return &Runner{db: db}
}

// MigrateUp runs all pending migrations
func (r *Runner) MigrateUp() error {
	log.Println("Running database migrations...")

	// Create migrations table if it doesn't exist
	if err := CreateMigrationsTable(r.db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrations := GetMigrations()

	for _, migration := range migrations {
		applied, err := IsMigrationApplied(r.db, migration.ID)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", migration.ID, err)
		}

		if applied {
			log.Printf("Migration %s already applied, skipping", migration.ID)
			continue
		}

		log.Printf("Applying migration: %s", migration.ID)

		// Start transaction
		tx, err := r.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction for migration %s: %w", migration.ID, err)
		}

		// Run migration
		if err := migration.Up(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to run migration %s: %w", migration.ID, err)
		}

		// Mark as applied
		if err := MarkMigrationApplied(tx, migration.ID); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to mark migration %s as applied: %w", migration.ID, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.ID, err)
		}

		log.Printf("Successfully applied migration: %s", migration.ID)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// MigrateDown rolls back the last migration
func (r *Runner) MigrateDown() error {
	log.Println("Rolling back last migration...")

	migrations := GetMigrations()
	if len(migrations) == 0 {
		log.Println("No migrations to roll back")
		return nil
	}

	// Find the last applied migration
	var lastMigration *Migration
	for i := len(migrations) - 1; i >= 0; i-- {
		applied, err := IsMigrationApplied(r.db, migrations[i].ID)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", migrations[i].ID, err)
		}
		if applied {
			lastMigration = &migrations[i]
			break
		}
	}

	if lastMigration == nil {
		log.Println("No applied migrations to roll back")
		return nil
	}

	log.Printf("Rolling back migration: %s", lastMigration.ID)

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction for rollback %s: %w", lastMigration.ID, err)
	}

	// Run rollback
	if err := lastMigration.Down(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to rollback migration %s: %w", lastMigration.ID, err)
	}

	// Mark as unapplied
	if err := MarkMigrationUnapplied(tx, lastMigration.ID); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to mark migration %s as unapplied: %w", lastMigration.ID, err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback %s: %w", lastMigration.ID, err)
	}

	log.Printf("Successfully rolled back migration: %s", lastMigration.ID)
	return nil
}

// Status shows the current migration status
func (r *Runner) Status() error {
	log.Println("Migration Status:")

	migrations := GetMigrations()

	for _, migration := range migrations {
		applied, err := IsMigrationApplied(r.db, migration.ID)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", migration.ID, err)
		}

		status := "Pending"
		if applied {
			status = "Applied"
		}

		log.Printf("  %s: %s", migration.ID, status)
	}

	return nil
}
