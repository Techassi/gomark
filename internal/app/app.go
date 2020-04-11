package app

import (
	// "github.com/Techassi/gomark/internal/authentication"
	"github.com/Techassi/gomark/internal/db"
	m "github.com/Techassi/gomark/internal/models"
	"github.com/Techassi/gomark/internal/scheduler"
)

// App represents dependencies which will be injected into the Server handlers.
type App struct {
	Config    *m.Config
	DB        *db.DB
	Scheduler *scheduler.Scheduler
	// Authenticator *authentication.Authenticator
	Settings m.Settings
}

// New initiates a new App instance and returns it.
func New(c string) *App {
	a := &App{}
	a.Init(c)

	return a
}

// Init sets up some basic parameters of the provided App instance, like the Config
// and the DB (Database) connection.
func (a *App) Init(c string) {
	a.Config = &m.Config{}
	a.Config.Init(c)
	a.Config.SetURL()

	a.DB = &db.DB{}
	a.DB.Init(a.Config)

	a.Scheduler = scheduler.New(a.Config, a.DB, 2)
	a.Scheduler.RegisterTasks(map[string]func(scheduler.Job){
		"download-meta":  scheduler.HandleMetaDownload,
		"download-image": scheduler.HandleImageDownload,
		"archive":        scheduler.HandleArchive,
		"save":           scheduler.HandleSave,
	})
	a.Scheduler.Start()

	a.Settings = a.DB.DefaultSettings()
}

func (a *App) GetConfig() *m.Config {
	return a.Config
}

func (a *App) RegisterEnabled() bool {
	return a.Settings.RegisterEnabled
}
