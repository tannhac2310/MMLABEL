#!/bin/bash

# exist if any script failed
set -euxo pipefail

command -v kubectl
command -v helm

validateInput() {
	if [ -z "$1" ]
		then
			echo "version can not be empty: usage: ./deployments/scripts/manual_deploy.bash v0.0.1"
	fi
}

validateInput $1
VERSION=$1

# create a tag
git tag $VERSION
git push origin $VERSION

DOCKER_IMAGE="192.168.88.33:32000/mm-printing/backend"
# build docker image
docker build -t ${DOCKER_IMAGE}:${VERSION} .
docker push ${DOCKER_IMAGE}:${VERSION}

# deploy
export KUBECONFIG=./deployments/config.dev.yaml
helm upgrade --install mm-printing \
        ./deployments/helm/backend/ \
        -n backend \
        --set=hydra.image.tag=${VERSION} \
        --set=appVersion=${VERSION}

