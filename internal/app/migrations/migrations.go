package migrations

import (
	"database/sql"
)

// Migration represents a database migration
type Migration struct {
	ID   string
	Up   func(*sql.Tx) error
	Down func(*sql.Tx) error
}

// GetMigrations returns all available migrations
func GetMigrations() []Migration {
	return []Migration{
		{
			ID: "001_create_users_table",
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec(`
					CREATE TABLE IF NOT EXISTS users (
						id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
						full_name VARCHAR(255) NOT NULL,
						email VARCHAR(255) NOT NULL UNIQUE,
						phone VARCHAR(20) NOT NULL,
						password VARCHAR(255) NOT NULL,
						avatar_url VARCHAR(255) NULL,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

						-- Indexes
						INDEX idx_users_full_name (full_name),
						INDEX idx_users_email (email),
						INDEX idx_users_phone (phone),
						INDEX idx_users_created_at (created_at)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
				`)
				return err
			},
			Down: func(tx *sql.Tx) error {
				_, err := tx.Exec(`DROP TABLE IF EXISTS users`)
				return err
			},
		},
		{
			ID: "002_create_contacts_table",
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec(`
					CREATE TABLE IF NOT EXISTS contacts (
						id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
						user_id INT UNSIGNED NOT NULL,
						full_name VARCHAR(255) NOT NULL,
						phone VARCHAR(20) NOT NULL,
						email VARCHAR(255) NULL,
						favorite BOOLEAN DEFAULT FALSE,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

						-- Foreign key constraint
						CONSTRAINT fk_contacts_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

						-- Indexes
						INDEX idx_contacts_user_id (user_id),
						INDEX idx_contacts_full_name (full_name),
						INDEX idx_contacts_phone (phone),
						INDEX idx_contacts_email (email),
						INDEX idx_contacts_favorite (favorite),
						INDEX idx_contacts_created_at (created_at),

						-- Composite index for common queries
						INDEX idx_contacts_user_favorite (user_id, favorite),
						INDEX idx_contacts_user_created (user_id, created_at DESC)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
				`)
				return err
			},
			Down: func(tx *sql.Tx) error {
				_, err := tx.Exec(`DROP TABLE IF EXISTS contacts`)
				return err
			},
		},
		{
			ID: "003_fix_schema_migrations_table",
			Up: func(tx *sql.Tx) error {
				// Check if the table exists and what columns it has
				var tableExists bool
				err := tx.QueryRow("SHOW TABLES LIKE 'schema_migrations'").Scan(&tableExists)
				if err != nil && err.Error() != "sql: no rows in result set" {
					// Table doesn't exist, just create it
					_, err = tx.Exec(`
						CREATE TABLE IF NOT EXISTS schema_migrations (
							version VARCHAR(255) PRIMARY KEY,
							name VARCHAR(255) NOT NULL,
							applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
						) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
					`)
					return err
				}

				// Table exists, check its structure
				rows, err := tx.Query("DESCRIBE schema_migrations")
				if err != nil {
					return err
				}
				defer rows.Close()

				var columns []string
				for rows.Next() {
					var field, typ, null, key, def, extra string
					if err := rows.Scan(&field, &typ, &null, &key, &def, &extra); err != nil {
						return err
					}
					columns = append(columns, field)
				}

				// Check if it has the old structure (id column)
				hasIdColumn := false
				hasNameColumn := false
				hasVersionColumn := false
				for _, col := range columns {
					switch col {
					case "id":
						hasIdColumn = true
					case "name":
						hasNameColumn = true
					case "version":
						hasVersionColumn = true
					}
				}

				if hasVersionColumn && hasNameColumn {
					// Table already has correct structure, nothing to do
					return nil
				}

				// Rename existing table
				_, err = tx.Exec("ALTER TABLE schema_migrations RENAME TO schema_migrations_old")
				if err != nil {
					return err
				}

				// Create new table with correct structure
				_, err = tx.Exec(`
					CREATE TABLE schema_migrations (
						version VARCHAR(255) PRIMARY KEY,
						name VARCHAR(255) NOT NULL,
						applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
					) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
				`)
				if err != nil {
					return err
				}

				// Migrate data based on old table structure
				if hasIdColumn && hasNameColumn {
					// Old table has id and name columns
					_, err = tx.Exec(`
						INSERT INTO schema_migrations (version, name, applied_at)
						SELECT id, name, applied_at FROM schema_migrations_old
					`)
				} else if hasIdColumn {
					// Old table has only id column, use id as name too
					_, err = tx.Exec(`
						INSERT INTO schema_migrations (version, name, applied_at)
						SELECT id, id, applied_at FROM schema_migrations_old
					`)
				}
				// If neither, table was empty anyway

				if err != nil {
					return err
				}

				// Drop the old table
				_, err = tx.Exec("DROP TABLE schema_migrations_old")
				return err
			},
			Down: func(tx *sql.Tx) error {
				// This migration is not reversible as it's a fix
				return nil
			},
		},
	}
}

// CreateMigrationsTable creates the migrations tracking table
func CreateMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	return err
}

// IsMigrationApplied checks if a migration has been applied
func IsMigrationApplied(db *sql.DB, migrationID string) (bool, error) {
	var count int
	// Try with 'version' column first (new structure)
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", migrationID).Scan(&count)
	if err != nil {
		// If that fails, try with 'id' column (old structure)
		err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE id = ?", migrationID).Scan(&count)
		if err != nil {
			return false, err
		}
	}
	return count > 0, nil
}

// MarkMigrationApplied marks a migration as applied
func MarkMigrationApplied(tx *sql.Tx, migrationID string) error {
	// Try with 'version' column first (new structure)
	_, err := tx.Exec("INSERT INTO schema_migrations (version, name, applied_at) VALUES (?, ?, NOW())", migrationID, migrationID)
	if err != nil {
		// If that fails, try with 'id' column (old structure)
		_, err = tx.Exec("INSERT INTO schema_migrations (id, name, applied_at) VALUES (?, ?, NOW())", migrationID, migrationID)
	}
	return err
}

// MarkMigrationUnapplied removes a migration from the applied list
func MarkMigrationUnapplied(tx *sql.Tx, migrationID string) error {
	// Try with 'version' column first (new structure)
	_, err := tx.Exec("DELETE FROM schema_migrations WHERE version = ?", migrationID)
	if err != nil {
		// If that fails, try with 'id' column (old structure)
		_, err = tx.Exec("DELETE FROM schema_migrations WHERE id = ?", migrationID)
	}
	return err
}
