# Ping pong app

## Deploy to k3d cluster

1. Build, import, deploy `make`

## postgresql

```bash
# run postgresql container
docker run --name some-postgres -e POSTGRES_USER=<username> -e POSTGRES_PASSWORD=<password> -e POSTGRES_DB=pingpongdb -d postgres
DB_URL=postgresql://<username>:<password>@172.17.0.2:5432/pingpongdb?sslmode=disable
```

Replace `username` and `password` with proper values

## Work with encrypted yaml

If new `key.txt` is created, files from `manifests/enc/deployment.yaml` should recreated

```bash
sops --encrypt \
	--age $(grep '# public key:' ../key.txt | cut -d ':' -f 2 | tr -d ' ') \
	--encrypted-regex '(value)' \
	./manifests/deployment.yaml > ./manifests/enc/deployment.yaml
```

## Deploy to Kubernetes

Check gke/README.md.

1. Build images and deploy to GKE `make docker-build`
2. Use `make gke` to apply all manifests
