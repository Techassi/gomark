package server

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strconv"

	"github.com/Techassi/gomark/internal/app"
	cnst "github.com/Techassi/gomark/internal/constants"
	tpl "github.com/Techassi/gomark/internal/templating"
	"github.com/Techassi/gomark/internal/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server is the top-level server instance.
type Server struct {
	Port int
	Mux  *echo.Echo
	App  *app.App
}

// New initiates a new Server instance and returns it.
func New(app *app.App) *Server {
	s := &Server{}
	s.Init(app)

	return s
}

// Init sets up some basic parameters of the provided Server instance.
func (s *Server) Init(app *app.App) {
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
		Templates: template.Must(tpl.CompileTemplates(util.AbsolutePath("templates/*/*.html"), f)),
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
	s.Mux.Static("/", util.AbsolutePath("app/dist"))
	s.Mux.Static("/js", util.AbsolutePath("app/dist/js"))
	s.Mux.Static("/css", util.AbsolutePath("app/dist/css"))
	s.Mux.Static("/img", util.AbsolutePath("app/dist/img"))

	// Image and archive routes
	s.Mux.Static("/image", filepath.Join(s.App.Config.WebRoot, cnst.FS_IMAGE_DIR))
	s.Mux.Static("/archive", filepath.Join(s.App.Config.WebRoot, cnst.FS_ARCHIVE_DIR))

	// Auth routes
	s.Mux.POST("/auth/login", s.App.AuthLogin)
	s.Mux.POST("/auth/register", s.App.AuthRegister)

	// Custom 404 error page
	// s.Mux.GET("/404", s.App.Ui404Page)

	// Configure the JWT setting to use in the middleware
	jwtConfig := middleware.JWTConfig{
		SigningKey:              []byte(s.App.Config.Security.Jwt.Secret),
		TokenLookup:             "cookie:Authorization",
		ErrorHandlerWithContext: s.App.AuthJWTError,
	}

	// API routes
	api := s.Mux.Group("/api")
	// api.Use(middleware.JWTWithConfig(jwtConfig))

	// v1 API routes
	v1 := api.Group("/v1")
	v1.GET("/recent", s.App.ApiGetRecentBookmarks)
	v1.GET("/bookmarks", s.App.ApiGetBookmarks)
	v1.GET("/bookmarks/:hash", s.App.ApiGetBookmark)
	v1.GET("/folders", s.App.ApiGetFolders)
	v1.GET("/folders/:hash", s.App.ApiGetSubFolders)

	v1.POST("/bookmark", s.App.ApiPostBookmark)
	v1.POST("/bookmark/:hash", s.App.ApiUpdateBookmark)
	v1.POST("/bookmark/:hash/share", s.App.ApiShareBookmark)
	v1.POST("/folder", s.App.ApiPostFolder)
	v1.POST("/folder/:hash", s.App.ApiPostEntityToFolder)

	// Event endpoint
	v1.POST("/event", s.App.ApiPostEvent)

	// Auth routes
	auth := s.Mux.Group("/auth")
	auth.Use(middleware.JWTWithConfig(jwtConfig))

	auth.POST("/code", s.App.Auth2FACode)
	auth.POST("/code/create", s.App.AuthCreate2FACode)

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
