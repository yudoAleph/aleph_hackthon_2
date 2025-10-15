package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"user-service/internal/app/models"
	"user-service/internal/app/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB holds test database connection and utilities
type TestDB struct {
	DB    *gorm.DB
	SqlDB *sql.DB
	Mock  sqlmock.Sqlmock
}

// SetupTestDB creates a test database connection
func SetupTestDB() (*TestDB, error) {
	// Use SQLite for testing
	dsn := "file::memory:?cache=shared"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	return &TestDB{
		DB:    db,
		SqlDB: sqlDB,
	}, nil
}

// SetupTestDBWithMock creates a test database with sqlmock
func SetupTestDBWithMock() (*TestDB, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlmock: %w", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm db: %w", err)
	}

	return &TestDB{
		DB:    gormDB,
		SqlDB: sqlDB,
		Mock:  mock,
	}, nil
}

// Close closes the test database connection
func (tdb *TestDB) Close() error {
	if tdb.SqlDB != nil {
		return tdb.SqlDB.Close()
	}
	return nil
}

// MigrateTestDB runs migrations on test database
func (tdb *TestDB) MigrateTestDB() error {
	// Auto-migrate the schema
	err := tdb.DB.AutoMigrate(&models.User{}, &models.Contact{})
	if err != nil {
		return fmt.Errorf("failed to migrate test database: %w", err)
	}
	return nil
}

// TestUser creates a test user for testing
func TestUser() *models.User {
	return &models.User{
		FullName: "Test User",
		Email:    "test@example.com",
		Phone:    "+1234567890",
		Password: "hashedpassword",
	}
}

// TestContact creates a test contact for testing
func TestContact(userID uint) *models.Contact {
	email := "contact@example.com"
	return &models.Contact{
		UserID:   userID,
		FullName: "Test Contact",
		Phone:    "+0987654321",
		Email:    &email,
		Favorite: false,
	}
}

// CreateTestUser creates a test user in the database
func CreateTestUser(ctx context.Context, repo repository.Repository) (*models.User, error) {
	user := TestUser()
	return repo.CreateUser(ctx, user)
}

// CreateTestContact creates a test contact in the database
func CreateTestContact(ctx context.Context, repo repository.Repository, userID uint) (*models.Contact, error) {
	contact := TestContact(userID)
	return repo.CreateContact(ctx, contact)
}

// SetupTestEnvironment sets up the complete test environment
func SetupTestEnvironment(t *testing.T) (*TestDB, repository.Repository, func()) {
	t.Helper()

	// Setup test database
	testDB, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Migrate database
	if err := testDB.MigrateTestDB(); err != nil {
		testDB.Close()
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Create repository
	repo := repository.NewRepository(testDB.DB)

	// Return cleanup function
	cleanup := func() {
		if err := testDB.Close(); err != nil {
			log.Printf("Error closing test database: %v", err)
		}
	}

	return testDB, repo, cleanup
}

// SetupTestEnvironmentWithMock sets up test environment with mocked database
func SetupTestEnvironmentWithMock(t *testing.T) (*TestDB, repository.Repository, func()) {
	t.Helper()

	// Setup test database with mock
	testDB, err := SetupTestDBWithMock()
	if err != nil {
		t.Fatalf("Failed to setup test database with mock: %v", err)
	}

	// Create repository
	repo := repository.NewRepository(testDB.DB)

	// Return cleanup function
	cleanup := func() {
		if err := testDB.Close(); err != nil {
			log.Printf("Error closing test database: %v", err)
		}
	}

	return testDB, repo, cleanup
}

// GetTestJWTToken returns a test JWT token for testing
func GetTestJWTToken() string {
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.test_signature"
}

// GetTestJWTSecret returns a test JWT secret
func GetTestJWTSecret() string {
	return "test_jwt_secret_key"
}
