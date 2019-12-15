package models

// App represents dependencies which will be injected into the Server handlers.
type App struct {
	Config *Config
	DB     *DB
}

// Init sets up some basic parameters of the provided App instance, like the Config
// and the DB (Database) connection.
func (a *App) Init(c string) {
	a.Config = &Config{}
	a.Config.Init(c)

	a.DB = &DB{}
	a.DB.Init(a.Config)
}
