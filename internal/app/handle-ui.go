package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

////////////////////////////////////////////////////////////////////////////////
///////////////////////////// AUTHENTICATION PAGES /////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (app *App) UiLoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{})
}

func (app *App) Ui2FACodePage(c echo.Context) error {
	return c.Render(http.StatusOK, "code.html", map[string]interface{}{})
}

func (app *App) UiRegisterPage(c echo.Context) error {
	if !app.RegisterEnabled() {
		return c.Redirect(http.StatusMovedPermanently, "/login")
	}

	return c.Render(http.StatusOK, "register.html", map[string]interface{}{})
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////// ERROR PAGES /////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (app *App) Ui404Page(c echo.Context) error {
	return c.Render(http.StatusOK, "404.html", map[string]interface{}{})
}

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// BOOKMARK PAGES ////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (app *App) UiSharedEntityPage(c echo.Context) error {
	return c.Render(http.StatusOK, "shared-entity.html", map[string]interface{}{
		"config": app.Config,
		"entity": app.DB.GetShared(c.Param("hash")),
	})
}

func (app *App) UiRecentBookmarksPage(c echo.Context) error {
	return c.Render(http.StatusOK, "recent-bookmarks.html", map[string]interface{}{})
}

func (app *App) UiBookmarksPage(c echo.Context) error {
	return c.Render(http.StatusOK, "bookmarks.html", map[string]interface{}{})
}

func (app *App) UiBookmarkPage(c echo.Context) error {
	return c.Render(http.StatusOK, "bookmark.html", map[string]interface{}{})
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////// NOTE PAGES //////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (app *App) UiNotePage(c echo.Context) error {
	return c.Render(http.StatusOK, "note.html", map[string]interface{}{})
}

func (app *App) UiNotesPage(c echo.Context) error {
	return c.Render(http.StatusOK, "notes.html", map[string]interface{}{})
}

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// DASHBOARD PAGES ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// UiHomePage renders the home page
func (app *App) UiHomePage(c echo.Context) error {
	// user := c.Get("user")
	// if user == nil {
	// 	return c.Redirect(http.StatusMovedPermanently, "/login")
	// }

	// token := user.(*jwt.Token)
	// claims := token.Claims.(jwt.MapClaims)
	// userID := uint(claims["userid"].(float64))

	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"config": app.Config,
		// "user":     claims["username"].(string),
		"entities": app.DB.GetBookmarks(),
	})
}

////////////////////////////////////////////////////////////////////////////////
///////////////////////////////// SHARED PAGES /////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (app *App) UiSharedPage(c echo.Context) error {
	return c.Render(http.StatusOK, "shared.html", map[string]interface{}{})
}
