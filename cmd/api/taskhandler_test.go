package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/appcontext"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/rba"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/types"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/testutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/validatorutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/worker"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/dto"
)

func TestTaskHandler_processMultipleTasks(t *testing.T) {
	ms := &testutils.MockService{}
	validator := validatorutils.NewValidator()
	jobQueue := make(chan worker.Job, 10)
	rb := rba.New()

	ctx := appcontext.WithRole(context.Background(), types.Admin)

	th := newTaskHandler(ms, validator, jobQueue, rb)

	testCases := []struct {
		Name             string
		RequestBody      []dto.Payload
		ExpectedHTTPCode int
		ExpectedResponse string
	}{
		{
			Name: "Valid Payloads",
			RequestBody: []dto.Payload{
				{
					OperationType: dto.OpCreate,
					Data: map[string]interface{}{
						"title":       "Task Title",
						"description": "Task Description",
					},
				},
				{
					OperationType: dto.OpPut,
					Data: map[string]interface{}{
						"id":     "123",
						"status": "active",
					},
				},
				{
					OperationType: dto.OpDelete,
					Data: map[string]interface{}{
						"id": "1",
					},
				},
			},
			ExpectedHTTPCode: http.StatusOK,
			ExpectedResponse: `{"success":"all"}`,
		},
		{
			Name: "Valid OpCreate Invalid OpPut",
			RequestBody: []dto.Payload{
				{
					OperationType: dto.OpCreate,
					Data: map[string]interface{}{
						"title":       "Task Title",
						"description": "Task Description",
					},
				},
				{
					OperationType: dto.OpPut,
					Data: map[string]interface{}{
						"id": "123",
					},
				},
			},
			ExpectedHTTPCode: http.StatusBadRequest,
			ExpectedResponse: `{"errors":[{"code":"task:map[id:123]-has-errors:-code=Status-is-required,-message=Status-is-required,-err=Key:-'TaskUpdateRequest.Status'-Error:Field-validation-for-'Status'-failed-on-the-'required'-tag","message":"task:map[id:123] has errors: code=Status-is-required, message=Status is required, err=Key: 'TaskUpdateRequest.Status' Error:Field validation for 'Status' failed on the 'required' tag"}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			payloadCollection := &dto.PayloadCollection{Payloads: tc.RequestBody}
			payloadCollectionBytes, err := json.Marshal(payloadCollection)
			if err != nil {
				log.Fatal("cannot marshal json")
			}

			req, err := http.NewRequest("POST", "/multiple/", bytes.NewBuffer(payloadCollectionBytes))
			assert.NoError(t, err)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			errorHandler(th.processMultipleTasks)(w, req)

			assert.Equal(t, tc.ExpectedHTTPCode, w.Code)

			assert.JSONEq(t, tc.ExpectedResponse, w.Body.String())
		})
	}
}
