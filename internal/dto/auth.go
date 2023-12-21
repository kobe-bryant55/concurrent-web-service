package dto

import "github.com/MehmetTalhaSeker/concurrent-web-service/internal/types"

// Claims is the JWT claims.
type Claims struct {
	Role types.Role `json:"role"`
}

// TokenResponse is the response body for the user login endpoint.
type TokenResponse struct {
	Token string `json:"token"`
}
