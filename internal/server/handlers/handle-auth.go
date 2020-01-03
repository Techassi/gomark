package handlers

import (
	"net/http"

	m "github.com/Techassi/gomark/internal/models"
	"github.com/Techassi/gomark/internal/server/status"

	"github.com/labstack/echo/v4"
	// qrcode "github.com/skip2/go-qrcode"
)

// AUTH_JWTError handles the redirect to the login page if no JWT token is
// present
func AUTH_JWTError(e error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/login")
}

// AUTH_JWTRegister handles the register process of a new user
func AUTH_Register(c echo.Context) error {
	// Initiate Authentication
	auth := &m.Authentication{}
	auth.Init("register", c)

	if !auth.App.RegisterEnabled() {
		return c.JSON(200, status.ADMIN_RegisterDisabled())
	}

	// Extract the provided user information to check for validity
	username := c.FormValue("username")
	password := c.FormValue("password")

	email := c.FormValue("email")
	lastname := c.FormValue("lastname")
	firstname := c.FormValue("firstname")

	// Check if username and password are valid (username is not used already
	// and password checks all requirements)
	if !auth.ValidateNewCredentials(username, password) {
		return c.JSON(200, status.AUTH_InvalidNewCredentials())
	}

	// Register the new account
	auth.Register(username, password, lastname, firstname, email)

	return c.JSON(200, status.AUTH_SuccessfullyRegistered())
}

// AUTH_Login handles the user authentication via the DB to login the user
func AUTH_Login(c echo.Context) error {
	// Initiate Authentication
	auth := &m.Authentication{}
	auth.Init("login", c)
	if err := auth.CheckAuthorizationCookie(); err != nil {
		return c.JSON(200, err)
	}

	// Extract the provided user information
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Check if the provided credentials are valid
	valid, user := auth.ValidateCredentials(username, password)
	if !valid {
		return c.JSON(200, status.AUTH_InvalidCredentials())
		return nil
	}

	auth.SetUser(&user)

	// Check if the user has 2FA activated, if yes proceed to 2FA code authentication.
	// If not continue with JWT authentication
	if auth.Uses2FA() {
		// Set temporary token to validate the user can access the 2FA code page
		if err := auth.SetTemp2FAToken(); err != nil {
			return c.JSON(200, err)
		}

		// Redirect to code page
		return c.Redirect(http.StatusMovedPermanently, "/code")
	}

	// Continue with the JWT login flow
	return JWTLoginFlow(auth)
}

// AUTH_2FACode handles the 2FA code validation
func AUTH_2FACode(c echo.Context) error {
	// Initiate Authentication
	auth := &m.Authentication{}
	auth.Init("2fa", c)
	err := auth.CheckAuthorizationCookie()
	if err != nil {
		return c.JSON(200, err)
	}

	// Check Temp2FAToken if valid
	auth.CheckTemp2FAToken()

	// Validate if the provided 2FA code is valid
	valid := auth.Validate2FACode()
	if !valid {
		return c.JSON(200, status.AUTH_2FAAuthenticationError())
	}

	// Continue with the JWT login flow
	return JWTLoginFlow(auth)
}

func JWTLoginFlow(auth *m.Authentication) error {
	// Set JWT token
	if err := auth.SetJWTToken(); err != nil {
		return auth.Context.JSON(200, err)
	}

	return auth.Context.JSON(200, status.AUTH_SuccessfullySignedIn())
}
