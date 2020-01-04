package handlers

import (
	"net/http"

	m "github.com/Techassi/gomark/internal/models"
	"github.com/Techassi/gomark/internal/server/status"
	"github.com/Techassi/gomark/internal/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////// BOOKMARK FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func API_GetRecentBookmarks(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

func API_GetBookmarks(c echo.Context) error {
	app := c.Get("app").(*m.App)

	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	b := app.DB.GetBookmarksByUserID(uint(claims["userid"].(float64)))
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    http.StatusOK,
		"bookmarks": b,
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
	app := c.Get("app").(*m.App)

	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	bookmarkName := c.FormValue("bookmark-name")
	bookmarkDesc := c.FormValue("bookmark-desc")
	bookmarkUrl := c.FormValue("bookmark-url")
	bookmarkHash := util.EntityHash(bookmarkName, bookmarkUrl)

	b := m.Entity{
		Type:      "bookmark",
		Name:      bookmarkName,
		Hash:      bookmarkHash,
		Shared:    false,
		ClickedOn: 0,
		BookmarkData: &m.BookmarkData{
			Description: bookmarkDesc,
			URL:         bookmarkUrl,
			ImageURL:    "",
		},
	}

	app.DB.SaveEntity(b, uint(claims["userid"].(float64)))

	return c.JSON(http.StatusOK, status.API_GeneralSuccess())
}

func API_UpdateBookmark(c echo.Context) error {
	return nil
}

func API_PostBookmarkTags(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// FOLDER FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func API_PostFolder(c echo.Context) error {
	app := c.Get("app").(*m.App)

	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	folderName := c.FormValue("folder-name")
	folderHash := util.EntityHashPlusString(folderName)

	f := m.Entity{
		Type:      "folder",
		Name:      folderName,
		Hash:      folderHash,
		Shared:    false,
		ClickedOn: 0,
		FolderData: &m.FolderData{
			ChildEntitiesCount: 0,
			ChildEntities:      []m.Entity{},
		},
	}

	app.DB.SaveEntity(f, uint(claims["userid"].(float64)))

	return c.JSON(http.StatusOK, status.API_GeneralSuccess())
}

func API_GetFolders(c echo.Context) error {
	app := c.Get("app").(*m.App)

	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	f := app.DB.GetFoldersByUserID(uint(claims["userid"].(float64)))
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"folders": f,
	})
}

func API_PostEntityToFolder(c echo.Context) error {
	app := c.Get("app").(*m.App)

	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	}

	bookmarkName := c.FormValue("bookmark-name")
	bookmarkDesc := c.FormValue("bookmark-desc")
	bookmarkUrl := c.FormValue("bookmark-url")
	bookmarkHash := util.EntityHash(bookmarkName, bookmarkUrl)

	e := m.Entity{
		Type:      "bookmark",
		Name:      bookmarkName,
		Hash:      bookmarkHash,
		Shared:    false,
		ClickedOn: 0,
		BookmarkData: &m.BookmarkData{
			Description: bookmarkDesc,
			URL:         bookmarkUrl,
			ImageURL:    "",
		},
	}

	app.DB.SaveEntityToFolder(e, c.Param("hash"))

	return c.JSON(http.StatusOK, status.API_GeneralSuccess())
}
