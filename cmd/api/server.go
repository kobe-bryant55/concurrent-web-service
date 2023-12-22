package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/MehmetTalhaSeker/concurrent-web-service/application/taskservice"
	taskdomain "github.com/MehmetTalhaSeker/concurrent-web-service/domain/task"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/logger"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/rba"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/shared/config"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/validatorutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/worker"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	addr      string
	cfg       *config.Config
	db        *sql.DB
	validator validatorutils.Validator
	jobQueue  chan worker.Job
	rba       rba.RBA
	lg        *logger.Logger
}

func NewServer(addr string, cfg *config.Config, db *sql.DB, jobQueue chan worker.Job, rba rba.RBA, vl validatorutils.Validator, lg *logger.Logger) *Server {
	return &Server{
		addr:      addr,
		cfg:       cfg,
		db:        db,
		jobQueue:  jobQueue,
		rba:       rba,
		validator: vl,
		lg:        lg,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()

	tr := taskdomain.NewPostgresTaskRepository(s.db)
	ts := taskservice.NewService(tr)
	th := newTaskHandler(ts, s.validator, s.jobQueue, s.rba)

	ah := newAuthHandler()

	lh := newLogHandler(s.lg)

	// Create and Start dispatcher.
	dispatcher := worker.NewDispatcher(s.cfg.Worker.MaxWorker, s.jobQueue, ts, s.lg)
	dispatcher.Run()

	mux.Handle("/tasks/", lh.log(authenticate(th)))
	mux.Handle("/tasks", lh.log(authenticate(th)))
	mux.Handle("/auth", lh.log(ah))
	mux.Handle("/auth/", lh.log(ah))

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server running on: ", s.addr)

	return http.ListenAndServe(s.addr, mux)
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
