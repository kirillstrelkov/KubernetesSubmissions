all: docker-build import k3d-apply k3d-apply-decrypt

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
	@echo "Please run manually next command. Check README.md for more information."
	@echo 'DB_URL="<postgres_url>" go run main.go'

docker-build:
	docker build -t ping-pong-app .

docker-run:
	docker run --rm ping-pong-app

import:
	k3d image import ping-pong-app

k3d-apply:
	kubectl apply -f manifests

k3d-apply-decrypt:
	export SOPS_AGE_KEY_FILE=../key.txt; \
	find ./manifests/enc -name '*.yaml' -exec sh -c 'sops --decrypt "{}" | kubectl apply -f -' \;
