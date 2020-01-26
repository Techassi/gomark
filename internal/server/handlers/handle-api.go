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
/////////////////////////////// ENTITY FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// APIGetSubEntities gets all sub entities from a parent folder
func APIGetSubEntities(c echo.Context) error {
	app := c.Get("app").(*m.App)

	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["userid"].(float64))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":   http.StatusOK,
		"entities": app.DB.GetSubEntities(c.Param("hash"), userID),
	})
}

// APIPostEntityToFolder saves any type of entity to a folder
func APIPostEntityToFolder(c echo.Context) error {
	app := c.Get("app").(*m.App)

	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	parent := c.Param("hash")
	entity := c.QueryParam("type")
	if parent == "" {
		return c.JSON(http.StatusOK, status.API_NoHashProvided())
	}

	switch entity {
	case "folder":
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
	case "bookmark":
		bookmarkName := c.FormValue("bookmark-name")
		bookmarkDesc := c.FormValue("bookmark-desc")
		bookmarkURL := c.FormValue("bookmark-url")
		bookmarkHash := util.EntityHash(bookmarkName, bookmarkURL)
		bookmarkOwner := uint(claims["userid"].(float64))

		b := m.Bookmark{
			OwnerID:     bookmarkOwner,
			Name:        bookmarkName,
			Hash:        bookmarkHash,
			Shared:      false,
			Pinned:      false,
			ClickedOn:   0,
			Description: bookmarkDesc,
			URL:         bookmarkURL,
			ImageURL:    "",
			HasParent:   true,
		}

		app.DB.SaveBookmarkToFolder(c.Param("hash"), b)
		return c.JSON(http.StatusOK, status.API_GeneralSuccess())
	default:
		return c.JSON(http.StatusOK, status.API_WrongType())
	}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////// BOOKMARK FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// APIGetRecentBookmarks gets recent bookmarks
func APIGetRecentBookmarks(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

// APIGetBookmarks gets all bookmarks
func APIGetBookmarks(c echo.Context) error {
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

// APIGetBookmark gets a single bookmark matching the hash
func APIGetBookmark(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

// APIGetBookmarkTags gets tags attached to a single bookmark
func APIGetBookmarkTags(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

// APIPostBookmark saves a bookmark
func APIPostBookmark(c echo.Context) error {
	app := c.Get("app").(*m.App)

	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	bookmarkName := c.FormValue("bookmark-name")
	bookmarkDesc := c.FormValue("bookmark-desc")
	bookmarkURL := c.FormValue("bookmark-url")
	bookmarkHash := util.EntityHash(bookmarkName, bookmarkURL)
	bookmarkOwner := uint(claims["userid"].(float64))

	b := m.Bookmark{
		OwnerID:     bookmarkOwner,
		Name:        bookmarkName,
		Hash:        bookmarkHash,
		Shared:      false,
		Pinned:      false,
		ClickedOn:   0,
		Description: bookmarkDesc,
		URL:         bookmarkURL,
		ImageURL:    "",
		HasParent:   false,
	}

	app.DB.SaveBookmark(b)
	return c.JSON(http.StatusOK, status.API_GeneralSuccess())
}

// APIUpdateBookmark updates a bookmark
func APIUpdateBookmark(c echo.Context) error {
	return nil
}

// APIPostBookmarkTags saves one or more tags to a bookmark
func APIPostBookmarkTags(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// FOLDER FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// APIGetFolders gets all folders
func APIGetFolders(c echo.Context) error {
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

// APIGetSubFolders gets all subfolders from a parent folder
func APIGetSubFolders(c echo.Context) error {
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
		"folders": app.DB.GetSubFolders(c.Param("hash"), userID),
	})
}

// APIPostFolder saves a folder
func APIPostFolder(c echo.Context) error {
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
