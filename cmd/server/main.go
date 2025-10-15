package main

import (
	"log"
	"net/http"
	"user-service/configs"
	"user-service/docs"
	"user-service/internal/app"
	middlewareApps "user-service/internal/middleware"
	"user-service/pkg/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg := configs.LoadConfig()

	// Initialize DB
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	db.AutoMigrate(database)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Pass database to handlers
	h := app.NewHandler(cfg, database) // make sure this function accepts *gorm.DB
	// Public Routes
	e.GET("/ping", h.Ping)

	groupMobile := e.Group("/api/v1/mobile", middleware.BasicAuth(middlewareApps.BasicAuthValidator))
	groupMobile.POST("/users/:id", h.Get)

	// Swagger UI and raw OpenAPI JSON
	docs.SwaggerInfo.BasePath = "/"
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/openapi.json", func(c echo.Context) error { return c.JSONBlob(http.StatusOK, []byte(docs.SwaggerInfo.SwaggerTemplate)) })

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
