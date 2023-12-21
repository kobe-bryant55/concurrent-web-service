package apiutils

import (
	"encoding/json"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/logger"
	"net/http"
	"time"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	start := time.Now()

	jb, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			logger.ErrorLog.Println(err)
		}
		return
	}

	w.WriteHeader(status)
	_, err = w.Write(jb)
	if err != nil {
		logger.ErrorLog.Println(err)
		return
	}

	logger.InfoLog.Printf("STATUS:%d TIME:%v DURATION:%v", status, time.Now(), time.Since(start))
}
