package providers

import (
	"encoding/json"
	"errors"
	"os"
)
func SaveToken(provider, token string) error {
	data := map[string]string{}

	b, _ := os.ReadFile("tokens.json")
	if len(b) > 0 {
		_ = json.Unmarshal(b, &data)
	}

	data[provider] = token

	out, _ := json.MarshalIndent(data, "", "  ")
	return os.WriteFile("tokens.json", out, 0600)
}

func LoadToken(provider string) (string, error) {
	b, err := os.ReadFile("tokens.json")
	if err != nil {
		return "", err
	}

	data := map[string]string{}
	if err := json.Unmarshal(b, &data); err != nil {
		return "", err
	}

	tok, ok := data[provider]
	if !ok {
		return "", errors.New("token not found")
	}

	return tok, nil
}