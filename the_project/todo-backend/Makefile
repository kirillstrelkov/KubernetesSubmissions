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
	echo "NATS port should be forwarded from k3d cluster to localhost:42222"
	echo "postgresql should be running in docker container on 172.17.0.2:5432 and database pingpongdb should exist"
	PORT=8085 NATS_URL=http://localhost:42222/ DB_URL=postgresql://username:password@172.17.0.2:5432/pingpongdb?sslmode=disable go run main.go

docker-build:
	docker build -t todo-backend-app .

docker-run:
	docker run --rm -e PORT=8085 -p 8085:8085 todo-backend-app

deploy: docker-build
	k3d image import todo-backend-app
