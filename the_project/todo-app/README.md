# Todo app

## Deploy to k3d cluster

1. Import image `k3d image import todo-app-app`
2. Deploy `kubectl apply -f manifests/deployment.yaml`
3. Find pod with `kubectl get pods`
4. Get logs with `kubectl logs -f <pod-name>`

### To forward port

```bash
kubectl port-forward <pod-name> 3000:8080
```
