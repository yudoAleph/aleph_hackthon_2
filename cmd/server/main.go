package main

import (
	"log"
	"user-service/configs"
	"user-service/internal/app/handlers"
	"user-service/internal/app/repository"
	"user-service/internal/app/routes"
	"user-service/internal/app/service"
	"user-service/pkg/db"

	"github.com/gin-gonic/gin"
)

// @title Contact Management API
// @version 1.0
// @description This is a contact management server.
// @termsOfService http://swagger.io/terms/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg := configs.LoadConfig()

	// Initialize DB
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Run migrations
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Initialize repository
	repo := repository.NewRepository(database)

	// Initialize service
	svc := service.NewService(repo, cfg.JWTSecret)

	// Initialize handler
	handler := handlers.NewHandler(svc, cfg.JWTSecret)

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Initialize Gin router
	router := gin.New()

	// Configure routes
	routes.SetupRoutes(router, handler, cfg.JWTSecret)

	// Start server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
