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

docker-build:
	docker build -t todo-app-app .

docker-run:
	docker run --rm todo-app-app

deploy: docker-build
	docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/todoapp
	k3d image import todo-app-app
	kubectl apply -f manifests
