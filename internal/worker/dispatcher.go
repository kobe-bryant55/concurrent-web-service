package worker

import (
	"github.com/MehmetTalhaSeker/concurrent-web-service/application/task"
	"log"
)

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	maxWorkers int
	WorkerPool chan chan Job
	jobQueue   chan Job
	service    task.Service
}

func NewDispatcher(maxWorkers int, jobQueue chan Job, service task.Service) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool, maxWorkers: maxWorkers, jobQueue: jobQueue, service: service}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.service)
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
