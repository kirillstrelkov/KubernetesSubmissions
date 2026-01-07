all: apply

ns:
	kubens exercises  || kubectl create namespace exercises

apply: ns
	kubectl apply -f manifests/

clean:
	kubectl delete -f manifests/ || true
