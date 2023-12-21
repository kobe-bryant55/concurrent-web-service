package errorutils

import "errors"

// Auth Errors.
var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrExpiredToken      = errors.New("expired token")
	ErrMissingAuthHeader = errors.New("missing authorization header")
)

// API Errors.
var (
	ErrBadRequest       = errors.New("bad request")
	ErrMethodNotAllowed = errors.New("method not allowed")
)

// Common errors.
var (
	ErrJSONDecode   = errors.New("json decode error")
	ErrJSONEncode   = errors.New("json encode error")
	ErrInvalidID    = errors.New("invalid ID")
	ErrUnauthorized = errors.New("unauthorized user")
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

// Auth Error Codes.
const (
	ErrCodeInvalidToken      = "auth/invalid-token"
	ErrCodeExpiredToken      = "auth/expired-token"
	ErrCodeMissingAuthHeader = "auth/missing-header"
)

// API Error Codes.
const (
	ErrCodeBadRequest       = "req/bad-request"
	ErrCodeMethodNotAllowed = "req/method-not-allowed"
)

// Common Error Codes.
const (
	ErrCodeJSONDecode   = "com/json-decode"
	ErrCodeJSONEncode   = "com/json-encode"
	ErrCodeInvalidID    = "com/invalid-id"
	ErrCodeUnauthorized = "un/unauthorized"
)

// Task Error Codes.
const (
	ErrCodeTaskCount    = "task/count-failed"
	ErrCodeTaskCreate   = "task/create-failed"
	ErrCodeTaskDelete   = "task/delete-failed"
	ErrCodeTaskRead     = "task/read-failed"
	ErrCodeTaskReads    = "task/reads-failed"
	ErrCodeTaskUpdate   = "task/update-failed"
	ErrCodeTaskNotFound = "task/not-found"
)
