all: ns

ns:
	kubectl create namespace default || true
	kubens default

install: ns
	helm install -f myvalues.yaml my-nats oci://registry-1.docker.io/bitnamicharts/nats

update: ns
	helm upgrade -f myvalues.yaml my-nats oci://registry-1.docker.io/bitnamicharts/nats

uninstall: ns
	helm uninstall my-nats

monitor:
	kubectl port-forward --namespace default svc/my-nats 8222:8222

deploy-app:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes-hy/material-example/refs/heads/master/app9/manifests/deployments.yaml

metrics:
	kubectl port-forward --namespace default svc/my-nats-metrics 7777:7777
