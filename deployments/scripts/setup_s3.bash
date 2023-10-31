#!/bin/bash

set -euo pipefail

NAMESPACE=backend
if [[ "$(kubectl get namespace | grep -c $NAMESPACE)" != "1" ]]; then
		kubectl create namespace $NAMESPACE
fi

if [[ "$(helm repo ls | grep -c minio)" != "1" ]]; then
		helm repo add minio https://helm.min.io/
		helm repo update
fi

case $ENV in

  "dev")
		helm upgrade --install s3-storage -n $NAMESPACE \
			minio/minio \
			-f ./deployments/helm/s3/values.yaml
    ;;

  "staging")
 		helm upgrade --install s3-storage -n $NAMESPACE \
			minio/minio \
			-f ./deployments/helm/s3/values-staging.yaml \
    ;;

  *)
		echo "not register env $ENV"
    ;;
esac

