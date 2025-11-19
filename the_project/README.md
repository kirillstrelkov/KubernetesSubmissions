# Todo app

## Deploy to k3d cluster

1. Build, import and deploy `make`

NOTE: create cluster with `k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2`
NOTE 2: make folder in agent0 `docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/todoapp`

## postgresql

```bash
DB_URL=postgresql://<username>:<password>@172.17.0.2:5432/pingpongdb?sslmode=disable
```

Replace `username` and `password` with proper values

## Work with encrypted yaml

If new `key.txt` is created, `manifests/deployment.yaml` should be created from `manifests/enc/deployment.yaml` but with proper `value` for `DB_URL`

```bash
sops --encrypt \
	--age $(grep '# public key:' ../key.txt | cut -d ':' -f 2 | tr -d ' ') \
	--encrypted-regex '(Data)$' \
	./manifests/secrets.yaml > ./manifests/enc/secrets.yaml
```
