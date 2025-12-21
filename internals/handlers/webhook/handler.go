package webhook

import (
	"net/http"
	"strings"
)

type WebhookHandler struct {
	reg *WebhookRegistry
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	provider := strings.TrimPrefix(r.URL.Path, "/webhook/")
	if provider == "" {
		http.Error(w, "provider missing", http.StatusBadRequest)
		return
	}

	processor, ok := h.reg.Get(provider)
	if !ok {
		http.Error(w, "unknown provider", http.StatusNotFound)
		return
	}

	processor.Process(w, r)
}

func NewWebhookHandler(reg *WebhookRegistry) *WebhookHandler {
	return &WebhookHandler{reg: reg}
}