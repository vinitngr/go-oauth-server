BIN := server
MAIN := ./cmd/main.go
TUNNEL_CFG := $(HOME)/.cloudflared/8080.yml
TUNNEL_NAME := 8080

.PHONY: build tunnel run

build:
	go build -o $(BIN) $(MAIN)

tunnel:
	cloudflared tunnel --config $(TUNNEL_CFG) run $(TUNNEL_NAME)

run: build tunnel