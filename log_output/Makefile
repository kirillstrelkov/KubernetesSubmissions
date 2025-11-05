all: docker-build import k3d-apply

.PHONY: all fmt lint build run import docker-build docker-run

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

import:
	k3d image import log-output-app

docker-build:
	docker build -t log-output-app .

docker-run:
	docker run -d --rm log-output-app

k3d-apply:
	kubectl apply -f manifests
