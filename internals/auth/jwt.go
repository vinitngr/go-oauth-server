package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignJWT(claims map[string]any, secret []byte, ttl time.Duration) (string, error) {
	claims["exp"] = time.Now().Add(ttl).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	return t.SignedString(secret)
}
