package webhooks

import (
	"io"
	"net/http"
	"net/url"
)

type GitHubWebhook struct{}

func (g *GitHubWebhook) Process(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	event := r.Header.Get("X-GitHub-Event")
	if event == "" {
		http.Error(w, "missing event", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read error", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		values, err := url.ParseQuery(string(body))
		if err != nil {
			http.Error(w, "invalid form payload", http.StatusBadRequest)
			return
		}
		body = []byte(values.Get("payload"))
	}

	logWebhookLine("github", event)

	w.WriteHeader(http.StatusOK)
}
