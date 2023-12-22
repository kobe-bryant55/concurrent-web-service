package main

import (
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/logger"
	"net/http"
)

type logHandler struct {
	lg *logger.Logger
}

func newLogHandler(lg *logger.Logger) *logHandler {
	return &logHandler{lg: lg}
}

func (lh *logHandler) log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lh.lg.InfoLog.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
