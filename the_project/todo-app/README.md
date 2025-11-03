# Todo app

## Deploy to k3d cluster

1. Build, import and deploy `make deploy`

NOTE: create cluster with `k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2`
NOTE 2: make folder in agent0 `docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/todoapp`
