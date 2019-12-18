package server

import (
	"fmt"
	"html/template"
	"io"
	"strconv"

	m "github.com/Techassi/gomark/internal/models"
	handle "github.com/Techassi/gomark/internal/server/handlers"
	"github.com/Techassi/gomark/internal/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server is the top-level instance.
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

	// Register template renderer
	s.Mux.Renderer = &TemplateRenderer{
		Templates: template.Must(template.ParseGlob(util.GetAbsPath("templates/*/*.html"))),
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

// Startup spins up the router, adds all routes and listens on the provided port.
// It panics when the router couldn't be started. Panics in the call chain
// will recover, print a stack trace and the HTTPErrorHandler handles the panic.
func (s *Server) Run() {
	// Static routes
	s.Mux.Static("/static", util.GetAbsPath("public"))
	s.Mux.Static("/2fa", util.GetAbsPath("public/2fa"))

	// Unprotected routes
	s.Mux.GET("/login", handle.UI_LoginPage)
	s.Mux.GET("/code", handle.UI_TwoFACodePage)
	s.Mux.GET("/register", handle.UI_RegisterPage)
	s.Mux.GET("/s/:hash", handle.UI_SharedBookmarkPage)

	// Custom 404 error page
	s.Mux.GET("/404", handle.UI_404Page)

	s.Mux.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", s.App)
			return next(c)
		}
	})

	// Protected routes
	pr := s.Mux.Group("")
	pr.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:              []byte(s.App.Config.Security.Jwt.Secret),
		TokenLookup:             "cookie:Authorization",
		ErrorHandlerWithContext: handle.AUTH_JWTError,
	}))

	pr.GET("/", handle.UI_DashboardPage)
	pr.GET("/notes", handle.UI_NotesPage)
	pr.GET("/shared", handle.UI_SharedPage)
	pr.GET("/recent", handle.UI_RecentBookmarksPage)
	pr.GET("/bookmarks", handle.UI_BookmarksPage)
	pr.GET("/b/:hash", handle.UI_BookmarkPage)
	pr.GET("/n/:hash", handle.UI_NotePage)

	// API routes
	api := s.Mux.Group("/api")

	// v1 API routes
	v1 := api.Group("/v1")
	v1.GET("/recent", handle.API_GetRecentBookmarks)
	v1.GET("/bookmarks", handle.API_GetBookmarks)
	v1.GET("/bookmarks/:hash", handle.API_GetBookmark)
	v1.GET("/bookmarks/:hash/tags", handle.API_GetBookmarkTags)

	v1.POST("/bookmarks", handle.API_PostBookmark)
	v1.POST("/bookmarks/:hash/tags", handle.API_PostBookmarkTags)

	// Auth routes
	auth := s.Mux.Group("/auth")
	auth.POST("/login", handle.AUTH_JWTLogin)
	auth.POST("/logout", handle.AUTH_JWTLogout)
	auth.POST("/register", handle.AUTH_JWTRegister)

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
