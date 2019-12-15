package status

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AUTH_InvalidCredentials(c echo.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusUnauthorized,
		"scope":   "auth",
		"error":   "invalid_credentials",
		"message": "Your credentials are invalid.",
	})
}
