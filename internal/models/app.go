package models

// App represents dependencies which will be injected into the Server handlers.
type App struct {
	Config    *Config
	DB        *DB
	Scheduler *Scheduler
	Settings  Settings
}

// Init sets up some basic parameters of the provided App instance, like the Config
// and the DB (Database) connection.
func (a *App) Init(c string) {
	a.Config = &Config{}
	a.Config.Init(c)
	a.Config.SetURL()

	a.DB = &DB{}
	a.DB.Init(a.Config)

	a.Scheduler = &Scheduler{}
	a.Scheduler.Init(a.DB)

	a.Settings = a.DB.DefaultSettings()
}

func (a *App) GetConfig() *Config {
	return a.Config
}

func (a *App) RegisterEnabled() bool {
	return a.Settings.RegisterEnabled
}
