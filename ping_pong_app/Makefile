all: docker-build import k3d-apply

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

docker-build:
	docker build -t ping-pong-app .

docker-run:
	docker run --rm ping-pong-app

import:
	k3d image import ping-pong-app

k3d-apply:
	kubectl apply -f manifests
