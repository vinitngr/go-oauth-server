package auth

import "net/http"

func SessionCookie(jwt string) *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Value:    jwt,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	}
}
