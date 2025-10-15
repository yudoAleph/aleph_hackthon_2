# Database Migrations

This project uses a custom migration system to manage database schema changes. The migration system provides version control for your database schema and allows you to apply, rollback, and track database changes safely.

## Migration Files

Migration files are located in `internal/app/migrations/`:

- `migrations.go` - Contains all migration definitions
- `runner.go` - Contains the migration runner logic

## Available Commands

### Using Make Commands

```bash
# Apply all pending migrations
make migrate-up

# Rollback the last migration
make migrate-down

# Check migration status
make migrate-status
```

### Using Direct Commands

```bash
# Apply all pending migrations
go run ./cmd/migrate/main.go -command=up

# Rollback the last migration
go run ./cmd/migrate/main.go -command=down

# Check migration status
go run ./cmd/migrate/main.go -command=status
```

## How It Works

### Migration Structure

Each migration has:
- `ID`: Unique identifier for the migration
- `Up`: Function to apply the migration
- `Down`: Function to rollback the migration

### Migration Tracking

The system uses a `schema_migrations` table to track which migrations have been applied:

```sql
CREATE TABLE schema_migrations (
    id VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Current Migrations

1. **001_create_users_table** - Creates the users table with all necessary columns and indexes
2. **002_create_contacts_table** - Creates the contacts table with foreign key relationship to users

## Adding New Migrations

To add a new migration:

1. Add a new migration entry to the `GetMigrations()` function in `internal/app/migrations/migrations.go`
2. Follow the naming convention: `XXX_description` where XXX is a zero-padded number
3. Implement both `Up` and `Down` functions
4. Test the migration with `make migrate-up` and `make migrate-down`

Example:

```go
{
    ID: "003_add_user_status",
    Up: func(tx *sql.Tx) error {
        _, err := tx.Exec(`ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active'`)
        return err
    },
    Down: func(tx *sql.Tx) error {
        _, err := tx.Exec(`ALTER TABLE users DROP COLUMN status`)
        return err
    },
}
```

## Best Practices

1. **Test Migrations**: Always test both up and down migrations
2. **Idempotent**: Migrations should be safe to run multiple times
3. **Transactional**: Each migration runs in a transaction for atomicity
4. **Version Control**: Commit migration files with your code changes
5. **Backup**: Always backup your database before running migrations in production

## Troubleshooting

### Migration Fails

If a migration fails:
1. Check the error message
2. Fix the issue in the migration code
3. The migration system will retry from the failed migration

### Foreign Key Issues

Ensure data types match between referenced and referencing columns:
- Use `BIGINT UNSIGNED` for foreign keys referencing GORM's default `uint` primary keys

### Rollback Issues

If rollback fails, you may need to manually fix the database state before retrying.