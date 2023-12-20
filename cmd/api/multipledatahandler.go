package main

import (
	"encoding/json"
	"fmt"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/dto"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apiutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/worker"
	"io"
	"net/http"
	"sync"
)

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
func (s *Server) processMultipleTasks(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		apiutils.WriteJSON(w, http.StatusMethodNotAllowed, errorutils.New(errorutils.ErrMethodNotAllowed, nil))
		return nil
	}

	content := new(dto.PayloadCollection)
	if err := json.NewDecoder(io.LimitReader(r.Body, 2048)).Decode(&content); err != nil {
		return err
	}

	if err := s.validator.Validate(content); err != nil {
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
			if err := s.processTask(payload); err != nil {
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

func (s *Server) processTask(payload dto.Payload) *errorutils.APIError {
	if err := s.validator.Validate(payload); err != nil {
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

	if err := s.validator.Validate(taskRequest); err != nil {
		return errorutils.New(fmt.Errorf("task:%+v has errors: %v", payload.Data, err), nil)
	}

	go func() {
		w := worker.Job{Payload: payload}
		s.jobQueue <- w
	}()

	return nil
}
