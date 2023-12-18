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
	}

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

	apiutils.WriteJSON(w, http.StatusOK, dto.ResponseOK{Success: "All"})
	return nil
}
