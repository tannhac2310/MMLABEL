#!/bin/bash

set -euo pipefail

NAMESPACE=nats-streaming
if [[ "$(kubectl get namespace | grep -c $NAMESPACE)" != "1" ]]; then
		kubectl create namespace $NAMESPACE
fi

case $ENV in

  "dev")
#		helm upgrade --install nats-streaming -n $NAMESPACE \
#			./deployments/helm/nats-streaming/
    echo "not register env $ENV"
    ;;

  "staging")
		helm upgrade --install nats-streaming -n $NAMESPACE \
			./deployments/helm/nats-streaming/ \
			-f ./deployments/helm/nats-streaming/values-staging.yaml
    ;;

  *)
		echo "not register env $ENV"
    ;;
esac

