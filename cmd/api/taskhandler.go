package main

import (
	"encoding/json"
	"fmt"
	"github.com/MehmetTalhaSeker/concurrent-web-service/application/taskservice"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/dto"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/rba"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/types"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apiutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/validatorutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/worker"
	"io"
	"net/http"
	"regexp"
	"sync"
)

var (
	multiple   = regexp.MustCompile(`^/tasks/multiple/*$`)
	taskCreate = regexp.MustCompile(`^/tasks/*$`)
	taskReads  = regexp.MustCompile(`^/tasks/*$`)
	taskRead   = regexp.MustCompile(`^/tasks/(\d+)$`)
	taskUpdate = regexp.MustCompile(`^/tasks/(\d+)$`)
	taskDelete = regexp.MustCompile(`^/tasks/(\d+)$`)
)

type taskHandler struct {
	service   taskservice.Service
	validator validatorutils.Validator
	queue     chan worker.Job
	rba       rba.RBA
}

func newTaskHandler(service taskservice.Service, validator validatorutils.Validator, queue chan worker.Job, rba rba.RBA) *taskHandler {
	return &taskHandler{
		service:   service,
		validator: validator,
		queue:     queue,
		rba:       rba,
	}
}

func (h *taskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && multiple.MatchString(r.URL.Path):
		errorHandler(h.processMultipleTasks)(w, r)
		return

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
	if err := h.rba.CheckHasRole(r.Context(), types.Admin); err != nil {
		return err
	}

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
	if err := h.rba.CheckHasRole(r.Context(), types.Admin); err != nil {
		return err
	}

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
	if err := h.rba.CheckHasRole(r.Context(), types.Admin); err != nil {
		return err
	}

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

// processMultipleTasks godoc
// @Summary Process Multiple Tasks
// @Description Create, update and delete tasks with single request.
// @Tags tags
// @Accept  json
// @Param processMultipleTasks body dto.PayloadCollection true "Input"
// @Produce  json
// @Success 200 {object} dto.ResponseOK
// @Failure      400  {object}  errorutils.APIErrors
// @Failure      422  {object}  errorutils.APIErrors
// @Failure      404  {object}  errorutils.APIErrors
// @Failure      500  {object}  errorutils.APIErrors
// @Router /multiple [post]
func (h *taskHandler) processMultipleTasks(w http.ResponseWriter, r *http.Request) error {
	if err := h.rba.CheckHasRole(r.Context(), types.Admin); err != nil {
		return err
	}

	content := new(dto.PayloadCollection)
	if err := json.NewDecoder(io.LimitReader(r.Body, 2048)).Decode(&content); err != nil {
		return err
	}

	if err := h.validator.Validate(content); err != nil {
		return err
	}

	var errs []*errorutils.APIError
	var wg sync.WaitGroup
	errCh := make(chan *errorutils.APIError, len(content.Payloads))

	for _, payload := range content.Payloads {
		payload := payload
		wg.Add(1)
		go func(payload dto.Payload) {
			defer wg.Done()
			if err := h.processTask(payload); err != nil {
				errCh <- err
			}
		}(payload)
	}

	wg.Wait()

	close(errCh)

	for err := range errCh {
		errs = append(errs, err)
	}

	if errs != nil {
		return errorutils.ValidationError(errs)
	}

	apiutils.WriteJSON(w, http.StatusOK, dto.ResponseOK{Success: "all"})
	return nil
}

func (h *taskHandler) processTask(payload dto.Payload) *errorutils.APIError {
	if err := h.validator.Validate(payload); err != nil {
		return errorutils.New(fmt.Errorf("task:%+v has errors: %v", payload, err), nil)
	}

	var taskRequest interface{}

	switch payload.OperationType {
	case dto.OpCreate:
		taskRequest = new(dto.TaskCreateRequest)
	case dto.OpPut:
		taskRequest = new(dto.TaskUpdateRequest)
	case dto.OpDelete:
		taskRequest = new(dto.RequestWithID)
	default:
		return errorutils.New(fmt.Errorf("invalid operation for: %+v", payload), nil)
	}

	if err := apputils.InterfaceToStruct(payload.Data, taskRequest); err != nil {
		return errorutils.New(fmt.Errorf("task:%+v has errors: %v", payload.Data, err), nil)
	}

	if err := h.validator.Validate(taskRequest); err != nil {
		return errorutils.New(fmt.Errorf("task:%+v has errors: %v", payload.Data, err), nil)
	}

	go func() {
		w := worker.Job{Payload: payload}
		h.queue <- w
	}()

	return nil
}
