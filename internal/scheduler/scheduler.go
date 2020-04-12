package scheduler

import (
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/Techassi/gomark/internal/db"
	m "github.com/Techassi/gomark/internal/models"

	"github.com/pkg/errors"
)

// Scheduler is the top level struct
type Scheduler struct {
	Config         *m.Config
	DB             *db.DB
	Tasks          map[string]func(Job)
	Queue          chan Job
	Ready          chan chan Job
	Workers        []*Worker
	MaxWorkers     int
	HTTPClient     http.Client
	workersStopped sync.WaitGroup
	dispatcherStop sync.WaitGroup
	stop           chan bool
}

// Job represents the work each worker consumes
type Job struct {
	Work      string
	Data      string
	Entity    m.Entity
	Result    Result
	Archive   Archive
	Scheduler *Scheduler
}

type Result struct {
	Title       string
	Description string
	Image       string
}

type Archive struct {
	Body []byte
}

var (
	// ErrorInavlidMaxWorkers will be thrown if the amount of max workers is < -1
	ErrorInavlidMaxWorkers = errors.New("Inavlid max workers amount")

	ErrorNoTaskFound = errors.New("Did not find task with this name")
)

// New returns a new scheduler instance with n workers, provide n = -1 to use
// runtime.NumCPU - 1 workers.
func New(config *m.Config, db *db.DB, n int) *Scheduler {
	if n < -1 {
		n = 1
	}

	if n == -1 || n > runtime.NumCPU() {
		n = runtime.NumCPU() - 1
	}

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	return &Scheduler{
		Config:         config,
		DB:             db,
		Queue:          make(chan Job),
		Ready:          make(chan chan Job, n),
		MaxWorkers:     n,
		HTTPClient:     client,
		workersStopped: sync.WaitGroup{},
		dispatcherStop: sync.WaitGroup{},
		stop:           make(chan bool),
	}
}

func (s *Scheduler) RegisterTasks(t map[string]func(Job)) {
	s.Tasks = t
}

func (s *Scheduler) Start() {
	w := make([]*Worker, s.MaxWorkers, s.MaxWorkers)
	for i := 0; i < s.MaxWorkers; i++ {
		w[i] = NewWorker(s.Ready, s.workersStopped, s)
		w[i].Start()
	}

	s.Workers = w
	go s.dispatch()
}

func (s *Scheduler) Job(work, data string) Job {
	return Job{
		Work:      work,
		Data:      data,
		Scheduler: s,
	}
}

// Schedule schedules a new job in the queue to be cunsumed by workers
func (s *Scheduler) Schedule(job Job) {
	s.Queue <- job
}

func (s *Scheduler) GetTask(work string) (func(Job), error) {
	if f, ok := s.Tasks[work]; ok {
		return f, nil
	}

	return nil, ErrorNoTaskFound
}

func (s *Scheduler) dispatch() {
	s.dispatcherStop.Add(1)
	for {
		select {
		case job := <-s.Queue:
			w := <-s.Ready
			w <- job
		case <-s.stop:
			for _, worker := range s.Workers {
				worker.Stop()
			}
			s.workersStopped.Wait()
			s.dispatcherStop.Done()
			return
		}
	}
}
