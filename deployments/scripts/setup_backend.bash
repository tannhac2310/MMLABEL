#!/bin/bash

set -euo pipefail

case $ENV in

  "dev")
		helm upgrade --install mm-printing \
    			./deployments/helm/backend \
    			-f ./deployments/helm/backend/values-dev.yaml \
    			 -n backend-dev \
    			--set=hydra.env=${ENV} \
    			--set=aurora.env=${ENV} \
    			--set=iot.env=${ENV} \
    			--set=hydra.image.tag=${TAG} \
    			--set=aurora.image.tag=${TAG} \
    			--set=iot.image.tag=${TAG} \
    			--set=appVersion=${TAG}
    ;;

  "staging")
		helm upgrade --install mm-printing \
        ./deployments/helm/backend \
        -f ./deployments/helm/backend/values-staging.yaml \
         -n backend \
        --set=hydra.env=${ENV} \
        --set=aurora.env=${ENV} \
        --set=iot.env=${ENV} \
        --set=hydra.image.tag=${TAG} \
        --set=aurora.image.tag=${TAG} \
        --set=iot.image.tag=${TAG} \
        --set=appVersion=${TAG}
    ;;
  *)
		echo "not register env $ENV"
    ;;
esac


#ENV=staging TAG=v1.6.58 ./deployments/scripts/setup_backend.bash