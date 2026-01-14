clean-part1:
	kubectl delete -f manifests/part1/ || true

clean-part2:
	kubectl delete -f manifests/part2/ || true

clean: clean-part1 clean-part2
	@echo "Cleaned up all resources."

deploy-part1: clean
	kubectl apply -f manifests/part1/

deploy-part2: clean
	kubectl apply -f manifests/part2/

deploy-part3: clean
	kubectl apply -f manifests/part2/dep.yaml
	kubectl apply -f manifests/part3/

deploy-part4: clean
	kubectl apply -f manifests/part2/dep.yaml
	kubectl apply -f manifests/part4/

forward-ingressgateway:
	kubectl port-forward svc/istio-ingressgateway 8080:80 -n istio-ingress

docker:
	sed -i "/res.send/c\  res.send('Hello, World\!');" app/index.js
	docker build -t hello:1 ./app
	sed -i "/res.send/c\  res.send('Hello, New World\!');" app/index.js
	docker build -t hello:2 ./app
	sed -i "/res.send/c\  res.send('Hello, World\!');" app/index.js
	k3d image import -c k3s-iostio hello:1
	k3d image import -c k3s-iostio hello:2

load:
	@echo "Sending 100 requests..."
	@for i in $$(seq 1 100); do \
		echo "";\
		curl http://localhost:9080/; \
	done
	@echo "Done."
