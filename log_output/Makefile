all: docker-build import

.PHONY: build import docker-build docker-run


import:
	k3d image import log-output-writer-app
	k3d image import log-output-reader-app


docker-build:
	docker build -t log-output-writer-app ./writer
	docker build -t log-output-reader-app ./reader

docker-run:
	docker run -d --rm log-output-writer-app
	docker run -d --rm log-output-reader-app
