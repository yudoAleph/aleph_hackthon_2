package middleware

import (
	"os"

	"github.com/labstack/echo/v4"
)

func BasicAuthValidator(username, password string, c echo.Context) (bool, error) {
	authUsername := os.Getenv("BASIC_AUTH_USERNAME")
	authPassword := os.Getenv("BASIC_AUTH_PASSWORD")

	lang := c.Request().Header.Get("Accept-Language")
	if lang == "" {
		return false, nil
	}

	// Save downstream
	c.Set("lang", lang)
	c.Set("auth", c.Request().Header.Get(echo.HeaderAuthorization))

	return username == authUsername && password == authPassword, nil
}
