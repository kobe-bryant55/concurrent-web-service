package main

import (
	"fmt"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/dto"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/types"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apiutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
	"net/http"
	"regexp"
)

var (
	generateAdminJWT = regexp.MustCompile(`^/auth/admin/*$`)
	generateUserJWT  = regexp.MustCompile(`^/auth/user/*$`)
)

type authHandler struct {
}

func newAuthHandler() *authHandler {
	return &authHandler{}
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && generateAdminJWT.MatchString(r.URL.Path):
		errorHandler(h.generateAdminToken)(w, r)
		return

	case r.Method == http.MethodGet && generateUserJWT.MatchString(r.URL.Path):
		errorHandler(h.generateUserToken)(w, r)
		return

	default:
		apiutils.WriteJSON(w, http.StatusMethodNotAllowed, errorutils.New(errorutils.ErrMethodNotAllowed, nil))
		return
	}
}

func (h *authHandler) generateAdminToken(w http.ResponseWriter, r *http.Request) error {
	token, err := apputils.CreateJWT(types.Admin)
	if err != nil {
		return errorutils.New(fmt.Errorf("something went wrong"), nil)
	}

	apiutils.WriteJSON(w, http.StatusOK, dto.TokenResponse{Token: token})
	return nil
}

func (h *authHandler) generateUserToken(w http.ResponseWriter, r *http.Request) error {
	token, err := apputils.CreateJWT(types.Registered)
	if err != nil {
		return errorutils.New(fmt.Errorf("something went wrong"), nil)
	}

	apiutils.WriteJSON(w, http.StatusOK, dto.TokenResponse{Token: token})
	return nil
}
