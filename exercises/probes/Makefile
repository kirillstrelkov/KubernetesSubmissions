all: apply

ns: 
	kubens exercises || kubectl create namespace exercises

setup: ns
	kubectl create namespace argo-rollouts
	kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml

apply: ns
	kubectl apply -n exercises -f .
