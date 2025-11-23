all: monitoring

init-repos:
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update

create-namespaces:
	kubectl create namespace loki-stack
	kubectl create namespace prometheus

loki:
	helm install -n loki-stack --values loki-values.yaml loki grafana/loki

alloy:
	helm install -n prometheus --values alloy-values.yaml alloy grafana/alloy 

prometheus:
	helm install prometheus-community/kube-prometheus-stack --generate-name --namespace prometheus

monitoring: init-repos create-namespaces prometheus loki alloy
	kubectl --namespace prometheus get secrets $(kubectl --namespace prometheus get secrets -oname | grep grafana | sed 's/secret\///g') -o jsonpath="{.data.admin-password}" | base64 -d ; echo
	
