package routes

import (
	"time"
	"user-service/internal/app/handlers"
	"user-service/internal/logger"
	"user-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, h *handlers.Handler, jwtSecretKey string) {
	// Add middlewares
	router.Use(middleware.SecureHeaders())
	router.Use(middleware.TimeoutMiddleware(30 * time.Second)) // 30 second timeout
	router.Use(logger.JSONLogMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", h.Register)
		public.POST("/auth/login", h.Login)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(jwtSecretKey))
	{
		// User routes
		protected.GET("/me", h.GetProfile)
		protected.PUT("/me", h.UpdateProfile)

		// Contact routes
		contacts := protected.Group("/contacts")
		{
			contacts.GET("", h.ListContacts)
			contacts.POST("", h.CreateContact)
			contacts.GET("/:id", h.GetContact)
			contacts.PUT("/:id", h.UpdateContact)
			contacts.DELETE("/:id", h.DeleteContact)
		}
	}
}
