package app

import (
	"user-service/configs"
	"user-service/internal/app/handlers"
	"user-service/internal/app/repository"
	"user-service/internal/app/service"

	"gorm.io/gorm"
)

// NewHandler creates a new handler instance
func NewHandler(cfg configs.Config, db *gorm.DB) *handlers.Handler {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo, cfg.JWTSecret)
	return handlers.NewHandler(svc, cfg.JWTSecret)
}
