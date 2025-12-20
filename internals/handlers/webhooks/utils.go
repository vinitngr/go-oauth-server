package webhooks

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func ParseJSON(body []byte) (map[string]any, error) {
	var data map[string]any
	err := json.Unmarshal(body, &data)
	return data, err
}

func logWebhookLine(provider, event string) {
	f, err := os.OpenFile(
		"webhook.event.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return
	}
	defer f.Close()

	fmt.Fprintf(f, "%s %s %s\n",
		time.Now().Format(time.RFC3339),
		provider,
		event,
	)
}
