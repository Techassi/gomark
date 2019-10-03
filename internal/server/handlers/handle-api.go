package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func API_GetRecentBookmarks(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "bookmarks": ""})
}

func API_GetBookmarks(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "bookmarks": ""})
}

func API_GetBookmark(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "bookmarks": ""})
}

func API_GetBookmarkTags(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "tags": ""})
}

func API_PostBookmark(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "action": ""})
}

func API_PostBookmarkTags(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "action": ""})
}
