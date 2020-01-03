package models

import (
	// "fmt"
	"net/http"
	"time"

	"github.com/Techassi/gomark/internal/server/status"
	"github.com/Techassi/gomark/internal/util"

	"github.com/dgrijalva/jwt-go"
	dgoogauth "github.com/dgryski/dgoogauth"
	"github.com/labstack/echo/v4"
)

type Authentication struct {
	Type    string
	User    *User
	Context echo.Context
	App     *App
}

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// GENERAL FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// Init initiates the authentication struct
func (a *Authentication) Init(t string, c echo.Context) {
	a.Type = t
	a.Context = c
	a.App = c.Get("app").(*App)
}

// CheckCookie checks if there is an Authorization cookie and if there is checks
// if it is valid
func (a *Authentication) CheckAuthorizationCookie() map[string]interface{} {
	_, err := a.Context.Cookie("Authorization")
	if err == nil {
		user := a.Context.Get("user")
		if user == nil {
			return status.AUTH_InvalidCredentials()
		}

		return status.AUTH_AlreadySignedIn()
	}

	return nil
}

// SetUser sets the user in the authentication struct
func (a *Authentication) SetUser(u *User) {
	a.User = u
}

// ValidateCredentials validates if the provided username and password are correct
func (a *Authentication) ValidateCredentials(u, p string) (bool, User) {
	return a.App.DB.ValidateCredentials(u, p)
}

// ValidateNewCredentials validates if the provided username and password are ready
// to be used for a new user
func (a *Authentication) ValidateNewCredentials(u, p string) bool {
	return a.App.DB.ValidateNewCredentials(u, p)
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////// REGISTER FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// TODO: Add error handling
func (a *Authentication) Register(user, pass, last, first, mail string) {
	a.App.DB.Register(user, pass, last, first, mail)
}

////////////////////////////////////////////////////////////////////////////////
///////////////////////////////// 2FA FUNCTIONS ////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// Uses2FA returns if a user uses 2FA
func (a *Authentication) Uses2FA() bool {
	return a.User.TwoFA
}

// SetTemp2FAToken sets a temporary 2FA token to validate if the user can access
// the 2FA code page
func (a *Authentication) SetTemp2FAToken() map[string]interface{} {
	expirationTime := time.Now().Add(time.Minute * 5)
	temp2FAToken, temp2FATokenErr := util.RandomCryptoString(10)
	if temp2FATokenErr != nil {
		return status.AUTH_2FATempTokenCreateError()
	}

	c := new(http.Cookie)
	c.Name = "TempTwoFAToken"
	c.Path = "/"
	c.Value = temp2FAToken
	c.Expires = expirationTime
	a.Context.SetCookie(c)

	a.App.DB.SetTemp2FAToken(a.User, temp2FAToken, expirationTime)

	return nil
}

// CheckTemp2FAToken checks if the TempTwoFAToken cookie is set to validate if
// the user can access this page. If valid the user gets set in the authentication
// struct
func (a *Authentication) CheckTemp2FAToken() map[string]interface{} {
	tempToken, tempTokenErr := a.Context.Cookie("TempTwoFAToken")
	if tempTokenErr != nil {
		return status.AUTH_2FATempTokenError()
	}

	a.User = a.App.DB.CheckTemp2FAToken(tempToken.Value)
	if a.User.ID == 0 {
		return status.AUTH_2FATempTokenError()
	}

	return nil
}

// Validate2FACode checks if the provided 2FA code is valid
func (a *Authentication) Validate2FACode() bool {
	// Set up OTPConfig
	otpc := &dgoogauth.OTPConfig{
		Secret:      a.User.TwoFAKey,
		WindowSize:  3,
		HotpCounter: 0,
	}

	// Check 2FA code provided by user (input)
	code := a.Context.FormValue("twofacode")
	valid, err := otpc.Authenticate(code)
	if err != nil || !valid {
		return false
	}

	// Code is valid, so return true
	return true
}

////////////////////////////////////////////////////////////////////////////////
///////////////////////////////// JWT FUNCTIONS ////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// SetJWTToken creates and sets a JWT token
func (a *Authentication) SetJWTToken() map[string]interface{} {
	// Create a new JWT token and get the current time + 24 hours to set the
	// expiration time
	token := jwt.New(jwt.SigningMethodHS256)
	expirationTime := time.Now().Add(time.Hour * 24)
	expirationTimeUnix := expirationTime.Unix()

	// Set JWT claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = a.User.Username
	claims["exp"] = expirationTimeUnix

	// Sign the JWT token
	config := a.App.GetConfig()
	t, err := token.SignedString([]byte(config.Security.Jwt.Secret))
	if err != nil {
		return status.AUTH_JWTTokenSigningError(err)
	}

	// Set response cookie with token
	c := new(http.Cookie)
	c.Name = "Authorization"
	c.Path = "/"
	c.Value = t
	c.Expires = expirationTime
	a.Context.SetCookie(c)

	return nil
}