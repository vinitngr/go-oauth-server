BIN := server
MAIN := ./cmd/main.go
TUNNEL_CFG := $(HOME)/.cloudflared/8080.yml
TUNNEL_NAME := 8080

.PHONY: build run tunnel dev

build:
	go build -o $(BIN) $(MAIN)

run:
	./$(BIN)

tunnel:
	cloudflared tunnel --config $(TUNNEL_CFG) run $(TUNNEL_NAME)

dev: build
	$(MAKE) -j2 run tunnel
