# GKE

## Create cluster

```bash
gcloud container clusters create dwk-cluster --zone=europe-north1-b --cluster-version=1.32 --disk-size=32 --num-nodes=3 --machine-type=e2-micro
```

```bash
# env variables
export PROJECT_ID=<add project id>
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
