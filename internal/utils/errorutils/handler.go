package errorutils

import (
	"errors"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apiutils"
	"net/http"
)

func Handler(err error, w http.ResponseWriter) {
	var ae *APIError

	ok := errors.As(err, &ae)
	if ok {
		apiutils.WriteJSON(w, StatusCode(ae.Code), ae)
		return
	}

	var aes *APIErrors

	ok = errors.As(err, &aes)
	if ok {
		apiutils.WriteJSON(w, http.StatusBadRequest, aes)
		return
	}

	code := Code(err)

	apiutils.WriteJSON(w, StatusCode(code), &APIError{
		Code:    code,
		Message: err.Error(),
		Err:     err,
	})
}
