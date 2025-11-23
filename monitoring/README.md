# Monitoring

If monitor is needed just run `make` it would run all needed command.

## Check logs via UI

1. Run Grafana and login

   ```bash
   # get password for admin:
   kubectl --namespace prometheus get secrets $(kubectl --namespace prometheus get secrets -oname | grep grafana | sed 's/secret\///g') -o jsonpath="{.data.admin-password}" | base64 -d ; echo

   # start
   export POD_NAME=$(kubectl --namespace prometheus get pod -l "app.kubernetes.io/name=grafana" -oname)
   kubectl --namespace prometheus port-forward $POD_NAME 3000
   ```

2. Add loki to data source:
   URL: `http://loki-gateway.loki-stack.svc.cluster.local`
   Header: `X-Scope-OrgId` value `grafana`
3. Explore -> loki -> Label browser -> find backend

## Manual installation

### Install helm

```bash
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-4
chmod 700 get_helm.sh
./get_helm.sh
```

### Install dependencies

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# prometheus stack
kubectl create namespace prometheus
helm install prometheus-community/kube-prometheus-stack --generate-name --namespace prometheus
```

### Install loki

```bash
kubectl create namespace loki-stack
helm install -n loki-stack --values loki-values.yaml loki grafana/loki
```

### Install grafana alloy

```bash
helm install -n prometheus --values alloy-values.yaml alloy grafana/alloy
```

### Add loki via [UI](http://localhost:3000/)

Use:

- http://loki-gateway.loki-stack.svc.cluster.local
- `X-Scope-OrgId` header to `grafana`
