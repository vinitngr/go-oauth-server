package auth

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func Middleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("session")
			if err != nil {
				http.Error(w, "unauthorized", 401)
				return
			}

			t, err := jwt.Parse(c.Value, func(*jwt.Token) (any, error) {
				return secret, nil
			})
			if err != nil || !t.Valid {
				http.Error(w, "invalid token", 401)
				return
			}

			ctx := context.WithValue(r.Context(), "user", t.Claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
