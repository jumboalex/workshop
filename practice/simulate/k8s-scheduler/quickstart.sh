#!/bin/bash

set -e

echo "╔════════════════════════════════════════════════════════╗"
echo "║   Custom Kubernetes Scheduler - Quick Start           ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl not found. Please install kubectl first."
    exit 1
fi

# Check if connected to cluster
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Cannot connect to Kubernetes cluster. Please check your kubeconfig."
    exit 1
fi

echo "✅ Connected to Kubernetes cluster"
echo ""

# Build the scheduler
echo "📦 Building the scheduler..."
go build -o my-scheduler main.go
echo "✅ Build successful"
echo ""

# Option to deploy to cluster or run locally
echo "How would you like to run the scheduler?"
echo "1. Run locally (development)"
echo "2. Deploy to cluster (production)"
read -p "Choice (1-2): " choice

if [ "$choice" == "1" ]; then
    echo ""
    echo "🚀 Starting scheduler locally..."
    echo "Press Ctrl+C to stop"
    echo ""
    ./my-scheduler
elif [ "$choice" == "2" ]; then
    echo ""
    echo "🚀 Deploying scheduler to cluster..."

    # Build Docker image
    read -p "Enter your Docker registry (e.g., docker.io/username): " registry
    image_name="${registry}/my-scheduler:latest"

    echo "Building Docker image: ${image_name}"
    docker build -t ${image_name} .

    echo "Pushing to registry..."
    docker push ${image_name}

    # Update deployment YAML
    sed -i "s|your-registry/my-scheduler:latest|${image_name}|g" deployment.yaml

    echo "Deploying to Kubernetes..."
    kubectl apply -f deployment.yaml

    echo ""
    echo "✅ Deployment complete!"
    echo ""
    echo "Check status with:"
    echo "  kubectl get pods -n kube-system -l app=my-scheduler"
    echo ""
    echo "View logs with:"
    echo "  kubectl logs -n kube-system -l app=my-scheduler -f"
else
    echo "Invalid choice"
    exit 1
fi

echo ""
echo "📝 To test the scheduler, create a pod with:"
echo "  kubectl apply -f example-pod.yaml"
echo ""
echo "Monitor scheduling with:"
echo "  kubectl get events --watch"
