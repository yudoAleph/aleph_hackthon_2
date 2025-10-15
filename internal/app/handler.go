package app

import (
	"net/http"

	"user-service/configs"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	db      *gorm.DB
	service *Service
}

func NewHandler(cfg configs.Config, db *gorm.DB) *Handler {
	repo := NewUserRepository(db)
	service := NewService(repo, cfg.JWTSecret)
	return &Handler{db: db, service: service}
}

func (h *Handler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}

func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")

	if len(id) == 0 {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	user, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "not found")
	}

	return c.JSON(http.StatusOK, user)
}
