#!/bin/bash

set -euo pipefail

NAMESPACE=monitoring
if [[ "$(kubectl get namespace | grep -c $NAMESPACE)" != "1" ]]; then
		kubectl create namespace $NAMESPACE
fi

case $ENV in

  "dev")
		helm upgrade --install -n monitoring grafana \
			grafana/grafana \
			--set=persistence.enabled=true

		helm install -n monitoring jaeger-tracing ./deployments/helm/monitoring/jaeger/
		helm upgrade --install -n monitoring prometheus prometheus-community/prometheus
		helm upgrade --install -n monitoring loki ./ --set "fluent-bit.enabled=true,promtail.enabled=false,loki.persistence.enabled=true,loki.persistence.size=20Gi"
		helm install -n monitoring pyroscope ./deployments/helm/monitoring/pyroscope/
    ;;

  "staging")
  	helm upgrade --install -n monitoring grafana \
			grafana/grafana \
			--set=persistence.enabled=true \
			--set=persistence.storageClassName=managed-premium-retain \
			--set=persistence.size=20Gi

		helm install -n monitoring jaeger-tracing ./deployments/helm/monitoring/jaeger/ \
			-f ./deployments/helm/monitoring/jaeger/values-staging.yaml
		helm upgrade --install -n monitoring prometheus prometheus-community/prometheus \
			--set "server.persistentVolume.size=20Gi,server.persistentVolume.storageClass=managed-premium-retain" \
			--set "alertmanager.enabled=false"
		helm upgrade --install -n monitoring loki ./deployments/helm/monitoring/loki-stack \
			-f ./deployments/helm/monitoring/loki-stack/values-staging.yaml \
			--set "fluent-bit.enabled=true,promtail.enabled=false,loki.persistence.enabled=true,loki.persistence.size=50Gi,loki.persistence.storageClassName=managed-premium-retain"
		helm install -n monitoring pyroscope ./deployments/helm/monitoring/pyroscope/
    ;;

  *)
		echo "not register env $ENV"
    ;;
esac

