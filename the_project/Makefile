all: ns deploy

.PHONY: docker-build deploy

ns:
	kubectl create namespace project || true
	kubens project

install-nats: ns
	helm install -f config/nats.yaml my-nats oci://registry-1.docker.io/bitnamicharts/nats

docker-build-gke:
	docker build -t $$DOCKER_FULL_REPO/todo-app-app:v1 ./todo-app
	docker push $$DOCKER_FULL_REPO/todo-app-app:v1
	docker build -t $$DOCKER_FULL_REPO/todo-backend-app:v1 ./todo-backend
	docker push $$DOCKER_FULL_REPO/todo-backend-app:v1
	docker build -t $$DOCKER_FULL_REPO/create-post-run:v1 ./scripts/post
	docker push $$DOCKER_FULL_REPO/create-post-run:v1
	docker build -t $$DOCKER_FULL_REPO/pg-backup:v2 ./scripts/pg-backup
	docker push $$DOCKER_FULL_REPO/pg-backup:v2

docker-build:
	docker build -t todo-app-app ./todo-app
	docker build -t todo-backend-app ./todo-backend
	docker build -t todo-broadcaster-app ./todo-broadcaster
	docker build -t create-post-run ./scripts/post
	docker build -t pg-backup ./scripts/pg-backup

deploy: docker-build decrypt
	docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/todoapp
	k3d image import todo-app-app
	k3d image import todo-backend-app
	k3d image import todo-broadcaster-app
	k3d image import create-post-run
	k3d image import pg-backup
	kubectl kustomize kust-local --load-restrictor LoadRestrictionsNone | kubectl apply -n project -f -

encrypt:
	@PUBKEY=$$(grep '# public key:' ../key.txt | cut -d ':' -f 2 | tr -d ' ') ; \
	sops --encrypt --age $$PUBKEY --encrypted-regex '(stringData)' ./manifests/dec/secrets.yaml > ./manifests/enc/secrets.yaml

decrypt:
	mkdir -p ./manifests/dec/
	export SOPS_AGE_KEY_FILE=../key.txt; \
	sops --decrypt ./manifests/enc/secrets.yaml > ./manifests/dec/secrets.yaml

clean:
	kubectl delete namespace project || true
	kubectl delete pvc todo-app-pvc || true
	kubectl delete pv todo-app-pv || true
	kubectl delete storageclasses todo-app-pv-class || true

gke: docker-build-gke decrypt
	kubectl apply -k .

create-todos:
	bash ./scripts/create_todos.sh

label-nats-metrics:
	kubectl label servicemonitors.monitoring.coreos.com -n prometheus my-nats-metrics $(kubectl -n prometheus get prometheus -ojson | grep -A 5 "serviceMonitorSelector" | grep kube | sed 's/"//g' | sed 's/: /=/')

argocd: ns
	# kubectl kustomize kust-argocd --load-restrictor LoadRestrictionsNone | kubectl apply -n project -f - || true
	kubectl apply -f manifests/argocd-app.yaml
