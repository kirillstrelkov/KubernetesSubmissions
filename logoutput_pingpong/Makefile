all: docker-build import kube-apply

docker-build:
	docker build -t log-output-writer-app ../log_output/writer
	docker build -t log-output-reader-app ../log_output/reader
	docker build -t ping-pong-app ../ping_pong_app

import:
	k3d image import log-output-writer-app
	k3d image import log-output-reader-app
	k3d image import ping-pong-app

kube-apply:
	docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/kube/pingpong
	kubectl apply -f manifests
	kubectl apply -f ../log_output/manifests
	kubectl apply -f ../ping_pong_app/manifests
