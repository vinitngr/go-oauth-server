package oauth

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/vinitngr/go-oauth-server/internals/auth"
	"github.com/vinitngr/go-oauth-server/internals/config"
)

func GithubLogin(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := fmt.Sprintf(
			"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user",
			cfg.GithubClientID,
			url.QueryEscape(cfg.GithubRedirectURI),
		)
		http.Redirect(w, r, u, http.StatusFound)
	}
}

func GithubCallback(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "missing code", 400)
			return
		}

		accessToken, err := exchangeCode(code, cfg)
		if err != nil {
			http.Error(w, "token exchange failed", 500)
			return
		}

		user, err := fetchGitHubUser(accessToken)
		if err != nil {
			http.Error(w, "user fetch failed", 500)
			return
		}

		claims := map[string]any{
			"github_id": user.GithubID,
			"login":     user.Login,
			"name":      user.Name,
			"avatar":    user.Avatar,
			"email":     user.Email,
		}

		jwt, _ := auth.SignJWT(claims, cfg.JWTSecret, 24*time.Hour)
		http.SetCookie(w, auth.SessionCookie(jwt))

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
