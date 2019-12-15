package handlers

import (
	"fmt"
	"net/http"

	m "github.com/Techassi/gomark/internal/models"
	"github.com/Techassi/gomark/internal/server/status"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// AUTH_JWTError handles the redirect to the login page if no JWT token is
// present
func AUTH_JWTError(e error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/login")
}

// AUTH_JWTRegister handles the regsiter process of a new user
func AUTH_JWTRegister(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/login")
}

// AUTH_JWTLogin handles the user authentication via the DB to login the user
func AUTH_JWTLogin(c echo.Context) error {
	app := c.Get("app").(*m.App)
	u := m.User{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}

	valid := app.DB.ValidCredentials(u)
	if !valid {
		status.AUTH_InvalidCredentials(c)
		return nil
	}

	fmt.Println(valid)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = "Joe"

	return nil
}

// AUTH_JWTLogout handles the logout process of the user
func AUTH_JWTLogout(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/login")
}
