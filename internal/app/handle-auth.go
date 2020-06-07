package app

import (
	"net/http"
	"time"

	"github.com/Techassi/gomark/internal/server/status"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	// qrcode "github.com/skip2/go-qrcode"
)

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// GENERAL FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// AuthJWTError handles the redirect to the login page if no JWT token is
// present
func (app *App) AuthJWTError(e error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/login")
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////// REGISTER FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// AuthRegister handles the register process of a new user
func (app *App) AuthRegister(c echo.Context) error {
	// TODO: Implement authentication logic using goauth
	var (
		username  = c.FormValue("username")
		password  = c.FormValue("password")
		firstname = c.FormValue("firstname")
		lastname  = c.FormValue("lastname")
		email     = c.FormValue("email")
	)

	ok := app.DB.ValidateNewCredentials(username, password)
	if !ok {
		return c.JSON(http.StatusBadRequest, status.AuthInvalidNewCredentials())
	}

	err := app.DB.Register(username, password, lastname, firstname, email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, status.AuthNotRegistered())
	}

	return c.JSON(http.StatusOK, status.AuthSuccessfullyRegistered())
}

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// LOGIN FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// AuthLogin handles the user authentication via the DB to login the user
func (app *App) AuthLogin(c echo.Context) error {
	// TODO: Implement authentication logic using goauth
	// We are using the default JWT login flow until goauth is ready
	var (
		token       = jwt.New(jwt.SigningMethodHS256)
		expTime     = time.Now().Add(time.Hour * 48)
		expTimeUnix = expTime.Unix()
		username    = c.FormValue("username")
		password    = c.FormValue("password")
	)

	ok, user := app.DB.ValidateCredentials(username, password)
	if !ok {
		return c.JSON(http.StatusUnauthorized, status.AuthInvalidCredentials())
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = user.Username
	claims["exp"] = expTimeUnix

	t, err := token.SignedString([]byte(app.Config.Security.Jwt.Secret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, status.AuthJWTTokenSigningError(err))
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Path = "/"
	cookie.Value = t
	cookie.Expires = expTime
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, status.AuthSuccessfullySignedIn())
}

////////////////////////////////////////////////////////////////////////////////
///////////////////////////////// 2FA FUNCTIONS ////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// AuthCreate2FACode handles the 2FA code validation
func (app *App) AuthCreate2FACode(c echo.Context) error {
	// Initiate Authentication
	// a := &auth.Authenticator{}
	// a.Init("2fa", c)

	// user := c.Get("user")
	// if user == nil {
	// 	return c.Redirect(http.StatusMovedPermanently, "/login")
	// }

	// token := user.(*jwt.Token)
	// claims := token.Claims.(jwt.MapClaims)
	// userID := uint(claims["userid"].(float64))
	// username := claims["username"].(string)

	// code, err := a.Create2FACode(userID, username)
	// if err != nil {
	// 	return c.JSON(http.StatusOK, status.Auth2FAQRCodeError())
	// }

	// return c.JSON(http.StatusOK, map[string]interface{}{
	// 	"status": http.StatusOK,
	// 	"code":   code,
	// })
	return nil
}

// Auth2FACode handles the 2FA code validation
func (app *App) Auth2FACode(c echo.Context) error {
	// Initiate Authentication
	// a := &auth.Authenticator{}
	// a.Init("2fa", c)

	// // Check Temp2FAToken if valid
	// a.CheckTemp2FAToken()

	// // Validate if the provided 2FA code is valid
	// valid := a.Validate2FACode()
	// if !valid {
	// 	return c.JSON(200, status.Auth2FAAuthenticationError())
	// }

	// // Continue with the JWT login flow
	// return app.JWTLoginFlow(a)
	return nil
}
