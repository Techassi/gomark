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

// ApiGetSubEntities gets all sub entities from a parent folder
func (app *App) ApiGetSubEntities(c echo.Context) error {
	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.ApiGeneralAccesError())
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["userid"].(float64))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":   http.StatusOK,
		"entities": app.DB.GetSubEntities(c.Param("hash"), userID),
	})
}

// ApiPostEntityToFolder saves any type of entity to a folder
func (app *App) ApiPostEntityToFolder(c echo.Context) error {
	// user := c.Get("user")
	// if user == nil {
	// 	return c.JSON(http.StatusOK, status.ApiGeneralAccesError())
	// }

	// token := user.(*jwt.Token)
	// claims := token.Claims.(jwt.MapClaims)

	// parent := c.Param("hash")
	// entity := c.QueryParam("type")
	// if parent == "" {
	// 	return c.JSON(http.StatusOK, status.ApiNoHashProvided())
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
	// 		return c.JSON(http.StatusOK, status.ApiGeneralAccesError())
	// 	}
	// 	return c.JSON(http.StatusOK, status.ApiGeneralSuccess())
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
	// 	return c.JSON(http.StatusOK, status.ApiGeneralSuccess())
	// default:
	// 	return c.JSON(http.StatusOK, status.ApiWrongType())
	// }
	return nil
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////// BOOKMARK FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// ApiGetRecentBookmarks gets recent bookmarks
func (app *App) ApiGetRecentBookmarks(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"entities": app.DB.GetRecentBookmarks(1),
	})
}

// ApiGetBookmarks gets all bookmarks
func (app *App) ApiGetBookmarks(c echo.Context) error {
	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.ApiGeneralAccesError())
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["userid"].(float64))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    http.StatusOK,
		"bookmarks": app.DB.GetBookmarks(userID),
	})
}

// ApiGetBookmark gets a single bookmark matching the hash
func (app *App) ApiGetBookmark(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

// ApiGetBookmarkTags gets tags attached to a single bookmark
func (app *App) ApiGetBookmarkTags(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

// ApiPostBookmark saves a bookmark
func (app *App) ApiPostBookmark(c echo.Context) error {
	// user := c.Get("user")
	// if user == nil {
	// 	return c.JSON(http.StatusOK, status.ApiGeneralAccesError())
	// }

	// token := user.(*jwt.Token)
	// claims := token.Claims.(jwt.MapClaims)

	entityName := c.FormValue("bookmark-name")
	entityDesc := c.FormValue("bookmark-desc")
	entityURL := c.FormValue("bookmark-url")
	entityHash := util.EntityHash(entityName, entityURL)
	// entityOwner := uint(claims["userid"].(float64))

	e := m.Entity{
		OwnerID: 1,
		Type:    "bookmark",
		Hash:    entityHash,
		Name:    entityName,
		Bookmark: &m.Bookmark{
			Description: entityDesc,
			URL:         entityURL,
		},
	}

	app.DB.SaveEntity(e)
	go app.Scheduler.Schedule(app.Scheduler.Job("download-meta", entityHash))
	return c.JSON(http.StatusOK, status.ApiGeneralSuccess())
}

// ApiUpdateBookmark updates a bookmark
func (app *App) ApiUpdateBookmark(c echo.Context) error {
	return nil
}

// ApiPostBookmarkTags saves one or more tags to a bookmark
func (app *App) ApiPostBookmarkTags(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}

func (app *App) ApiShareBookmark(c echo.Context) error {
	shareHash, err := app.DB.ShareBookmark(c.Param("hash"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, status.ApiInternalError(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"share_hash": shareHash,
	})
}

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// FOLDER FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// ApiGetFolders gets all folders
func (app *App) ApiGetFolders(c echo.Context) error {
	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.ApiGeneralAccesError())
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["userid"].(float64))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"folders": app.DB.GetFolders(userID),
	})
}

// ApiGetSubFolders gets all subfolders from a parent folder
func (app *App) ApiGetSubFolders(c echo.Context) error {
	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusOK, status.ApiGeneralAccesError())
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["userid"].(float64))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"folders": app.DB.GetSubFolders(c.Param("hash"), userID),
	})
}

// ApiPostFolder saves a folder
func (app *App) ApiPostFolder(c echo.Context) error {
	// app := c.Get("app").(*app.App)

	// user := c.Get("user")
	// if user == nil {
	// 	return c.JSON(http.StatusOK, status.ApiGeneralAccesError())
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
	return c.JSON(http.StatusOK, status.ApiGeneralSuccess())
}

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// EVENT FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (app *App) ApiPostEvent(c echo.Context) error {
	job := app.Scheduler.Job(c.FormValue("event-type"), c.FormValue("event-data"))
	go app.Scheduler.Schedule(job)
	return c.JSON(http.StatusOK, status.ApiGeneralSuccess())
}
