# GKE

## Create cluster

```bash
gcloud container clusters create dwk-cluster --zone=europe-north1-b --cluster-version=1.32 --disk-size=32 --num-nodes=3 --machine-type=e2-micro
```

```bash
# env variables
export PROJECT_ID=dwk-gke-479211
export REGION=europe-north1
export DOCKER_PKG_DEV=$REGION-docker.pkg.dev
export DOCKER_REPO=my-docker-repo
export DOCKER_FULL_REPO=$DOCKER_PKG_DEV/${PROJECT_ID}/$DOCKER_REPO

gcloud artifacts repositories create $DOCKER_REPO \
   --repository-format=docker \
   --location=europe-north1 \
   --description="Docker repository"

gcloud auth configure-docker $DOCKER_PKG_DEV

gcloud artifacts repositories add-iam-policy-binding $DOCKER_REPO \
    --location=$REGION \
    --member=serviceAccount:170849643312-compute@developer.gserviceaccount.com \
    --role="roles/artifactregistry.reader"

```

## Add monitoring and NATS

```bash
# check gke cluster name
kubectx

# add monitoring
cd monitoring
kubectx <gke cluster>
kubens default
make monitoring

# add NATS
cd exercises/nats
kubectx <gke cluster>
kubens default
make install
```

## Add to ArgoCD

`argocd` cli should be preinstalled.

```bash
argocd login localhost:40917 --username admin --password <password for login> --insecure
argocd cluster add <gke cluster>
```
