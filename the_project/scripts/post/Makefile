
docker-build:
	docker build -t create-post-run .
	k3d image import create-post-run

docker-run:
	docker run --network host --rm -e URL=http://localhost:8081/posts create-post-run
