package worker

import (
	"github.com/MehmetTalhaSeker/concurrent-web-service/application/taskservice"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/logger"
	"log"
)

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	maxWorkers int
	WorkerPool chan chan Job
	jobQueue   chan Job
	service    taskservice.Service
	lg         *logger.Logger
}

func NewDispatcher(maxWorkers int, jobQueue chan Job, service taskservice.Service, lg *logger.Logger) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool, maxWorkers: maxWorkers, jobQueue: jobQueue, service: service, lg: lg}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.service, d.lg)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	log.Println("Worker que dispatcher started...")

	for {
		select {
		case job := <-d.jobQueue:
			go func(job Job) {
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
