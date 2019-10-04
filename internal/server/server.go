package server

import (
    "log"
    "time"
    "net/http"

    handle "github.com/Techassi/gomark/internal/server/handlers"
    mw "github.com/Techassi/gomark/internal/server/middlewares"
    m "github.com/Techassi/gomark/internal/server/models"
    "github.com/Techassi/gomark/internal/util"

    "github.com/appleboy/gin-jwt/v2"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
)

var identityKey = "id"

func Startup(db *gorm.DB, port string) {
    r := gin.Default()

    // Load templates
    r.LoadHTMLGlob(util.GetAbsPath("templates/*/*.html"))

    // Static routes
    r.Static("/style", util.GetAbsPath("public/scss"))
	r.Static("/js", util.GetAbsPath("public/js/dist"))
	r.Static("/favicon", util.GetAbsPath("public/assets/favicon"))
	r.Static("/font", util.GetAbsPath("public/assets/fonts"))
	// r.Static("/image", helper.JoinPaths(conf.Paths.WWWDir.Path, "/images"))

    // Inject Database Connection into handlers
	r.Use(mw.InjectDB(db))

    authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "Gomark",
		Key:         []byte("JcExCbVJC>Jc4vLcSBG13l4TF2ZayRaXfWF18NaLR!k87fGPR9t2wVGVKWp3k5VA"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*m.User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &m.User{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals m.Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

            // rewrite this to check if empty
			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				return &m.User{
					UserName:  userID,
					LastName:  "Bo-Yi",
					FirstName: "Wu",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
            // check for account in database
			if v, ok := data.(*m.User); ok && v.UserName == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.Redirect(http.StatusMovedPermanently, "/login")
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

    // Unprotected frontend routes
    r.GET("/login", handle.UI_LoginPage)
    r.GET("/register", handle.UI_RegisterPage)
    r.GET("/s/:hash", handle.UI_SharedBookmarkPage)

    // Protected frontend routes
    protected := r.Group("")
    protected.Use(authMiddleware.MiddlewareFunc())
    {
        protected.GET("/", handle.UI_HomePage)
        protected.GET("/shared", handle.UI_SharedPage)
        protected.GET("/recent", handle.UI_RecentBookmarksPage)
        protected.GET("/bookmarks", handle.UI_BookmarksPage)
        protected.GET("/b/:hash", handle.UI_BookmarkPage)
    }

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
        auth.GET("/refresh", authMiddleware.RefreshHandler)

        auth.POST("/login", authMiddleware.LoginHandler)
        auth.POST("/register", handle.AUTH_Register)
        auth.POST("/logout", handle.AUTH_Logout)
    }

    r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

    r.Run(port)
}
