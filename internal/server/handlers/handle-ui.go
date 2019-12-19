package handlers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Authentication pages

func UI_LoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{})
}

func UI_TwoFACodePage(c echo.Context) error {
	return c.Render(http.StatusOK, "code.html", map[string]interface{}{})
}

func UI_RegisterPage(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", map[string]interface{}{})
}

// Error pages

func UI_404Page(c echo.Context) error {
	return c.Render(http.StatusOK, "404.html", map[string]interface{}{})
}

// Bookmark pages

func UI_SharedBookmarkPage(c echo.Context) error {
	return c.Render(http.StatusOK, "shared-bookmark.html", map[string]interface{}{})
}

func UI_RecentBookmarksPage(c echo.Context) error {
	return c.Render(http.StatusOK, "recent-bookmarks.html", map[string]interface{}{})
}

func UI_BookmarksPage(c echo.Context) error {
	return c.Render(http.StatusOK, "bookmarks.html", map[string]interface{}{})
}

func UI_BookmarkPage(c echo.Context) error {
	return c.Render(http.StatusOK, "bookmark.html", map[string]interface{}{})
}

// Note pages

func UI_NotePage(c echo.Context) error {
	return c.Render(http.StatusOK, "note.html", map[string]interface{}{})
}

func UI_NotesPage(c echo.Context) error {
	return c.Render(http.StatusOK, "notes.html", map[string]interface{}{})
}

// Dashboard pages

func UI_DashboardPage(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	return c.String(http.StatusOK, claims["name"].(string))
	// return c.Render(http.StatusOK, "dashboard.html", map[string]interface{}{})
}

// Shared pages

func UI_SharedPage(c echo.Context) error {
	return c.Render(http.StatusOK, "shared.html", map[string]interface{}{})
}
