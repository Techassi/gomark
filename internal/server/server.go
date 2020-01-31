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
	s.Mux.GET("/login", handle.UI_LoginPage)
	s.Mux.GET("/register", handle.UI_RegisterPage)
	s.Mux.GET("/s/:hash", handle.UI_SharedBookmarkPage)

	s.Mux.POST("/login", handle.AUTH_Login)
	s.Mux.POST("/register", handle.AUTH_Register)

	// Custom 404 error page
	s.Mux.GET("/404", handle.UI_404Page)

	s.Mux.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", s.App)
			return next(c)
		}
	})

	// Configure the JWT setting to use in the middleware
	jwtConfig := middleware.JWTConfig{
		SigningKey:              []byte(s.App.Config.Security.Jwt.Secret),
		TokenLookup:             "cookie:Authorization",
		ErrorHandlerWithContext: handle.AUTH_JWTError,
	}

	// Protected routes
	pr := s.Mux.Group("/")
	pr.Use(middleware.JWTWithConfig(jwtConfig))

	pr.GET("", handle.UI_HomePage)
	pr.GET("/code", handle.UI_2FACodePage)
	pr.GET("/notes", handle.UI_NotesPage)
	pr.GET("/shared", handle.UI_SharedPage)
	pr.GET("/recent", handle.UI_RecentBookmarksPage)
	pr.GET("/bookmarks", handle.UI_BookmarksPage)
	pr.GET("/b/:hash", handle.UI_BookmarkPage)
	pr.GET("/n/:hash", handle.UI_NotePage)

	// API routes
	api := s.Mux.Group("/api")
	api.Use(middleware.JWTWithConfig(jwtConfig))

	// v1 API routes
	v1 := api.Group("/v1")
	v1.GET("/recent", handle.API_GetRecentBookmarks)
	v1.GET("/bookmarks", handle.API_GetBookmarks)
	v1.GET("/bookmarks/:hash", handle.API_GetBookmark)
	v1.GET("/bookmarks/:hash/tags", handle.API_GetBookmarkTags)
	v1.GET("/folders", handle.API_GetFolders)
	v1.GET("/folders/:hash", handle.API_GetSubFolders)

	v1.POST("/bookmark", handle.API_PostBookmark)
	v1.POST("/bookmark/:hash", handle.API_UpdateBookmark)
	v1.POST("/bookmark/:hash/tags", handle.API_PostBookmarkTags)
	v1.POST("/folder", handle.API_PostFolder)
	v1.POST("/folder/:hash", handle.API_PostEntityToFolder)

	// Auth routes
	auth := s.Mux.Group("/auth")
	auth.Use(middleware.JWTWithConfig(jwtConfig))

	auth.POST("/code", handle.AUTH_2FACode)
	auth.POST("/code/create", handle.AUTH_Create2FACode)

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
