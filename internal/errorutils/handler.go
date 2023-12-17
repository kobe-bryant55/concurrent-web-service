package errorutils

import (
	"errors"
	"net/http"
)

func Handler(err error) (int, *APIErrors) {
	var ae *APIError
	aes := new(APIErrors)

	aes.Errors = make([]*APIError, 0)

	ok := errors.As(err, &ae)
	if ok {
		aes.Errors = append(aes.Errors, ae)
		return StatusCode(ae.Code), aes
	}

	ok = errors.As(err, &aes)
	if ok {
		return http.StatusBadRequest, aes
	}

	code := Code(err)

	ne := &APIError{
		Code:    code,
		Message: err.Error(),
		Err:     err,
	}

	aes.Errors = append(aes.Errors, ne)

	return StatusCode(code), aes
}
