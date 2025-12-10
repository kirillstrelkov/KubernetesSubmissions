all: patch-argocd

ns:
	kubectl create namespace argocd || true
	kubens argocd

install: ns
	kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

get-credentials:
	kubectl get -n argocd secrets argocd-initial-admin-secret -o yaml | grep 'password' | awk '{{print $$2}}' | base64 -d
	@echo ""

upload-sops-key:
	kubectl create secret generic sops-age --from-file=key.txt=../key.txt -n argocd

patch-argocd: ns
	kubectl apply -f manifests/sops-cmp.yaml
	kubectl patch deployment argocd-repo-server -n argocd --patch-file manifests/argocd-repo-server-patch.yaml

sidecar-build:
	docker build -f Dockerfile -t custom-argocd-sidecar:v1 .
	k3d image import custom-argocd-sidecar:v1 -c k3s-default
