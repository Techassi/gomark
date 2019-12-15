package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func API_GetRecentBookmarks(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

func API_GetBookmarks(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

func API_GetBookmark(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

func API_GetBookmarkTags(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

func API_PostBookmark(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

func API_PostBookmarkTags(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}
