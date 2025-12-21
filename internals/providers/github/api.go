package github

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/vinitngr/go-oauth-server/internals/config"
	utils "github.com/vinitngr/go-oauth-server/internals/providers"
)

func FetchGitHubUser(token string) (User, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+ token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	var raw struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
	}

	json.NewDecoder(resp.Body).Decode(&raw)

	return User{
		GithubID: raw.ID,
		Login:    raw.Login,
		Name:     raw.Name,
		Avatar:   raw.AvatarURL,
		Email:    raw.Email,
	}, nil
}

func ExchangeCode(code string, cfg config.Config) (string, error) {
	resp, err := http.Post(
		"https://github.com/login/oauth/access_token",
		"application/x-www-form-urlencoded",
		strings.NewReader(url.Values{
			"client_id":     {cfg.GithubClientID},
			"client_secret": {cfg.GithubClientSecret},
			"code":          {code},
			"redirect_uri":  {cfg.GithubRedirectURI},
		}.Encode()),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	vals, _ := url.ParseQuery(string(body))

	_ = utils.SaveToken("github", vals.Get("access_token"))
	return vals.Get("access_token"), nil
}
