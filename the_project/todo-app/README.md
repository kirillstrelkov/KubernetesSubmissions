# Todo app

## Deploy to k3d cluster

1. Import image `k3d image import todo-app-app`
2. Deploy `kubectl apply -f manifests/deployment.yaml`
3. Find pod with `kubectl get pods`
4. Get logs with `kubectl logs -f <pod-name>`
5. Deploy service `kubectl apply -f manifests/service.yaml`

NOTE: create cluster with `k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2`
