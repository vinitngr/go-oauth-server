build:
    go build -o server ./cmd/main.go

start:
    ./server

tunnel:
    cloudflared tunnel --config ~/.cloudflared/8080.yml run 8080

dev: build
    just start & just tunnel & wait
