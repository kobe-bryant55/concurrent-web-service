package main

import (
	"bytes"
	"encoding/json"
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

	th := newTaskHandler(ms, validator, jobQueue)

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
			ExpectedHTTPCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			payloadCollection := &dto.PayloadCollection{Payloads: tc.RequestBody, Token: "token"}
			payloadCollectionBytes, err := json.Marshal(payloadCollection)
			if err != nil {
				log.Fatal("cannot marshal json")
			}

			req, err := http.NewRequest("POST", "/multiple/", bytes.NewBuffer(payloadCollectionBytes))
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			err = th.processMultipleTasks(w, req)

			assert.Equal(t, tc.ExpectedHTTPCode, w.Code)

			if tc.ExpectedResponse != "" {
				assert.JSONEq(t, tc.ExpectedResponse, w.Body.String())
			}
		})
	}
}
