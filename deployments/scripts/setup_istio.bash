#!/bin/bash

set -euo pipefail

ROOT=~/.mm-printing
if [[ ! -d $ROOT ]]; then
  mkdir -p $ROOT
fi

export ISTIO_VERSION=1.12.2
ISTIO_VERSION_NAME=${ISTIO_VERSION//./-}

install_istio() {
  ISTIO_DIR=$ROOT/istio-$ISTIO_VERSION
  if [[ ! -d $ISTIO_DIR ]]; then
    # run in subshell to keep the current directory
    (cd $ROOT && curl -L https://istio.io/downloadIstio | sh -)
  fi

  istioctl=$ISTIO_DIR/bin/istioctl
  $istioctl x precheck

  rev=$(kubectl -n istio-system get pod -l app=istio-ingressgateway -o "jsonpath={.items..metadata.labels['istio\.io/rev']}")
  if [[ "$rev" != "${ISTIO_VERSION//./-}" ]]; then
    echo "Installing istio version $ISTIO_VERSION"
    echo " $(pwd)"
    $istioctl install -y -f ./deployments/istio/$ENV/config-${ISTIO_VERSION_NAME}.yaml --set revision=${ISTIO_VERSION_NAME}
  fi
}

create_backend_namespace() {
  if [[ "$(kubectl get namespace | grep -c $NAMESPACE)" != "1" ]]; then
      kubectl create namespace $NAMESPACE
  fi

  istio_rev="$(kubectl get namespace $NAMESPACE -o "jsonpath={.metadata.labels['istio\.io/rev']}")"
  if [[ -z $istio_rev ]]; then
    kubectl label namespace $NAMESPACE istio.io/rev=${ISTIO_VERSION_NAME}
  elif [[ $istio_rev != "$ISTIO_VERSION_NAME" ]]; then
    kubectl label namespace $NAMESPACE istio.io/rev-
    kubectl label namespace $NAMESPACE istio.io/rev=${ISTIO_VERSION_NAME}
  fi
}

setup_cert_manager() {
  if [[ "$(kubectl get namespace | grep -c cert-manager)" != "1" ]]; then
      kubectl create namespace cert-manager
  fi

  if [[ "$(helm repo ls | grep -c jetstack)" != "1" ]]; then
      helm repo add jetstack https://charts.jetstack.io
      helm repo update
  fi

  helm upgrade --install cert-manager jetstack/cert-manager \
    --wait \
    --namespace cert-manager \
    --version v1.7.1 \
    --set installCRDs=true

	kubectl -n cert-manager apply -f ./deployments/cert-manager/issuer.yaml
}

install_istio_gateway() {
  helm upgrade --wait -n istio-system --install --timeout 1m30s \
    ${ENV}-${ORG}-gateway ./deployments/helm/platforms/gateway \
    --values ./deployments/helm/platforms/gateway/${ENV}-${ORG}-values.yaml
}

install_istio_ratelimit() {
	kubectl -n istio-system apply -f  ./deployments/istio/dev/rate-limit-service.yaml
	kubectl -n istio-system apply -f  ./deployments/istio/dev/filter-rate-limit.yaml
}
