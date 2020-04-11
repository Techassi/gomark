package app

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

// API_GetSubEntities gets all sub entities from a parent folder
func (app *App) API_GetSubEntities(c echo.Context) error {
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

// API_PostEntityToFolder saves any type of entity to a folder
func (app *App) API_PostEntityToFolder(c echo.Context) error {
	// user := c.Get("user")
	// if user == nil {
	// 	return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	// }

	// token := user.(*jwt.Token)
	// claims := token.Claims.(jwt.MapClaims)

	// parent := c.Param("hash")
	// entity := c.QueryParam("type")
	// if parent == "" {
	// 	return c.JSON(http.StatusOK, status.API_NoHashProvided())
	// }

	// switch entity {
	// case "folder":
	// 	folderName := c.FormValue("folder-name")
	// 	folderHash := util.EntityHashPlusString(folderName)
	// 	folderOwner := uint(claims["userid"].(float64))

	// 	f := m.Folder{
	// 		OwnerID:       folderOwner,
	// 		Name:          folderName,
	// 		Hash:          folderHash,
	// 		Shared:        false,
	// 		ClickedOn:     0,
	// 		ChildrenCount: 0,
	// 		HasParent:     true,
	// 	}

	// 	if err := app.DB.SaveSubFolder(c.Param("hash"), f); err != nil {
	// 		return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	// 	}
	// 	return c.JSON(http.StatusOK, status.API_GeneralSuccess())
	// case "bookmark":
	// 	bookmarkName := c.FormValue("bookmark-name")
	// 	bookmarkDesc := c.FormValue("bookmark-desc")
	// 	bookmarkURL := c.FormValue("bookmark-url")
	// 	bookmarkHash := util.EntityHash(bookmarkName, bookmarkURL)
	// 	bookmarkOwner := uint(claims["userid"].(float64))

	// 	b := m.Bookmark{
	// 		OwnerID:     bookmarkOwner,
	// 		Name:        bookmarkName,
	// 		Hash:        bookmarkHash,
	// 		Shared:      false,
	// 		Pinned:      false,
	// 		ClickedOn:   0,
	// 		Description: bookmarkDesc,
	// 		URL:         bookmarkURL,
	// 		ImageURL:    "",
	// 		HasParent:   true,
	// 	}

	// 	app.DB.SaveBookmarkToFolder(c.Param("hash"), b)
	// 	return c.JSON(http.StatusOK, status.API_GeneralSuccess())
	// default:
	// 	return c.JSON(http.StatusOK, status.API_WrongType())
	// }
	return nil
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////// BOOKMARK FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// API_GetRecentBookmarks gets recent bookmarks
func (app *App) API_GetRecentBookmarks(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

// API_GetBookmarks gets all bookmarks
func (app *App) API_GetBookmarks(c echo.Context) error {
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

// API_GetBookmark gets a single bookmark matching the hash
func (app *App) API_GetBookmark(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

// API_GetBookmarkTags gets tags attached to a single bookmark
func (app *App) API_GetBookmarkTags(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

// API_PostBookmark saves a bookmark
func (app *App) API_PostBookmark(c echo.Context) error {
	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	entityName := c.FormValue("bookmark-name")
	entityDesc := c.FormValue("bookmark-desc")
	entityURL := c.FormValue("bookmark-url")
	entityHash := util.EntityHash(entityName, entityURL)
	entityOwner := uint(claims["userid"].(float64))

	e := m.Entity{
		OwnerID: entityOwner,
		Type:    "bookmark",
		Hash:    entityHash,
		Name:    entityName,
		Bookmark: &m.Bookmark{
			Description: entityDesc,
			URL:         entityURL,
		},
	}

	app.DB.SaveEntity(e)
	return c.JSON(http.StatusOK, status.API_GeneralSuccess())
}

// API_UpdateBookmark updates a bookmark
func (app *App) API_UpdateBookmark(c echo.Context) error {
	return nil
}

// API_PostBookmarkTags saves one or more tags to a bookmark
func (app *App) API_PostBookmarkTags(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

func (app *App) API_ShareBookmark(c echo.Context) error {
	shareHash, err := app.DB.ShareBookmark(c.Param("hash"))
	if err != nil {
		return c.JSON(http.StatusOK, status.API_InternalError(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     http.StatusOK,
		"share_hash": shareHash,
	})
}

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// FOLDER FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// API_GetFolders gets all folders
func (app *App) API_GetFolders(c echo.Context) error {
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

// API_GetSubFolders gets all subfolders from a parent folder
func (app *App) API_GetSubFolders(c echo.Context) error {
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

// API_PostFolder saves a folder
func (app *App) API_PostFolder(c echo.Context) error {
	// app := c.Get("app").(*app.App)

	// user := c.Get("user")
	// if user == nil {
	// 	return c.JSON(http.StatusOK, status.API_GeneralAccesError())
	// }

	// token := user.(*jwt.Token)
	// claims := token.Claims.(jwt.MapClaims)

	// folderName := c.FormValue("folder-name")
	// folderHash := util.EntityHashPlusString(folderName)
	// folderOwner := uint(claims["userid"].(float64))

	// f := m.Folder{
	// 	OwnerID:       folderOwner,
	// 	Name:          folderName,
	// 	Hash:          folderHash,
	// 	Shared:        false,
	// 	ClickedOn:     0,
	// 	ChildrenCount: 0,
	// 	HasParent:     false,
	// }

	// app.DB.SaveFolder(f)
	return c.JSON(http.StatusOK, status.API_GeneralSuccess())
}

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// EVENT FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (app *App) API_PostEvent(c echo.Context) error {
	job := app.Scheduler.Job(c.FormValue("event-type"), c.FormValue("event-data"))
	go app.Scheduler.Schedule(job)
	return c.JSON(http.StatusOK, status.API_GeneralSuccess())
}
