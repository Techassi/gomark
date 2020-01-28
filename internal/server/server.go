package server

import (
	"fmt"
	"html/template"
	"io"
	"strconv"

	m "github.com/Techassi/gomark/internal/models"
	handle "github.com/Techassi/gomark/internal/server/handlers"
	tpl "github.com/Techassi/gomark/internal/server/templating"
	"github.com/Techassi/gomark/internal/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server is the top-level server instance.
type Server struct {
	Port int
	Mux  *echo.Echo
	App  *m.App
}

// New initiates a new Server instance and returns it.
func New(app *m.App) *Server {
	s := &Server{}
	s.Init(app)

	return s
}

// Init sets up some basic parameters of the provided Server instance.
func (s *Server) Init(app *m.App) {
	s.Port = app.Config.Server.Port
	s.Mux = echo.New()

	// Set debug mode (only for development)
	s.Mux.Debug = true

	// Hide startup message
	s.Mux.HideBanner = true

	f := template.FuncMap{
		"FormatNoImage": tpl.FormatNoImage,
		"FormatColor":   tpl.FormatColor,
	}

	// Register template renderer
	s.Mux.Renderer = &TemplateRenderer{
		Templates: template.Must(template.New("").Funcs(f).ParseGlob(util.GetAbsPath("templates/*/*.html"))),
	}

	// Register middlewares
	s.Mux.Use(middleware.Logger())
	s.Mux.Use(middleware.Recover())

	// Enable GZIP compression
	s.Mux.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	s.App = app
}

// Run spins up the router, adds all routes and listens on the provided port.
// It panics when the router couldn't be started. Panics in the call chain
// will recover, print a stack trace and the HTTPErrorHandler handles the panic.
func (s *Server) Run() {
	// Static routes
	s.Mux.Static("/js", util.GetAbsPath("public/js/dist"))
	s.Mux.Static("/css", util.GetAbsPath("public/scss"))
	s.Mux.Static("/assets", util.GetAbsPath("public/assets"))
	s.Mux.Static("/font", util.GetAbsPath("public/assets/fonts"))

	// Unprotected routes
	s.Mux.GET("/login", handle.UILoginPage)
	s.Mux.GET("/register", handle.UIRegisterPage)
	s.Mux.GET("/s/:hash", handle.UISharedBookmarkPage)

	s.Mux.POST("/login", handle.AuthLogin)
	s.Mux.POST("/register", handle.AuthRegister)

	// Custom 404 error page
	s.Mux.GET("/404", handle.UI404Page)

	s.Mux.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", s.App)
			return next(c)
		}
	})

	// Protected routes
	pr := s.Mux.Group("/")
	pr.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:              []byte(s.App.Config.Security.Jwt.Secret),
		TokenLookup:             "cookie:Authorization",
		ErrorHandlerWithContext: handle.AuthJWTError,
	}))

	pr.GET("", handle.UIHomePage)
	pr.GET("/code", handle.UI2FACodePage)
	pr.GET("/notes", handle.UINotesPage)
	pr.GET("/shared", handle.UISharedPage)
	pr.GET("/recent", handle.UIRecentBookmarksPage)
	pr.GET("/bookmarks", handle.UIBookmarksPage)
	pr.GET("/b/:hash", handle.UIBookmarkPage)
	pr.GET("/n/:hash", handle.UINotePage)

	// API routes
	api := s.Mux.Group("/api")
	api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:              []byte(s.App.Config.Security.Jwt.Secret),
		TokenLookup:             "cookie:Authorization",
		ErrorHandlerWithContext: handle.AuthJWTError,
	}))

	// v1 API routes
	v1 := api.Group("/v1")
	v1.GET("/recent", handle.APIGetRecentBookmarks)
	v1.GET("/bookmarks", handle.APIGetBookmarks)
	v1.GET("/bookmarks/:hash", handle.APIGetBookmark)
	v1.GET("/bookmarks/:hash/tags", handle.APIGetBookmarkTags)
	v1.GET("/folders", handle.APIGetFolders)
	v1.GET("/folders/:hash", handle.APIGetSubFolders)

	v1.POST("/bookmark", handle.APIPostBookmark)
	v1.POST("/bookmark/:hash", handle.APIUpdateBookmark)
	v1.POST("/bookmark/:hash/tags", handle.APIPostBookmarkTags)
	v1.POST("/folder", handle.APIPostFolder)
	v1.POST("/folder/:hash", handle.APIPostEntityToFolder)

	// Auth routes
	auth := s.Mux.Group("/auth")
	auth.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:              []byte(s.App.Config.Security.Jwt.Secret),
		TokenLookup:             "cookie:Authorization",
		ErrorHandlerWithContext: handle.AuthJWTError,
	}))

	auth.POST("/code", handle.Auth2FACode)
	auth.POST("/code/create", handle.AuthCreate2FACode)

	// Startup the router
	port := fmt.Sprintf(":%s", strconv.Itoa(s.Port))
	s.Mux.Logger.Fatal(s.Mux.Start(port))
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	Templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.Templates.ExecuteTemplate(w, name, data)
}
