all: deploy k3d-apply-decrypt

.PHONY: docker-build deploy

docker-build:
	docker build -t todo-app-app ./todo-app
	docker build -t todo-backend-app ./todo-backend

deploy: docker-build
	docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/todoapp
	k3d image import todo-app-app
	k3d image import todo-backend-app
	kubectl apply -f manifests

k3d-apply-decrypt:
	export SOPS_AGE_KEY_FILE=../key.txt; \
	find ./manifests/enc -name '*.yaml' -exec sh -c 'sops --decrypt "{}" | kubectl apply -f -' \;

decrypt:
	export SOPS_AGE_KEY_FILE=../key.txt; \
	sops --decrypt ./manifests/enc/secrets.yaml
