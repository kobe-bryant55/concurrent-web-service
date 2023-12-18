package worker

import (
	"github.com/MehmetTalhaSeker/concurrent-web-service/application/task"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/dto"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apputils"
)

var (
	MaxWorker       = 3  //os.Getenv("MAX_WORKERS")
	MaxQueue        = 20 //os.Getenv("MAX_QUEUE")
	MaxLength int64 = 2048
)

// Job represents the job to be run
type Job struct {
	Payload dto.Payload
}

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan struct{}
	service    task.Service
}

func NewWorker(workerPool chan chan Job, service task.Service) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan struct{}),
		service:    service,
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				switch job.Payload.OperationType {
				case dto.OpCreate:
					d := new(dto.TaskCreateRequest)
					apputils.InterfaceToStruct(job.Payload.Data, d)
					go w.service.Create(d)

				case dto.OpPut:
					d := new(dto.TaskUpdateRequest)
					apputils.InterfaceToStruct(job.Payload.Data, d)
					go w.service.Update(d)

				case dto.OpDelete:
					d := new(dto.RequestWithID)
					apputils.InterfaceToStruct(job.Payload.Data, d)
					go w.service.Delete(d)
				}

			case <-w.quit:
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- struct{}{}
	}()
}
