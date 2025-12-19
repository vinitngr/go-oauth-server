package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func User(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(jwt.MapClaims)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claims)
}
