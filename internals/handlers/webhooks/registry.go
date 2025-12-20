package webhooks

import "net/http"

type WebhookProcessor interface {
		Process(w http.ResponseWriter, r *http.Request)
}
type WebhookRegistry struct {
	handlers map[string]WebhookProcessor
}

func NewWebhookRegistry() *WebhookRegistry {
	return &WebhookRegistry{
		handlers: map[string]WebhookProcessor{
			"github": &GitHubWebhook{},
		},
	}
}

func (r *WebhookRegistry) Get(provider string) (WebhookProcessor, bool) {
	h, ok := r.handlers[provider]
	return h, ok
}