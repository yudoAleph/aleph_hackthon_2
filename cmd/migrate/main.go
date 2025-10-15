package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"user-service/configs"
	"user-service/internal/app/migrations"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var command string
	flag.StringVar(&command, "command", "up", "Migration command: up, down, status")
	flag.Parse()

	// Load configuration
	cfg := configs.LoadConfig()

	// Build MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// Initialize database connection
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Test connection
	if err := database.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create migration runner
	runner := migrations.NewRunner(database)

	// Execute command
	switch command {
	case "up":
		if err := runner.MigrateUp(); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
	case "down":
		if err := runner.MigrateDown(); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
	case "status":
		if err := runner.Status(); err != nil {
			log.Fatalf("Migration status failed: %v", err)
		}
	default:
		log.Fatalf("Unknown command: %s. Use: up, down, or status", command)
	}
}
