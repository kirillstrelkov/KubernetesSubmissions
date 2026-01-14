all: docker-build import k3d-apply

SHELL := /bin/bash
.ONESHELL:
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
	docker build -t log-output-app:v4 .

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

install-argo-rollouts:
	pushd ../exercises/probes
	make setup || true
	popd

deploy-pingpong-to-k3s-istio:
	pushd ../ping_pong_app
	make docker-build
	k3d image import -c k3s-iostio ping-pong-app:v1
	k3d image import -c k3s-iostio ping-pong-app:v2
	make k3d-apply || true
	popd

deploy-logoutput-to-k3s-istio:
	make docker-build
	k3d image import -c k3s-iostio log-output-app:v4
	make k3d-apply || true


deploy-to-k3s-istio: install-argo-rollouts deploy-pingpong-to-k3s-istio deploy-logoutput-to-k3s-istio
	@echo "Done"
