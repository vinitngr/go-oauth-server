package oauth

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/vinitngr/go-oauth-server/internals/auth"
	"github.com/vinitngr/go-oauth-server/internals/config"
	"github.com/vinitngr/go-oauth-server/internals/providers/google"
)

type GoogleOAuth struct {
	cfg config.Config
}

func NewGoogle(cfg config.Config) *GoogleOAuth {
	return &GoogleOAuth{cfg: cfg}
}

func (g *GoogleOAuth) Connect(w http.ResponseWriter, r *http.Request) {
	u := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&access_type=offline&prompt=consent",
		g.cfg.GoogleClientID,
		url.QueryEscape(g.cfg.GoogleRedirectURI),
		url.QueryEscape("openid email profile"),
	)

	http.Redirect(w, r, u, http.StatusFound)
}

func (g *GoogleOAuth) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	accessToken, err := google.ExchangeCode(code, g.cfg)
	if err != nil {
		http.Error(w, "token exchange failed", http.StatusInternalServerError)
		return
	}

	user, err := google.FetchGoogleUser(accessToken)
	if err != nil {
		http.Error(w, "user fetch failed", http.StatusInternalServerError)
		return
	}

	claims := map[string]any{
		"google_id": user.ID,
		"email":     user.Email,
		"name":      user.Name,
		"avatar":    user.Picture,
	}

	jwt, _ := auth.SignJWT(claims, g.cfg.JWTSecret, 24*time.Hour)
	http.SetCookie(w, auth.SessionCookie(jwt))

	http.Redirect(w, r, "/", http.StatusFound)
}
