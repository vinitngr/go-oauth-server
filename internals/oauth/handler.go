package oauth

import (
	"net/http"
	"strings"
)

type Handler struct {
	reg *OauthRegistry
}

func NewHandler(reg *OauthRegistry) *Handler {
	return &Handler{reg: reg}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/oauth/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 {
		http.Error(w, "invalid oauth path", http.StatusBadRequest)
		return
	}

	action := parts[0]
	provider := parts[1]

	processor, ok := h.reg.Get(provider)
	if !ok {
		http.Error(w, "unknown provider", http.StatusNotFound)
		return
	}

	switch action {
	case "connection":
		processor.Connect(w, r)
	case "callback":
		processor.Callback(w, r)
	default:
		http.Error(w, "invalid oauth action", http.StatusBadRequest)
	}
}
