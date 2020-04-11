package scheduler

import (
	"sync"
)

type Worker struct {
	Queue     chan Job
	Ready     chan chan Job
	Scheduler *Scheduler
	done      sync.WaitGroup
	stop      chan bool
}

func NewWorker(ready chan chan Job, done sync.WaitGroup, s *Scheduler) *Worker {
	return &Worker{
		Queue:     make(chan Job),
		Ready:     ready,
		Scheduler: s,
		done:      done,
		stop:      make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		w.done.Add(1)
		for {
			w.Ready <- w.Queue
			select {
			case job := <-w.Queue:
				f, err := w.Scheduler.GetTask(job.Work)
				if err != nil {
					continue
				}

				f(job)
			case <-w.stop:
				w.done.Done()
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.stop <- true
}
