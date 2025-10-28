## Log output app

### Deploy to k3d cluster

1. Import image `k3d image import log-output-app`
2. Deploy everything `kubectl apply -f manifests`
3. Find pod with `kubectl get pods`
4. Get logs with `kubectl logs -f <pod-name>`
