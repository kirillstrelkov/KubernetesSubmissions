# Ping pong app

## Deploy to k3d cluster

1. Build, import, deploy `make`

## postgresql

```bash
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
