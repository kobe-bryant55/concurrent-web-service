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

	case r.Method == http.MethodGet && generateUserJWT.MatchString(r.URL.Path):
		errorHandler(h.generateUserToken)(w, r)

	default:
		apiutils.WriteJSON(w, http.StatusMethodNotAllowed, errorutils.New(errorutils.ErrMethodNotAllowed, nil))
	}
}

// generateAdminToken godoc
// @Summary Generate admin token
// @Tags AUTH
// @Description Generate an admin token.
// @Produce  json
// @Success 200 {object} dto.TokenResponse
// @Failure 500  {object}  errorutils.APIErrors
// @Router /auth/admin [get]
func (h *authHandler) generateAdminToken(w http.ResponseWriter, r *http.Request) error {
	token, err := apputils.CreateJWT(types.Admin)
	if err != nil {
		return errorutils.New(fmt.Errorf("something went wrong"), nil)
	}

	apiutils.WriteJSON(w, http.StatusOK, dto.TokenResponse{Token: token})
	return nil
}

// generateUserToken godoc
// @Summary Generate user token
// @Tags AUTH
// @Description Generate an user token.
// @Produce  json
// @Success 200 {object} dto.TokenResponse
// @Failure 500  {object}  errorutils.APIErrors
// @Router /auth/user [get]
func (h *authHandler) generateUserToken(w http.ResponseWriter, r *http.Request) error {
	token, err := apputils.CreateJWT(types.Registered)
	if err != nil {
		return errorutils.New(fmt.Errorf("something went wrong"), nil)
	}

	apiutils.WriteJSON(w, http.StatusOK, dto.TokenResponse{Token: token})
	return nil
}
