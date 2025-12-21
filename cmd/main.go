package main

import (
	"log"
	"net/http"

	"github.com/vinitngr/go-oauth-server/internals/auth"
	"github.com/vinitngr/go-oauth-server/internals/config"
	"github.com/vinitngr/go-oauth-server/internals/handlers"
	"github.com/vinitngr/go-oauth-server/internals/handlers/webhooks"
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

	webhookReg := webhooks.NewWebhookRegistry()
	mux.Handle("/webhook/", webhooks.NewWebhookHandler(webhookReg))

	oauthReg := oauth.NewOauthRegistry(cfg)
	handler := oauth.NewHandler(oauthReg)
	mux.Handle("/oauth/connection/", handler)
	mux.Handle("/oauth/callback/", handler)

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
