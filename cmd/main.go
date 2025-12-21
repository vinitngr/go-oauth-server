package main

import (
	"log"
	"net/http"

	"github.com/vinitngr/go-oauth-server/internals/auth"
	"github.com/vinitngr/go-oauth-server/internals/config"
	"github.com/vinitngr/go-oauth-server/internals/handlers"
	"github.com/vinitngr/go-oauth-server/internals/handlers/webhook"
	"github.com/vinitngr/go-oauth-server/internals/oauth"
	// oauth "github.com/vinitngr/go-oauth-server/internals/oauth/github"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()

	mux.Handle("/api/user",
		auth.Middleware(cfg.JWTSecret)(
			http.HandlerFunc(handlers.User),
		),
	)

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	mux.Handle("/logout",
		auth.Middleware(cfg.JWTSecret)(
			http.HandlerFunc(handlers.Logout),
		),
	)

	webhookReg := webhook.NewWebhookRegistry()
	mux.Handle("/webhook/", webhook.NewWebhookHandler(webhookReg))

	oauthReg := oauth.NewOauthRegistry(cfg)
	handler := oauth.NewHandler(oauthReg)
	mux.Handle("/oauth/connection/", handler)
	mux.Handle("/oauth/callback/", handler)

	log.Println("https://8080.vinitngr.xyz")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
