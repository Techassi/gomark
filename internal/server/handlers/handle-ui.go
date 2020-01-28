package handlers

import (
	"net/http"

	m "github.com/Techassi/gomark/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

////////////////////////////////////////////////////////////////////////////////
///////////////////////////// AUTHENTICATION PAGES /////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func UILoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{})
}

func UI2FACodePage(c echo.Context) error {
	return c.Render(http.StatusOK, "code.html", map[string]interface{}{})
}

func UIRegisterPage(c echo.Context) error {
	app := c.Get("app").(*m.App)

	if !app.RegisterEnabled() {
		return c.Redirect(http.StatusMovedPermanently, "/login")
	}

	return c.Render(http.StatusOK, "register.html", map[string]interface{}{})
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////// ERROR PAGES /////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func UI404Page(c echo.Context) error {
	return c.Render(http.StatusOK, "404.html", map[string]interface{}{})
}

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// BOOKMARK PAGES ////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func UISharedBookmarkPage(c echo.Context) error {
	return c.Render(http.StatusOK, "shared-bookmark.html", map[string]interface{}{})
}

func UIRecentBookmarksPage(c echo.Context) error {
	return c.Render(http.StatusOK, "recent-bookmarks.html", map[string]interface{}{})
}

func UIBookmarksPage(c echo.Context) error {
	return c.Render(http.StatusOK, "bookmarks.html", map[string]interface{}{})
}

func UIBookmarkPage(c echo.Context) error {
	return c.Render(http.StatusOK, "bookmark.html", map[string]interface{}{})
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////// NOTE PAGES //////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func UINotePage(c echo.Context) error {
	return c.Render(http.StatusOK, "note.html", map[string]interface{}{})
}

func UINotesPage(c echo.Context) error {
	return c.Render(http.StatusOK, "notes.html", map[string]interface{}{})
}

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// DASHBOARD PAGES ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// UIHomePage renders the home page
func UIHomePage(c echo.Context) error {
	app := c.Get("app").(*m.App)

	user := c.Get("user")
	if user == nil {
		return c.Redirect(http.StatusMovedPermanently, "/login")
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["userid"].(float64))

	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"config":    app.Config,
		"user":      claims["username"].(string),
		"bookmarks": app.DB.GetBookmarksByUserID(userID),
	})
}

////////////////////////////////////////////////////////////////////////////////
///////////////////////////////// SHARED PAGES /////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func UISharedPage(c echo.Context) error {
	return c.Render(http.StatusOK, "shared.html", map[string]interface{}{})
}
