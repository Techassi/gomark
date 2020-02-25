package app

import (
	"github.com/Techassi/gomark/internal/db"
	m "github.com/Techassi/gomark/internal/models"
	"github.com/Techassi/gomark/internal/scheduler"
)

// App represents dependencies which will be injected into the Server handlers.
type App struct {
	Config    *m.Config
	DB        *db.DB
	Scheduler *scheduler.Scheduler
	Settings  m.Settings
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

	a.Scheduler = &scheduler.Scheduler{}
	a.Scheduler.Init(a.DB)

	a.Settings = a.DB.DefaultSettings()
}

func (a *App) GetConfig() *m.Config {
	return a.Config
}

func (a *App) RegisterEnabled() bool {
	return a.Settings.RegisterEnabled
}
