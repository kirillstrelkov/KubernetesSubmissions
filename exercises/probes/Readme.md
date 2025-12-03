# Probes

Check `Makefile` for more information

## Argo rollouts

```bash
kubectl create namespace argo-rollouts
kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml
```

## Argo kubernetes pluging

### Installation

<https://argoproj.github.io/argo-rollouts/installation/>

### Monitoring

```bash
# start in terminal
kubectl argo rollouts get rollout flaky-update-dep --watch

# start web interface
kubectl argo rollouts dashboard
```
