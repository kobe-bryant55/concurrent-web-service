package errorutils

import (
	"fmt"
	"net/http"
	"strings"
)

// errorCodes a map to store error codes.
var errorCodes = map[error]string{
	// Auth
	ErrInvalidToken:      ErrCodeInvalidToken,
	ErrExpiredToken:      ErrCodeExpiredToken,
	ErrMissingAuthHeader: ErrCodeMissingAuthHeader,

	// API
	ErrBadRequest:       ErrCodeBadRequest,
	ErrMethodNotAllowed: ErrCodeMethodNotAllowed,

	// Common
	ErrJSONDecode:   ErrCodeJSONDecode,
	ErrJSONEncode:   ErrCodeJSONEncode,
	ErrInvalidID:    ErrCodeInvalidID,
	ErrUnauthorized: ErrCodeUnauthorized,

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
	// Auth
	ErrCodeInvalidToken:      http.StatusUnauthorized,
	ErrCodeExpiredToken:      http.StatusUnauthorized,
	ErrCodeMissingAuthHeader: http.StatusUnauthorized,

	// API
	ErrCodeBadRequest:       http.StatusBadRequest,
	ErrCodeMethodNotAllowed: http.StatusMethodNotAllowed,

	// Common
	ErrCodeJSONDecode:   http.StatusUnprocessableEntity,
	ErrCodeJSONEncode:   http.StatusUnprocessableEntity,
	ErrCodeInvalidID:    http.StatusBadRequest,
	ErrCodeUnauthorized: http.StatusUnauthorized,

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
