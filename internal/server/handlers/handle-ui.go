package handlers

import (
    "net/http"

    "github.com/Techassi/gomark/internal/models"

    "github.com/gin-gonic/gin"
)

func UI_HomePage(c *gin.Context) {
    c.HTML(http.StatusOK, "home.html", gin.H{
        "entities": models.GetEntities(),
    })
}

func UI_LoginPage(c *gin.Context) {
    c.HTML(http.StatusOK, "login.html", gin.H{})
}

func UI_RegisterPage(c *gin.Context) {
    c.HTML(http.StatusOK, "register.html", gin.H{})
}

func UI_NotesPage(c *gin.Context) {
    c.HTML(http.StatusOK, "login.html", gin.H{
        "entities": models.GetNotes(),
    })
}

func UI_SharedPage(c *gin.Context) {
    c.HTML(http.StatusOK, "login.html", gin.H{
        "entities": models.GetSharedEntities(),
    })
}

func UI_RecentBookmarksPage(c *gin.Context) {
    c.HTML(http.StatusOK, "recent.html", gin.H{
        "entities": models.GetRecentEntities(),
    })
}

func UI_BookmarksPage(c *gin.Context) {
    c.HTML(http.StatusOK, "bookmarks.html", gin.H{
        "entities": models.GetBookmarks(),
    })
}

func UI_NotePage(c *gin.Context) {
    c.HTML(http.StatusOK, "bookmarks-tag.html", gin.H{
        "entity": models.GetNote(),
    })
}

func UI_BookmarkPage(c *gin.Context) {
    c.HTML(http.StatusOK, "bookmarks-tag.html", gin.H{})
}

func UI_SharedBookmarkPage(c *gin.Context) {
    c.HTML(http.StatusOK, "bookmarks-tag.html", gin.H{})
}

func UI_404ErrorPage(c *gin.Context) {
    c.HTML(http.StatusOK, "404.html", gin.H{})
}
