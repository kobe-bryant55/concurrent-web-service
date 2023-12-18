package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/MehmetTalhaSeker/concurrent-web-service/application/task"
	taskdomain "github.com/MehmetTalhaSeker/concurrent-web-service/domain/task"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/dto"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/shared/config"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apiutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/validatorutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/worker"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Server struct {
	addr      string
	cfg       *config.Config
	db        *sql.DB
	validator validatorutils.Validator
	jobQueue  chan worker.Job
}

func NewServer(addr string, cfg *config.Config, db *sql.DB, jobQueue chan worker.Job) *Server {
	return &Server{
		addr:      addr,
		cfg:       cfg,
		db:        db,
		validator: validatorutils.NewValidator(),
		jobQueue:  jobQueue,
	}
}

func (s *Server) Run() {
	mux := http.NewServeMux()

	tr := taskdomain.NewPostgresTaskRepository(s.db)
	ts := task.NewService(tr)
	th := newTaskHandler(ts)

	// Create and Start dispatcher.
	dispatcher := worker.NewDispatcher(worker.MaxWorker, s.jobQueue, ts)
	dispatcher.Run()

	mux.Handle("/tasks/", th)
	mux.Handle("/tasks", th)
	mux.Handle("/multiple/", errorHandler(s.processMultipleTasks))
	mux.Handle("/multiple", errorHandler(s.processMultipleTasks))

	log.Println("Server running on: ", s.addr)
	http.ListenAndServe(s.addr, mux)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func errorHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			errorutils.Handler(err, w)
		}
	}
}

func getID(r *http.Request, prefix string) string {
	id := strings.TrimPrefix(r.URL.Path, prefix)

	return id
}

func (s *Server) processMultipleTasks(w http.ResponseWriter, r *http.Request) error {
	content := new(dto.PayloadCollection)

	if err := json.NewDecoder(io.LimitReader(r.Body, worker.MaxLength)).Decode(&content); err != nil {
		return err
	}

	if err := s.validator.Validate(content); err != nil {
		return err
	}

	var errs []*errorutils.APIError
	errCh := make(chan *errorutils.APIError, len(content.Payloads))

	wg := sync.WaitGroup{}
	wg.Add(len(content.Payloads))

	for _, payload := range content.Payloads {
		payload := payload
		go func() {
			if err := s.validator.Validate(payload); err != nil {
				errCh <- errorutils.New(fmt.Errorf("task:%+v has errors: %v", payload, err), nil)
				wg.Done()
				return
			}

			switch payload.OperationType {

			case dto.OpCreate:
				d := new(dto.TaskCreateRequest)
				err := apputils.InterfaceToStruct(payload.Data, d)
				if err != nil {
					errCh <- errorutils.New(fmt.Errorf("task:%+v has errors: %v", payload.Data, err), nil)
					break
				}

				if err = s.validator.Validate(d); err != nil {
					errCh <- errorutils.New(fmt.Errorf("task:%+v has errors: %v", payload.Data, err), nil)
					break
				}

				w := worker.Job{Payload: payload}
				s.jobQueue <- w
				break

			case dto.OpPut:
				d := new(dto.TaskUpdateRequest)
				err := apputils.InterfaceToStruct(payload.Data, d)
				if err != nil {
					errCh <- errorutils.New(fmt.Errorf("task:%+v has errors: %v", payload.Data, err), nil)
					break
				}

				if err = s.validator.Validate(d); err != nil {
					errCh <- errorutils.New(fmt.Errorf("task:%+v has errors: %v", payload.Data, err), nil)
					break
				}

				w := worker.Job{Payload: payload}
				s.jobQueue <- w
				break

			case dto.OpDelete:
				d := new(dto.RequestWithID)
				err := apputils.InterfaceToStruct(payload.Data, d)
				if err != nil {
					errCh <- errorutils.New(fmt.Errorf("task:%+v has errors: %v", payload.Data, err), nil)
					break
				}

				if err = s.validator.Validate(d); err != nil {
					errCh <- errorutils.New(fmt.Errorf("task:%+v has errors: %v", payload.Data, err), nil)
					break
				}

				w := worker.Job{Payload: payload}
				s.jobQueue <- w
				break

			default:
				errCh <- errorutils.New(fmt.Errorf("invalid operation for: %+v", payload), nil)
				break
			}

			wg.Done()
		}()
	}

	wg.Wait()

	close(errCh)

	for err := range errCh {
		errs = append(errs, err)
	}

	if errs != nil {
		return errorutils.ValidationError(errs)
	}

	apiutils.WriteJSON(w, http.StatusOK, "OK")
	return nil
}
