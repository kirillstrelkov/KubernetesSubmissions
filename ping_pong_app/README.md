# Ping pong app

## Deploy to k3d cluster

0. Build docker images `make docker-build`
1. Import image `k3d image import ping-pong-app`
2. Deploy all `kubectl apply -f manifests`
