#!/bin/bash

set -euo pipefail

NAMESPACE=backend

kubectl -n $NAMESPACE apply -f ./deployments/helm/pv/s3.yaml