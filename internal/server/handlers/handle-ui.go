package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func UI_HomePage(c *gin.Context) {
    c.HTML(http.StatusOK, "home.html", gin.H{})
}

func UI_LoginPage(c *gin.Context) {
    c.HTML(http.StatusOK, "login.html", gin.H{})
}

func UI_RegisterPage(c *gin.Context) {
    c.HTML(http.StatusOK, "register.html", gin.H{})
}

func UI_NotesPage(c *gin.Context) {
    c.HTML(http.StatusOK, "login.html", gin.H{})
}

func UI_SharedPage(c *gin.Context) {
    c.HTML(http.StatusOK, "login.html", gin.H{})
}

func UI_RecentBookmarksPage(c *gin.Context) {
    c.HTML(http.StatusOK, "recent.html", gin.H{})
}

func UI_BookmarksPage(c *gin.Context) {
    c.HTML(http.StatusOK, "bookmarks.html", gin.H{})
}

func UI_NotePage(c *gin.Context) {
    c.HTML(http.StatusOK, "bookmarks-tag.html", gin.H{})
}

func UI_BookmarkPage(c *gin.Context) {
    c.HTML(http.StatusOK, "bookmarks-tag.html", gin.H{})
}

func UI_SharedBookmarkPage(c *gin.Context) {
    c.HTML(http.StatusOK, "bookmarks-tag.html", gin.H{})
}
