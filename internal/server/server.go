package server

import (
    "strings"
    "net/http"

    handle "github.com/Techassi/gomark/internal/server/handlers"
    mw "github.com/Techassi/gomark/internal/server/middlewares"
    "github.com/Techassi/gomark/internal/util"

    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
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
	r.Static("/font", util.GetAbsPath("public/assets/fonts"))
	// r.Static("/image", helper.JoinPaths(conf.Paths.WWWDir.Path, "/images"))

    r.Use(sessions.Sessions("gomark", cookie.NewStore([]byte("secret"))))

    // Unprotected frontend routes
    r.GET("/login", handle.UI_LoginPage)
    r.GET("/register", handle.UI_RegisterPage)
    r.GET("/s/:hash", handle.UI_SharedBookmarkPage)

    r.GET("404", handle.UI_404ErrorPage)

    // Protected frontend routes
    protected := r.Group("")
    protected.Use(mw.ValidateSession)
    {
        protected.GET("/", handle.UI_HomePage)
        protected.GET("/notes", handle.UI_NotesPage)
        protected.GET("/shared", handle.UI_SharedPage)
        protected.GET("/recent", handle.UI_RecentBookmarksPage)
        protected.GET("/bookmarks", handle.UI_BookmarksPage)
        protected.GET("/b/:hash", handle.UI_BookmarkPage)
        protected.GET("/n/:hash", handle.UI_NotePage)
    }

    // V1 API routes
    v1 := r.Group("/api/v1")
    {
        v1.GET("/recent", handle.API_GetRecentBookmarks)
        v1.GET("/bookmarks", handle.API_GetBookmarks)
        v1.GET("/bookmarks/:hash", handle.API_GetBookmark)
        v1.GET("/bookmarks/:hash/tags", handle.API_GetBookmarkTags)

        v1.POST("/bookmarks", handle.API_PostBookmark)
        v1.POST("/bookmarks/:hash/tags", handle.API_PostBookmarkTags)
    }

    // Authentication routes
    auth := r.Group("auth")
    {
        auth.GET("/refresh", handle.AUTH_Login)

        auth.POST("/login", handle.AUTH_Login)
        auth.POST("/register", handle.AUTH_Register)
        auth.POST("/logout", handle.AUTH_Logout)
    }

    r.NoRoute(func(c *gin.Context) {
        if strings.Contains(c.Request.Header.Get("Accept"), "text") {
            c.Redirect(http.StatusMovedPermanently, "/404")
            return
        }
        c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

    r.Run(port)
}
