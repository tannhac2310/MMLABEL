#!/bin/bash

set -euo pipefail

NAMESPACE=databases
if [[ "$(kubectl get namespace | grep -c $NAMESPACE)" != "1" ]]; then
		kubectl create namespace $NAMESPACE
fi

case $ENV in

  "dev")
		kubectl -n $NAMESPACE apply -f ./deployments/helm/pv/mongodb.yaml
		helm upgrade --install mongodb -n $NAMESPACE \
			bitnami/mongodb \
			-f ./deployments/helm/databases/mongodb/values.yaml
    ;;

  "staging")
		kubectl -n $NAMESPACE apply -f ./deployments/helm/pv/managed-premium-retain.yaml
		helm upgrade --install mongodb -n $NAMESPACE \
			bitnami/mongodb \
			-f ./deployments/helm/databases/mongodb/values-staging.yaml
    ;;

  *)
		echo "not register env $ENV"
    ;;
esac

# shared
#helm upgrade --install mongodb -n $NAMESPACE \
	#-f ./deployments/helm/databases/mongodb/mongo-shared.yaml \
	#bitnami/mongodb-sharded


