all: fmt lint

.PHONY: all fmt lint build run


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
