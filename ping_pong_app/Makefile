all: ns docker-build import k3d-apply

.PHONY: all fmt lint build run

ns:
	kubectl create namespace exercises || true
	kubens exercises

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

# Use env variable from gke/README.md
docker-build-gke:
	docker build -t $$DOCKER_FULL_REPO/ping-pong-app:v3 .
	docker push $$DOCKER_FULL_REPO/ping-pong-app:v3

docker-build:
	docker build -t ping-pong-app:v3 .
	docker build -t kirillstrelkov/ping-pong-app:v3 .
	docker push kirillstrelkov/ping-pong-app:v3

docker-run:
	docker run --rm ping-pong-app

import:
	k3d image import ping-pong-app:v3

k3d-apply: ns decrypt
	kubectl apply -k .

argocd-deploy: ns
	kubectl kustomize kust-argocd --load-restrictor LoadRestrictionsNone | kubectl apply -f -

argocd-add-app: ns
	kubectl apply -f manifests/argocd-app.yaml

gke: docker-build-gke
	kubectl apply -f manifests/namespace.yaml
	kubectl apply -f manifests/service.yaml
	export SOPS_AGE_KEY_FILE=../key.txt ; \
	sops --decrypt manifests/enc/statefulset.yaml | kubectl apply -f - ;\
	find ./manifests/gke/enc -name '*.yaml' -exec sh -c 'sops --decrypt "{}" | kubectl apply -f -' \;

decrypt:
	export SOPS_AGE_KEY_FILE=../key.txt ; \
	mkdir -p ./manifests/dec/ ; \
	sops --decrypt ./manifests/enc/deployment.yaml > ./manifests/dec/deployment.yaml ; \
	sops --decrypt ./manifests/enc/statefulset.yaml > ./manifests/dec/statefulset.yaml ; \
	sops --decrypt ./manifests/gke/enc/deployment.yaml > ./manifests/dec/deployment_gke.dec.yaml

encrypt:
	@PUBKEY=$$(grep '# public key:' ../key.txt | cut -d ':' -f 2 | tr -d ' ') ; \
	sops --encrypt --age $$PUBKEY --encrypted-regex '(value)' ./manifests/dec/deployment.yaml > ./manifests/enc/deployment.yaml ; \
	sops --encrypt --age $$PUBKEY --encrypted-regex '(value)' ./manifests/dec/statefulset.yaml > ./manifests/enc/statefulset.yaml ; \
	sops --encrypt --age $$PUBKEY --encrypted-regex '(value)' ./manifests/dec/deployment_gke.dec.yaml > ./manifests/gke/enc/deployment.yaml

clean:
	kubectl delete namespace exercises || true

