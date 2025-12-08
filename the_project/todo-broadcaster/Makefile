all: fmt lint

.PHONY: all fmt lint build run


fmt:
	go fmt ./...
	go run golang.org/x/tools/cmd/goimports@latest -w .

lint:
	go vet ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run
	golangci-lint run

build:
	go build

run:
	go get
	@echo "export TELEGRAM_BOT_TOKEN=<add bot token>"
	@echo "export TELEGRAM_CHAT_ID=<add chat id>"
	NATS_URL=http://localhost:42222/ PORT=8087 go run main.go

docker-build:
	docker build -t todo-broadcaster-app .

docker-run:
	docker run --rm -e PORT=8087 todo-broadcaster-app

deploy: docker-build
	k3d image import todo-broadcaster-app
