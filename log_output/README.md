## Log output app

### Deploy to k3d cluster

1. Build, import and deploy `make`

## Deploy to Kubernetes

Check gke/README.md.

1. Build images and deploy to GKE `make docker-build`
2. Use `make gke` to apply all manifests

To enable Gateway API:

```bash
gcloud container clusters update dwk-cluster --location=europe-north1-b --gateway-api=standard
```

## Deploy to cluster with service mesh k3s-iostio

```bash
make deploy-to-k3s-istio
```
