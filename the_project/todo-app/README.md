# Todo app

## Deploy to k3d cluster

1. Import image `k3d image import todo-app-app`
2. Deploy `kubectl create deployment todo-app-dep --image=todo-app-app`
3. Find pod with `kubectl get pods`
4. Get logs with `kubectl logs -f <pod-name>`
