package errorutils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// API Errors.
var (
	ErrBadRequest       = errors.New("bad request")
	ErrMethodNotAllowed = errors.New("method not allowed")
)

// Common errors.
var (
	ErrJSONDecode = errors.New("json decode error")
	ErrJSONEncode = errors.New("json encode error")
	ErrInvalidID  = errors.New("invalid ID")
)

// Task Errors.
var (
	ErrTaskCount    = errors.New("task count failed")
	ErrTaskCreate   = errors.New("task create failed")
	ErrTaskDelete   = errors.New("task delete failed")
	ErrTaskRead     = errors.New("task read failed")
	ErrTaskReads    = errors.New("task reads failed")
	ErrTaskUpdate   = errors.New("task update failed")
	ErrTaskNotFound = errors.New("task not found")
)

// errorCodes a map to store error codes.
var errorCodes = map[error]string{
	// API
	ErrBadRequest:       ErrCodeBadRequest,
	ErrMethodNotAllowed: ErrCodeMethodNotAllowed,

	// Common
	ErrJSONDecode: ErrCodeJSONDecode,
	ErrJSONEncode: ErrCodeJSONEncode,
	ErrInvalidID:  ErrCodeInvalidID,

	// Tasks
	ErrTaskCount:    ErrCodeTaskCount,
	ErrTaskCreate:   ErrCodeTaskCreate,
	ErrTaskDelete:   ErrCodeTaskDelete,
	ErrTaskRead:     ErrCodeTaskRead,
	ErrTaskReads:    ErrCodeTaskReads,
	ErrTaskUpdate:   ErrCodeTaskUpdate,
	ErrTaskNotFound: ErrCodeTaskNotFound,
}

// Code gets machine-readable error code from error.
func Code(err error) string {
	if code, ok := errorCodes[err]; ok {
		return code
	}
	return strings.ReplaceAll(err.Error(), " ", "-")
}

var statusCodeMap = map[string]int{
	// API
	ErrCodeBadRequest:       http.StatusBadRequest,
	ErrCodeMethodNotAllowed: http.StatusMethodNotAllowed,

	// Common
	ErrCodeJSONDecode: http.StatusUnprocessableEntity,
	ErrCodeJSONEncode: http.StatusUnprocessableEntity,
	ErrCodeInvalidID:  http.StatusBadRequest,

	// Task
	ErrCodeTaskCount:    http.StatusUnprocessableEntity,
	ErrCodeTaskCreate:   http.StatusUnprocessableEntity,
	ErrCodeTaskDelete:   http.StatusUnprocessableEntity,
	ErrCodeTaskRead:     http.StatusUnprocessableEntity,
	ErrCodeTaskReads:    http.StatusUnprocessableEntity,
	ErrCodeTaskUpdate:   http.StatusUnprocessableEntity,
	ErrCodeTaskNotFound: http.StatusNotFound,
}

// StatusCode gets HTTP status code from error code.
func StatusCode(code string) int {
	if status, ok := statusCodeMap[code]; ok {
		return status
	}

	return http.StatusInternalServerError
}

func ValidationError(errors []*APIError) error {
	return &APIErrors{
		errors,
	}
}

func Required(field string) error {
	return fmt.Errorf("%s is required", field)
}

func Max(field string) error {
	return fmt.Errorf("%s too long", field)
}

func Min(field string) error {
	return fmt.Errorf("%s too short", field)
}

func Len(field string) error {
	return fmt.Errorf("%s invalid", field)
}
