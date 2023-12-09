#!/bin/bash

set -euo pipefail

NAMESPACE=databases
if [[ "$(kubectl get namespace | grep -c $NAMESPACE)" != "1" ]]; then
		kubectl create namespace $NAMESPACE
fi

if [[ "$(helm repo ls | grep -c cockroachdb)" != "1" ]]; then
		helm repo add cockroachdb https://charts.cockroachdb.com/
		helm repo update
fi


case $ENV in

#  "dev")
#		helm -n $NAMESPACE upgrade --install cockroachdb cockroachdb/cockroachdb \
#			-f ./deployments/helm/databases/cockroachdb/values.yaml
#
#		curl -OOOOOOOOO \
#			https://raw.githubusercontent.com/cockroachdb/cockroach/master/cloud/kubernetes/client-secure.yaml
#
#		kubectl -n databases apply -f client-secure.yaml
#
#		kubectl -n databases exec -i cockroachdb-client-secure \
#			-- ./cockroach sql --certs-dir=/cockroach-certs --host=cockroachdb-public --execute="CREATE USER mm-printing WITH PASSWORD 'mm-printing';"
#		kubectl -n databases exec -i cockroachdb-client-secure \
#			-- ./cockroach sql --certs-dir=/cockroach-certs --host=cockroachdb-public --execute="GRANT admin TO mm-printing;"
#    ;;

  "staging")
		helm -n $NAMESPACE upgrade --install cockroachdb cockroachdb/cockroachdb \
    			-f ./deployments/helm/databases/cockroachdb/values-production.yaml

    kubectl -n databases apply -f client-secure-new.yaml

    kubectl -n databases exec -i cockroachdb-client-secure \
      -- ./cockroach sql --certs-dir=/cockroach-certs --host=cockroachdb-public --execute="CREATE USER mm-printing WITH PASSWORD 'M6763zUnr8tBdwyd';"
    kubectl -n databases exec -i cockroachdb-client-secure \
      -- ./cockroach sql --certs-dir=/cockroach-certs --host=cockroachdb-public --execute="GRANT admin TO mm-printing;"
      kubectl -n databases exec -i cockroachdb-client-secure \
      			-- ./cockroach sql --certs-dir=/cockroach-certs --host=cockroachdb-public --execute="GRANT admin TO tech4edu;"
    ;;

  "production")
		helm -n $NAMESPACE upgrade --install cockroachdb cockroachdb/cockroachdb \
			-f ./deployments/helm/databases/cockroachdb/values-production.yaml

		kubectl -n databases apply -f client-secure-new.yaml

		kubectl -n databases exec -i cockroachdb-client-secure \
			-- ./cockroach sql --certs-dir=/cockroach-certs --host=cockroachdb-public --execute="CREATE USER mm-printing WITH PASSWORD 'M6763zUnr8tBdwyd';"
		kubectl -n databases exec -i cockroachdb-client-secure \
			-- ./cockroach sql --certs-dir=/cockroach-certs --host=cockroachdb-public --execute="GRANT admin TO mm-printing;"
		;;
  *)
		echo "not register env $ENV"
    ;;
esac

# shared
#helm upgrade --install mongodb -n $NAMESPACE \
	#-f ./deployments/helm/databases/mongodb/mongo-shared.yaml \
	#bitnami/mongodb-sharded


