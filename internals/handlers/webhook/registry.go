package webhook

import (
	"net/http"

	"github.com/vinitngr/go-oauth-server/internals/providers/github"
)

type WebhookProcessor interface {
		Process(w http.ResponseWriter, r *http.Request)
}
type WebhookRegistry struct {
	handlers map[string]WebhookProcessor
}

func NewWebhookRegistry() *WebhookRegistry {
	return &WebhookRegistry{
		handlers: map[string]WebhookProcessor{
			"github": &github.GitHubWebhook{},
		},
	}
}

func (r *WebhookRegistry) Get(provider string) (WebhookProcessor, bool) {
	h, ok := r.handlers[provider]
	return h, ok
}