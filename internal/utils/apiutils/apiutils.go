package apiutils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")

	jb, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(status)
	w.Write(jb)
}
