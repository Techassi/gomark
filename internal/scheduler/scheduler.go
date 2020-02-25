package scheduler

import "github.com/Techassi/gomark/internal/db"

type Scheduler struct {
	DB *db.DB
}

func (s *Scheduler) Init(d *db.DB) {
	s.DB = d
}
