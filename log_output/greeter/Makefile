all: fmt lint build

.PHONY: all fmt lint build run import docker-build

fmt:
	go fmt ./...
	go run golang.org/x/tools/cmd/goimports@latest -w .

lint:
	go vet ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

build:
	go build

run:
	go run main.go

docker-build:
	docker build -t greeter:v1 .
