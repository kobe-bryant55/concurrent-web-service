package apputils

import (
	"fmt"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/types"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sign-in method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func CreateJWT(r types.Role) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": time.Now().Add(48 * time.Hour),
		"role":      r,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
