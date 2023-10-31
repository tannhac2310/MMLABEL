#!/bin/bash

# exist if any script failed
set -euxo pipefail

command -v kubectl
command -v helm

# environment variables for local env
export ENV=${ENV:-staging}
export ORG=${ORG:-mm-printing}
export TAG=${TAG:-v0.1.80-rc5}
export NAMESPACE=${NAMESPACE:-backend}

kubectl apply -f ./deployments/helm/pv/managed-premium-retain.yaml

### Working with your local k8s cluster
# Install some platform services
. ./deployments/scripts/setup_istio.bash

install_istio
install_istio_ratelimit
create_backend_namespace
setup_cert_manager

# Start all backend service in k8s

#install_istio_gateway

./deployments/scripts/setup_cockroach.bash
#./deployments/scripts/setup_influxdb.bash
./deployments/scripts/setup_redis.bash
# setup backend
#./deployments/scripts/setup_backend.bash

