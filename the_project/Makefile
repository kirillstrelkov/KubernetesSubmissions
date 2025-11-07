all: deploy

.PHONY: docker-build deploy

docker-build:
	docker build -t todo-app-app ./todo-app
	docker build -t todo-backend-app ./todo-backend

deploy: docker-build
	docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/todoapp
	k3d image import todo-app-app
	k3d image import todo-backend-app
	kubectl apply -f manifests
