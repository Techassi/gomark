package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/Techassi/gomark/internal/constants"
	m "github.com/Techassi/gomark/internal/models"
	"github.com/Techassi/gomark/internal/server/status"
	"github.com/Techassi/gomark/internal/util"

	"github.com/dgrijalva/jwt-go"
	dgoogauth "github.com/dgryski/dgoogauth"
	"github.com/labstack/echo/v4"
	qrcode "github.com/skip2/go-qrcode"
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
	// check if cookie is already set
	_, err := c.Cookie("Authorization")
	if err == nil {
		return status.AUTH_AlreadySignedIn(c)
	}

	app := c.Get("app").(*m.App)
	u := m.User{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}

	// check if the provided credentials are valid
	valid := app.DB.ValidCredentials(&u)
	if !valid {
		return status.AUTH_InvalidCredentials(c)
	}

	// Create 2FA Code if user isn't using 2FA already
	if !u.TwoFA {
		// Create 2FA secret
		twoFASecretBase32, twoFAErr := util.RandomCryptoString(10)
		if twoFAErr != nil {
			return status.AUTH_2FASecretError(c)
		}

		// Create the OTP Uri
		uri, uriErr := url.Parse("otpauth://totp")
		if uriErr != nil {
			return status.AUTH_2FAUriError(c)
		}

		uri.Path += fmt.Sprintf("/%s:%s", url.PathEscape(constants.TWOFA_ISSUER), u.Username)

		params := url.Values{}
		params.Add("secret", twoFASecretBase32)
		params.Add("issuer", constants.TWOFA_ISSUER)
		uri.RawQuery = params.Encode()

		// Generate QR Code for Google Authenticator
		qrErr := qrcode.WriteFile(uri.String(), qrcode.Medium, 256, filepath.Join(util.GetAbsPath("public/2fa"), "qr.png"))
		if qrErr != nil {
			log.Fatal(qrErr)
			return status.AUTH_2FAQRCodeError(c)
		}

		// Set that the user is now using 2FA and save the key
		app.DB.Update2FA(&u, twoFASecretBase32)
	}

	// Set temporary token to validate the user can access the 2FA code page
	tempTwoFAToken, tempTwoFATokenErr := util.RandomCryptoString(10)
	if tempTwoFATokenErr != nil {
		return status.AUTH_2FATempTokenCreateError(c)
	}

	currTime := time.Now().Add(time.Minute * 5)

	tokenCookie := new(http.Cookie)
	tokenCookie.Name = "TempTwoFAToken"
	tokenCookie.Path = "/"
	tokenCookie.Value = string(tempTwoFAToken)
	tokenCookie.Expires = currTime
	c.SetCookie(tokenCookie)

	app.DB.SetTempTwoFAToken(&u, string(tempTwoFAToken), currTime)

	return c.Redirect(http.StatusMovedPermanently, "/code")
}

// AUTH_JWT2FACode handles the 2FA code authentication process of the user
func AUTH_JWT2FACode(c echo.Context) error {
	// Check if cookie is already set
	_, err := c.Cookie("Authorization")
	if err == nil {
		return status.AUTH_AlreadySignedIn(c)
	}

	app := c.Get("app").(*m.App)

	// Check if the TempTwoFAToken cookie is set to validate if the user can
	// access this page
	tempToken, tempTokenErr := c.Cookie("TempTwoFAToken")
	if tempTokenErr != nil {
		return status.AUTH_2FATempTokenError(c)
	}

	// Check TempTwoFAToken if valid (exists and not expired)
	key := app.DB.CheckTempTwoFAToken(tempToken.Value)

	// TODO: Delete TempTwoFAToken after successfully using it

	// Set up OTPConfig
	otpc := &dgoogauth.OTPConfig{
		Secret:      key,
		WindowSize:  3,
		HotpCounter: 0,
	}

	// Check 2FA code provided by user (input)
	code := c.FormValue("twofacode")
	valid, err := otpc.Authenticate(code)
	if err != nil || !valid {
		return status.AUTH_2FAAuthenticationError(c)
	}

	// Create a new JWT token and get the current time + 24 hours to set the
	// expiration time
	token := jwt.New(jwt.SigningMethodHS256)
	currTime := time.Now().Add(time.Hour * 24)
	currTimeUnix := currTime.Unix()

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = "joe" // TODO: set this from the database
	claims["exp"] = currTimeUnix

	// Sign token
	t, err := token.SignedString([]byte(app.Config.Security.Jwt.Secret))
	if err != nil {
		return err
	}

	// Set response cookie with token
	tokenCookie := new(http.Cookie)
	tokenCookie.Name = "Authorization"
	tokenCookie.Path = "/"
	tokenCookie.Value = t
	tokenCookie.Expires = currTime
	c.SetCookie(tokenCookie)

	return status.AUTH_SuccessfullySignedIn(c)
}

// AUTH_JWTLogout handles the logout process of the user
func AUTH_JWTLogout(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/login")
}
