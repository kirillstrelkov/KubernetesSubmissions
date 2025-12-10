# ArgoCD

> **NOTE**: versions of argo server and sidecar pluging should match! ArgoCD sidecar is need to add SOPS for decrypting manifests.

Check [./Makefile](./Makefile) for more information

Most apps support kustomize but to make it work with ArgoCD patch is needed:

```bash
# create secret
make upload-sops-key

# patch argocd server to add sops
make patch-repo-server
```

This will add ability for ArgoCD to decrypt encrypted files

## Check versions of ArgoCD containers

```bash
kubectl get pod -n argocd -l app.kubernetes.io/name=argocd-server -o jsonpath='{.items[*].status.containerStatuses[*].image}'

kubectl logs  -n argocd -l app.kubernetes.io/name=argocd-repo-server -c cmp-sops
```
