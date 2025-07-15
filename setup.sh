#!/bin/bash

set -e

# Setup observability stack with kind and kube-prometheus

# Create kind cluster
kind create cluster --name observability

# Add Prometheus community Helm repository
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts

# Update Helm repositories
helm repo update

# Install kube-prometheus-stack chart
helm install kube-prometheus prometheus-community/kube-prometheus-stack --namespace monitoring --create-namespace

# Wait for pods to be ready
echo "Waiting for pods to be ready..."
kubectl --namespace monitoring wait --for=condition=ready pod --all --timeout=300s

# Get Grafana admin password
echo ""
echo "Grafana admin password:"
kubectl --namespace monitoring get secrets kube-prometheus-grafana -o jsonpath="{.data.admin-password}" | base64 -d
echo ""

# Instructions for accessing Grafana
echo ""
echo "To access Grafana UI:"
echo "1. Run: kubectl --namespace monitoring port-forward \$(kubectl --namespace monitoring get pod -l \"app.kubernetes.io/name=grafana,app.kubernetes.io/instance=kube-prometheus\" -o jsonpath=\"{.metadata.name}\") 3000:3000"
echo "2. Open: http://localhost:3000"
echo "3. Login with username: admin and the password shown above"