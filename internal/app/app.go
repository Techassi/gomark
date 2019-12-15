package app

import (
	m "github.com/Techassi/gomark/internal/models"
)

// New initiates a new App instance and returns it.
func New(c string) *m.App {
	a := &m.App{}
	a.Init(c)

	return a
}
