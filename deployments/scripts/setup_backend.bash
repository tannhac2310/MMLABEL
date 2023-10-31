#!/bin/bash

set -euo pipefail

NAMESPACE=backend

case $ENV in

  "dev")
		helm upgrade --install mm-printing -n $NAMESPACE \
			./deployments/helm/backend \
			--set=env=${ENV} \
			--set=hydra.env=${ENV} \
			--set=gezu.env=${ENV} \
			--set=aurora.env=${ENV} \
			--set=hydra.image.tag=${TAG} \
			--set=gezu.image.tag=${TAG}
			--set=aurora.image.tag=${TAG}
    ;;

  "staging")
		helm upgrade --install mm-printing -n $NAMESPACE \
			./deployments/helm/backend \
			-f ./deployments/helm/backend/values-staging.yaml \
			--set=env=${ENV} \
			--set=hydra.env=${ENV} \
			--set=aurora.env=${ENV} \
			--set=hydra.image.tag=${TAG} \
			--set=aurora.image.tag=${TAG} \
			--set=appVersion=${TAG}

# helm upgrade --install mm-printing \
#        ./deployments/helm/backend/ \
#        -f ./deployments/helm/backend/values-staging.yaml \
#        -n backend \
#        --set=hydra.env=staging \
#        --set=gezu.env=staging \
#        --set=aurora.env=staging \
#        --set=hydra.image.tag=${CI_COMMIT_TAG} \
#        --set=gezu.image.tag=${CI_COMMIT_TAG} \
#        --set=aurora.image.tag=${CI_COMMIT_TAG} \
#        --set=appVersion=${CI_COMMIT_TAG}
    ;;

  *)
		echo "not register env $ENV"
    ;;
esac
