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
	docker build -t $$DOCKER_FULL_REPO/log-output-app:v1 .
	docker push $$DOCKER_FULL_REPO/log-output-app:v1

docker-run:
	docker run -d --rm log-output-app

k3d-apply:
	kubectl apply -k .

gke:
	echo "NOTE: ping_pong_app should be deployed"

# 	remove old resources if exist
	kubectl delete -f manifests/ingress.yaml || true
	kubectl delete -f ../ping_pong_app/manifests/ingress.yaml || true
	kubectl delete -n exercises ingress log-output-pingpong-ingress || true
	
	kubectl apply -f manifests/namespace.yaml
	kubectl apply -f manifests/configmap.yaml
	kubectl apply -f manifests/service.yaml
	kubectl apply -f manifests/gke/

clean:
	kubectl delete namespace exercises || true
