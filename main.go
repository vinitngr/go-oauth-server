package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	clientID     string
	clientSecret string
	redirectURI  string

	jwtSecret []byte
)

type User struct {
	GithubID int    `json:"github_id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	clientID = os.Getenv("GITHUB_CLIENT_ID")
	clientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	redirectURI = os.Getenv("GITHUB_REDIRECT_URI")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		log.Fatal("missing env vars")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/oauth/github/login", githubLogin)
	mux.HandleFunc("/oauth/github/callback", githubCallback)
	mux.HandleFunc("/logout", logout)
	mux.Handle("/api/user", authMiddleware(http.HandlerFunc(userAPI)))

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func githubLogin(w http.ResponseWriter, r *http.Request) {
	u := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user",
		clientID,
		url.QueryEscape(redirectURI),
	)
	http.Redirect(w, r, u, http.StatusFound)
}

func githubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", 400)
		return
	}

	resp, err := http.Post(
		"https://github.com/login/oauth/access_token",
		"application/x-www-form-urlencoded",
		strings.NewReader(url.Values{
			"client_id":     {clientID},
			"client_secret": {clientSecret},
			"code":          {code},
			"redirect_uri":  {redirectURI},
		}.Encode()),
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	vals, _ := url.ParseQuery(string(body))
	accessToken := vals.Get("access_token")
	if accessToken == "" {
		http.Error(w, "token exchange failed", 500)
		return
	}

	user, err := fetchGitHubUser(accessToken)
	if err != nil {
		http.Error(w, "github user fetch failed", 500)
		return
	}

	claims := jwt.MapClaims{
		"github_id": user.GithubID,
		"login":     user.Login,
		"name":      user.Name,
		"avatar":    user.Avatar,
		"email":     user.Email,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "jwt sign failed", 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    signed,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "unauthorized", 401)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", 401)
			return
		}

		ctx := context.WithValue(r.Context(), "user", token.Claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func userAPI(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(jwt.MapClaims)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claims)
}

func fetchGitHubUser(token string) (User, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	var u struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
	}

	json.NewDecoder(resp.Body).Decode(&u)

	return User{
		GithubID: u.ID,
		Login:    u.Login,
		Name:     u.Name,
		Avatar:   u.AvatarURL,
		Email:    u.Email,
	}, nil
}

func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})
	w.WriteHeader(http.StatusNoContent)
}
