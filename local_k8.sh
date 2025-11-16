#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="gopherify"
IMAGE_NAME="users:local"

echo "Starting local Kubernetes dev environment..."

# 1. Start Minikube if it's not running
if ! minikube status &>/dev/null; then
  echo "Starting Minikube..."
  minikube start --driver=docker
else
  echo "Minikube already running."
fi

# 2. Enable NGINX Ingress
echo "Enabling NGINX ingress..."
minikube addons enable ingress

# 3. Use Minikube Docker daemon for local builds
echo "Using Minikube's Docker environment..."
eval "$(minikube docker-env)"

# 4. Build backend image
echo "Building local backend image..."
docker build -t ${IMAGE_NAME} -f backend/cmd/service/Dockerfile backend

# 5. Apply dev overlay
echo "Applying Kustomize dev overlay..."
kubectl apply -k k8s/dev

# 6. Wait for pods to become ready
echo "Waiting for pods in namespace '${NAMESPACE}' to be ready..."
kubectl wait --for=condition=available --timeout=180s deployment --all -n "${NAMESPACE}" || true

# 7. Show services
echo "Current services in namespace '${NAMESPACE}':"
kubectl get svc -n "${NAMESPACE}"

# 8. Show ingress info
echo "Ingress routes:"
kubectl get ingress -n "${NAMESPACE}"

INGRESS_HOST=$(kubectl get ingress users-ingress -n "${NAMESPACE}" -o jsonpath='{.spec.rules[0].host}' || echo "")
if [ -n "$INGRESS_HOST" ]; then
  echo ""
  echo "Backend should be accessible at:"
  echo "  http://${INGRESS_HOST}/api/v1/healthz"
  echo ""
  echo "Tip: Run 'sudo minikube tunnel' in another terminal if not already running."
else
  echo "Ingress host not found. Ensure the ingress resource exists and has a host configured."
fi