package main

import (
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/appcontext"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/dto"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, errorutils.ErrMissingAuthHeader.Error(), http.StatusUnauthorized)
			return
		}

		ts := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := apputils.ValidateJWT(ts)
		if err != nil {
			http.Error(w, errorutils.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		mc, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, errorutils.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		dts, ok := mc["expiresAt"].(string)
		if !ok {
			http.Error(w, errorutils.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		exp, err := time.Parse(time.RFC3339, dts)
		if err != nil {
			http.Error(w, errorutils.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		if exp.Unix() < time.Now().Unix() {
			http.Error(w, errorutils.ErrExpiredToken.Error(), http.StatusUnauthorized)
			return
		}

		claims := &dto.Claims{}

		err = apputils.InterfaceToStruct(mc, claims)
		if err != nil {
			http.Error(w, errorutils.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		ctx := appcontext.WithRole(r.Context(), claims.Role)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
