package models

type Scheduler struct {
	DB *DB
}

func (s *Scheduler) Init(d *DB) {
	s.DB = d
}
