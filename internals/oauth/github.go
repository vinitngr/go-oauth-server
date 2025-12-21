package oauth

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/vinitngr/go-oauth-server/internals/auth"
	"github.com/vinitngr/go-oauth-server/internals/config"
	"github.com/vinitngr/go-oauth-server/internals/providers/github"
)

type GitHubOAuth struct {
	cfg config.Config
}

func NewGithub(cfg config.Config) *GitHubOAuth {
	return &GitHubOAuth{cfg: cfg}
}

func (g *GitHubOAuth) Connect(w http.ResponseWriter, r *http.Request) {
	u := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user",
		g.cfg.GithubClientID,
		url.QueryEscape(g.cfg.GithubRedirectURI),
	)
	http.Redirect(w, r, u, http.StatusFound)
}

func (g *GitHubOAuth) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	accessToken, err := github.ExchangeCode(code, g.cfg)
	if err != nil {
		http.Error(w, "token exchange failed", http.StatusInternalServerError)
		return
	}

	user, err := github.FetchGitHubUser(accessToken)
	if err != nil {
		http.Error(w, "user fetch failed", http.StatusInternalServerError)
		return
	}

	claims := map[string]any{
		"github_id": user.GithubID,
		"login":     user.Login,
		"name":      user.Name,
		"avatar":    user.Avatar,
		"email":     user.Email,
	}

	jwt, _ := auth.SignJWT(claims, g.cfg.JWTSecret, 24*time.Hour)
	http.SetCookie(w, auth.SessionCookie(jwt))

	http.Redirect(w, r, "/", http.StatusFound)
}
