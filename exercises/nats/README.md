# NATS

Check `Makefile` for more info

## Install NATS with helm

```bash
make install
```

## Start monitoring

```bash
make monitor
```

## Deploy apps

```bash
make deploy-app
```

## Label nats prometheus metrics

```bash
# get label
kubectl get prometheus -n prometheus -o yaml | grep serviceMonitorSelector -A 5

# apply label
kubectl label servicemonitors.monitoring.coreos.com -n prometheus my-nats-metrics release=kube-prometheus-stack-1764690198
```
