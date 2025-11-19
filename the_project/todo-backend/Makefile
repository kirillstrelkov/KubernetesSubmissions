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
	go run main.go

docker-build:
	docker build -t todo-backend-app .

docker-run:
	docker run --rm -e PORT=8085 -p 8085:8085 todo-backend-app

deploy: docker-build
	k3d image import todo-backend-app
