package main

import (
	"database/sql"
	"github.com/MehmetTalhaSeker/concurrent-web-service/application/task"
	taskdomain "github.com/MehmetTalhaSeker/concurrent-web-service/domain/task"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/shared/config"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/validatorutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/worker"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"strings"
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

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

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
