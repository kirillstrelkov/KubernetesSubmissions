all: deploy k3d-apply-decrypt

.PHONY: docker-build deploy

docker-build:
	docker build -t todo-app-app ./todo-app
	docker build -t $$DOCKER_FULL_REPO/todo-app-app:v1 ./todo-app
	docker push $$DOCKER_FULL_REPO/todo-app-app:v1

	docker build -t todo-backend-app ./todo-backend
	docker build -t $$DOCKER_FULL_REPO/todo-backend-app:v1 ./todo-backend
	docker push $$DOCKER_FULL_REPO/todo-backend-app:v1
	
	docker build -t create-post-run ./scripts/post
	docker build -t $$DOCKER_FULL_REPO/create-post-run:v1 ./scripts/post
	docker push $$DOCKER_FULL_REPO/create-post-run:v1

deploy: docker-build
	docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/todoapp
	k3d image import todo-app-app
	k3d image import todo-backend-app
	k3d image import create-post-run
	kubectl apply -f manifests

k3d-apply-decrypt:
	export SOPS_AGE_KEY_FILE=../key.txt; \
	find ./manifests/enc -name '*.yaml' -exec sh -c 'sops --decrypt "{}" | kubectl apply -f -' \;

encrypt:
	@PUBKEY=$$(grep '# public key:' ../key.txt | cut -d ':' -f 2 | tr -d ' ') ; \
	sops --encrypt --age $$PUBKEY --encrypted-regex '(stringData)' ./manifests/dec/secrets.yaml > ./manifests/enc/secrets.yaml

decrypt:
	mkdir -p ./manifests/dec/
	export SOPS_AGE_KEY_FILE=../key.txt; \
	sops --decrypt ./manifests/enc/secrets.yaml > ./manifests/dec/secrets.yaml

clean:
	kubectl delete namespace project || true

gke: decrypt
	kubectl apply -k .
