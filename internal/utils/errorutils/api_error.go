package errorutils

import (
	"fmt"
	"strings"
)

// APIError returns human-readable and machine-readable error response.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"` // for developers to see in logs
}

// Error makes it compatible with `error` interface.
func (ae *APIError) Error() string {
	return fmt.Sprintf("code=%s, message=%v, err=%v", ae.Code, ae.Message, ae.Err)
}

// APIErrors returns multiple APIError.
type APIErrors struct {
	Errors []*APIError `json:"errors"`
}

// Error makes it compatible with `error` interface.
func (aes *APIErrors) Error() string {
	list := make([]string, 0)
	for _, ae := range aes.Errors {
		list = append(list, ae.Error())
	}

	return strings.Join(list, "|")
}

// New creates a new `error` compatible APIError.
func New(reason error, err error) *APIError {
	return &APIError{
		Code:    Code(reason),
		Message: reason.Error(),
		Err:     err,
	}
}
