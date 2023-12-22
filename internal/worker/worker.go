package worker

import (
	"github.com/MehmetTalhaSeker/concurrent-web-service/application/taskservice"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/dto"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/logger"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apputils"
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
	service    taskservice.Service
	lg         *logger.Logger
}

func NewWorker(workerPool chan chan Job, service taskservice.Service, lg *logger.Logger) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan struct{}),
		service:    service,
		lg:         lg,
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
					err := apputils.InterfaceToStruct(job.Payload.Data, d)
					if err != nil {
						w.lg.WarningLog.Println(err)
					}

					go func() {
						err := w.service.Create(d)
						if err != nil {
							w.lg.WarningLog.Println(err)
						}
					}()

				case dto.OpPut:
					d := new(dto.TaskUpdateRequest)
					err := apputils.InterfaceToStruct(job.Payload.Data, d)
					if err != nil {
						w.lg.WarningLog.Println(err)
					}

					go func() {
						_, err := w.service.Update(d)
						if err != nil {
							w.lg.WarningLog.Println(err)
						}
					}()

				case dto.OpDelete:
					d := new(dto.RequestWithID)
					err := apputils.InterfaceToStruct(job.Payload.Data, d)
					if err != nil {
						w.lg.WarningLog.Println(err)
					}
					go func() {
						_, err := w.service.Delete(d)
						if err != nil {
							w.lg.WarningLog.Println(err)
						}
					}()
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
