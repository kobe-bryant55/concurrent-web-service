package main

import (
	"database/sql"
	"encoding/json"
	"github.com/MehmetTalhaSeker/concurrent-web-service/application/task"
	taskdomain "github.com/MehmetTalhaSeker/concurrent-web-service/domain/task"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/errorutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/shared/config"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	addr string
	cfg  *config.Config
	db   *sql.DB
}

func NewServer(addr string, cfg *config.Config, db *sql.DB) *Server {
	return &Server{
		addr: addr,
		cfg:  cfg,
		db:   db,
	}
}

func (s *Server) Run() {
	mux := http.NewServeMux()

	tr := taskdomain.NewPostgresTaskRepository(s.db)
	ts := task.NewService(tr)
	th := newTaskHandler(ts)

	mux.Handle("/tasks/", th)
	mux.Handle("/tasks", th)

	log.Println("Server running on: ", s.addr)
	http.ListenAndServe(s.addr, mux)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")

	jb, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(status)
	w.Write(jb)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func errorHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			status, errs := errorutils.Handler(err)
			writeJSON(w, status, errs)
		}
	}
}

func getID(r *http.Request, prefix string) string {
	id := strings.TrimPrefix(r.URL.Path, prefix)

	return id
}
