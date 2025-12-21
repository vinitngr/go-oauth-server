package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GithubClientID     string
	GithubClientSecret string
	GithubRedirectURI  string
	JWTSecret          []byte
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	cfg := Config{
		GithubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		GithubRedirectURI:  os.Getenv("GITHUB_REDIRECT_URI"),
		JWTSecret:          []byte(os.Getenv("JWT_SECRET")),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURI:  os.Getenv("GOOGLE_REDIRECT_URI"),
	}

	must(cfg.GithubClientID, "GITHUB_CLIENT_ID")
	must(cfg.GithubClientSecret, "GITHUB_CLIENT_SECRET")
	must(cfg.GithubRedirectURI, "GITHUB_REDIRECT_URI")
	must(string(cfg.JWTSecret), "JWT_SECRET")

	return cfg
}

func must(v, name string) {
	if v == "" {
		log.Fatalf("missing env var: %s", name)
	}
}
