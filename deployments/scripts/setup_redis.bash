#!/bin/bash

set -euo pipefail

NAMESPACE=databases
if [[ "$(kubectl get namespace | grep -c $NAMESPACE)" != "1" ]]; then
		kubectl create namespace $NAMESPACE
fi

case $ENV in

  "dev")
		helm upgrade --install redis -n $NAMESPACE \
			bitnami/redis \
			-f ./deployments/helm/databases/redis/values.yaml
    ;;

  "staging")
		helm upgrade --install redis -n $NAMESPACE \
			bitnami/redis \
			-f ./deployments/helm/databases/redis/values-staging.yaml
    ;;

  *)
		echo "not register env $ENV"
    ;;
esac

