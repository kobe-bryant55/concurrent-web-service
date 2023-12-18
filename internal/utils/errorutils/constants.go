package errorutils

// API Error Codes.
const (
	ErrCodeBadRequest       = "req/bad-request"
	ErrCodeMethodNotAllowed = "req/method-not-allowed"
)

// Common Error Codes.
const (
	ErrCodeJSONDecode = "com/json-decode"
	ErrCodeJSONEncode = "com/json-encode"
	ErrCodeInvalidID  = "com/invalid-id"
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
