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
	userID := uint(claims["userid"].(float64))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    http.StatusOK,
		"bookmarks": app.DB.GetBookmarksByUserID(userID),
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
	bookmarkOwner := uint(claims["userid"].(float64))

	b := m.Bookmark{
		OwnerID:     bookmarkOwner,
		Name:        bookmarkName,
		Hash:        bookmarkHash,
		Shared:      false,
		Pinned:      false,
		ClickedOn:   0,
		Description: bookmarkDesc,
		URL:         bookmarkUrl,
		ImageURL:    "",
	}

	app.DB.SaveBookmark(b)
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
	folderOwner := uint(claims["userid"].(float64))

	f := m.Folder{
		OwnerID:       folderOwner,
		Name:          folderName,
		Hash:          folderHash,
		Shared:        false,
		ClickedOn:     0,
		ChildrenCount: 0,
		HasParent:     false,
	}

	app.DB.SaveFolder(f)
	return c.JSON(http.StatusOK, status.API_GeneralSuccess())
}

func API_PostSubFolder(c echo.Context) error {
	app := c.Get("app").(*m.App)

	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	folderName := c.FormValue("folder-name")
	folderHash := util.EntityHashPlusString(folderName)
	folderOwner := uint(claims["userid"].(float64))

	f := m.Folder{
		OwnerID:       folderOwner,
		Name:          folderName,
		Hash:          folderHash,
		Shared:        false,
		ClickedOn:     0,
		ChildrenCount: 0,
		HasParent:     true,
	}

	if err := app.DB.SaveSubFolder(c.Param("hash"), f); err != nil {
		return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	}
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
	userID := uint(claims["userid"].(float64))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"folders": app.DB.GetFolders(userID),
	})
}
