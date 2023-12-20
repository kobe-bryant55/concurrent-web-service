package main

import (
	"encoding/json"
	"github.com/MehmetTalhaSeker/concurrent-web-service/application/taskservice"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/dto"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apiutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
	"net/http"
	"regexp"
)

var (
	taskCreate = regexp.MustCompile(`^/tasks/*$`)
	taskReads  = regexp.MustCompile(`^/tasks/*$`)
	taskRead   = regexp.MustCompile(`^/tasks/(\d+)$`)
	taskUpdate = regexp.MustCompile(`^/tasks/(\d+)$`)
	taskDelete = regexp.MustCompile(`^/tasks/(\d+)$`)
)

type taskHandler struct {
	service taskservice.Service
}

func newTaskHandler(service taskservice.Service) *taskHandler {
	return &taskHandler{service: service}
}

func (h *taskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && taskCreate.MatchString(r.URL.Path):
		errorHandler(h.create)(w, r)
		return

	case r.Method == http.MethodGet && taskReads.MatchString(r.URL.Path):
		errorHandler(h.list)(w, r)
		return

	case r.Method == http.MethodGet && taskRead.MatchString(r.URL.Path):
		errorHandler(h.read)(w, r)
		return

	case r.Method == http.MethodPut && taskUpdate.MatchString(r.URL.Path):
		errorHandler(h.update)(w, r)
		return

	case r.Method == http.MethodDelete && taskDelete.MatchString(r.URL.Path):
		errorHandler(h.delete)(w, r)
		return

	default:
		apiutils.WriteJSON(w, http.StatusMethodNotAllowed, errorutils.New(errorutils.ErrMethodNotAllowed, nil))
		return
	}
}

func (h *taskHandler) create(w http.ResponseWriter, r *http.Request) error {
	d := new(dto.TaskCreateRequest)

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		return err
	}

	err := h.service.Create(d)
	if err != nil {
		return err
	}

	apiutils.WriteJSON(w, http.StatusOK, "OK")
	return nil
}

func (h *taskHandler) list(w http.ResponseWriter, r *http.Request) error {
	l, err := h.service.Reads()
	if err != nil {
		return err
	}

	apiutils.WriteJSON(w, http.StatusOK, l)

	return nil
}

func (h *taskHandler) read(w http.ResponseWriter, r *http.Request) error {
	d := new(dto.RequestWithID)
	id := getID(r, "/tasks/")
	d.ID = id

	l, err := h.service.Read(d)
	if err != nil {
		return err
	}

	apiutils.WriteJSON(w, http.StatusOK, l)

	return nil
}

func (h *taskHandler) update(w http.ResponseWriter, r *http.Request) error {
	d := new(dto.TaskUpdateRequest)
	id := getID(r, "/tasks/")

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		return err
	}

	d.ID = id

	res, err := h.service.Update(d)
	if err != nil {
		return err
	}

	apiutils.WriteJSON(w, http.StatusOK, res)
	return nil
}

func (h *taskHandler) delete(w http.ResponseWriter, r *http.Request) error {
	d := new(dto.RequestWithID)
	id := getID(r, "/tasks/")
	d.ID = id

	res, err := h.service.Delete(d)
	if err != nil {
		return err
	}

	apiutils.WriteJSON(w, http.StatusOK, res)
	return nil
}
