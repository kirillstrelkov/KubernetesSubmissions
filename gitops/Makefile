
install:
	kubectl create namespace argocd || true
	kubens argocd
	kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

get-credentials:
	kubectl get -n argocd secrets argocd-initial-admin-secret -o yaml | grep 'password' | awk '{{print $$2}}' | base64 -d
	echo
