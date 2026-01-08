
k3d-iostio:
	k3d cluster create --api-port 6550 -p '9080:80@loadbalancer' -p '9443:443@loadbalancer' --agents 2 --k3s-arg '--disable=traefik@server:*'
# 	helm install istio-cni istio/cni -n istio-system --set profile=ambient --set global.platform=k3d --wait
	istioctl install --set profile=ambient --set values.global.platform=k3d

k3d-iostio-apply-sample:
	kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.28/samples/bookinfo/platform/kube/bookinfo.yaml
	kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.28/samples/bookinfo/platform/kube/bookinfo-versions.yaml
# wait for pods to be ready
	sleep 30
	kubectl get pods
	kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.0/standard-install.yaml
	kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.28/samples/bookinfo/gateway-api/bookinfo-gateway.yaml
	kubectl annotate gateway bookinfo-gateway networking.istio.io/service-type=ClusterIP --namespace=default
# wait for gateway to get address
	sleep 30
	kubectl get gateway
# add namespace to ambient mesh
	kubectl label namespace default istio.io/dataplane-mode=ambient
# create curl pod
	kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.28/samples/curl/curl.yaml

kiali:
	istioctl dashboard kiali

install-kiali:
	kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.28/samples/addons/prometheus.yaml
	kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.28/samples/addons/kiali.yaml

port-forward-iostio:
	kubectl port-forward svc/bookinfo-gateway-istio 8080:80

send-100requests:
	@echo "Sending 100 requests..."
	@for i in $$(seq 1 100); do \
		curl -sSI -o /dev/null http://localhost:8080/productpage; \
	done

apply-auth-policy:
	kubectl apply -f auth-policy.yaml

test-curl-pod:
	kubectl exec deploy/curl -- curl -s "http://productpage:9080/productpage"

enable-l7:
	istioctl waypoint apply --enroll-namespace --wait
	kubectl get gtw waypoint
	kubectl apply -f  auth-policy-l7.yaml
	kubectl apply -f  auth-policy-l4.yaml

test-auth:
	kubectl exec deploy/curl -- curl -s "http://productpage:9080/productpage" -X DELETE || true
	kubectl exec deploy/reviews-v1 -- curl -s http://productpage:9080/productpage || true
	kubectl exec deploy/curl -- curl -s http://productpage:9080/productpage | grep -o "<title>.*</title>" || true

split-traffic:
	kubectl apply -f routing.yaml

test-traffic-split:
	kubectl exec deploy/curl -- sh -c "for i in \$$(seq 1 100); do curl -s http://productpage:9080/productpage | grep reviews-v.-; done"
