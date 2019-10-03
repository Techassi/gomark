package server

import (
    handle "github.com/Techassi/gomark/internal/server/handlers"
    "github.com/Techassi/gomark/internal/util"

    "github.com/gin-gonic/gin"
)

func Startup(port string) {
    r := gin.Default()

    // Load templates
    r.LoadHTMLGlob(util.GetAbsPath("templates/*/*.html"))

    // Static routes
    r.Static("/style", util.GetAbsPath("public/scss"))
	r.Static("/js", util.GetAbsPath("public/js/dist"))
	r.Static("/favicon", util.GetAbsPath("public/assets/favicon"))
	// r.Static("/image", helper.JoinPaths(conf.Paths.WWWDir.Path, "/images"))

    // Frontend routes
    r.GET("/", handle.UI_HomePage)
    r.GET("/login", handle.UI_LoginPage)
    r.GET("/recent", handle.UI_RecentBookmarksPage)
    r.GET("/bookmarks", handle.UI_BookmarksPage)
    r.GET("/bookmarks/:hash", handle.UI_BookmarkPage)

    // V1 API routes
    v1 := r.Group("/api/v1")
    {
        v1.GET("recent", handle.API_GetRecentBookmarks)
        v1.GET("bookmarks", handle.API_GetBookmarks)
        v1.GET("/bookmarks/:hash", handle.API_GetBookmark)
        v1.GET("/bookmarks/:hash/tags", handle.API_GetBookmarkTags)

        v1.POST("/bookmarks", handle.API_PostBookmark)
        v1.POST("/bookmarks/:hash/tags", handle.API_PostBookmarkTags)
    }

    // Authentication routes
    auth := r.Group("auth")
    {
        auth.POST("/login", handle.AUTH_Login)
        auth.POST("/logout", handle.AUTH_Logout)
    }

    r.Run(port)
}
