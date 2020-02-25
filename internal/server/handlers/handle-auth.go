package handlers

import (
	"net/http"

	auth "github.com/Techassi/gomark/internal/authentication"
	"github.com/Techassi/gomark/internal/server/status"
	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo/v4"
	// qrcode "github.com/skip2/go-qrcode"
)

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// GENERAL FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// AUTH_JWTError handles the redirect to the login page if no JWT token is
// present
func AUTH_JWTError(e error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/login")
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////// REGISTER FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// AUTH_Register handles the register process of a new user
func AUTH_Register(c echo.Context) error {
	// Initiate Authentication
	a := &auth.Authenticator{}
	a.Init("register", c)

	if !a.App.RegisterEnabled() {
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
	if !a.ValidateNewCredentials(username, password) {
		return c.JSON(200, status.AUTH_InvalidNewCredentials())
	}

	// Register the new account
	err := a.Register(username, password, lastname, firstname, email)
	if err != nil {
		return c.JSON(http.StatusOK, status.AUTH_NotRegistered)
	}
	return c.JSON(200, status.AUTH_SuccessfullyRegistered())
}

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// LOGIN FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// AUTH_Login handles the user authentication via the DB to login the user
func AUTH_Login(c echo.Context) error {
	// Initiate Authentication
	a := &auth.Authenticator{}
	a.Init("login", c)

	// Extract the provided user information
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Check if the provided credentials are valid
	valid, user := a.ValidateCredentials(username, password)
	if !valid {
		return c.JSON(200, status.AUTH_InvalidCredentials())
	}

	a.SetUser(&user)

	// Check if the user has 2FA activated, if yes proceed to 2FA code authentication.
	// If not continue with JWT authentication
	if a.Uses2FA() {
		// Set temporary token to validate the user can access the 2FA code page
		if err := a.SetTemp2FAToken(); err != nil {
			return c.JSON(200, err)
		}

		return c.JSON(200, status.AUTH_2FARequired())
	}

	// Continue with the JWT login flow
	return JWTLoginFlow(a)
}

func JWTLoginFlow(a *auth.Authenticator) error {
	// Set JWT token
	if err := a.SetJWTToken(); err != nil {
		return a.Context.JSON(200, err)
	}
	return a.Context.JSON(200, status.AUTH_SuccessfullySignedIn())
}

////////////////////////////////////////////////////////////////////////////////
///////////////////////////////// 2FA FUNCTIONS ////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// AUTH_Create2FACode handles the 2FA code validation
func AUTH_Create2FACode(c echo.Context) error {
	// Initiate Authentication
	a := &auth.Authenticator{}
	a.Init("2fa", c)

	user := c.Get("user")
	if user == nil {
		return c.Redirect(http.StatusMovedPermanently, "/login")
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["userid"].(float64))
	username := claims["username"].(string)

	code, err := a.Create2FACode(userID, username)
	if err != nil {
		return c.JSON(http.StatusOK, status.AUTH_2FAQRCodeError())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
		"code":   code,
	})
}

// AUTH_2FACode handles the 2FA code validation
func AUTH_2FACode(c echo.Context) error {
	// Initiate Authentication
	a := &auth.Authenticator{}
	a.Init("2fa", c)

	// Check Temp2FAToken if valid
	a.CheckTemp2FAToken()

	// Validate if the provided 2FA code is valid
	valid := a.Validate2FACode()
	if !valid {
		return c.JSON(200, status.AUTH_2FAAuthenticationError())
	}

	// Continue with the JWT login flow
	return JWTLoginFlow(a)
}
