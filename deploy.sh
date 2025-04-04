# Start Minikube if not already running
minikube status || minikube start

# Enable the Ingress addon
minikube addons enable ingress

# Set up Docker environment to use Minikube's Docker daemon
eval $(minikube docker-env)

# Build client image using Minikube's Docker daemon
echo "Building client image..."
docker build -t shrtn-client -f ./client/Dockerfile --target development .

# Build server image using Minikube's Docker daemon
echo "Building server image..."
docker build -t shrtn-server -f ./server/Dockerfile --target development ./server

# Load environment variables from .env file
echo "Loading environment variables..."
if [ -f .env ]; then
  # This properly exports all variables from .env file
  set -a
  source .env
  set +a
fi

# Deploy using command line flags for sensitive data
echo "Deploying Helm chart..."
helm upgrade --install shrtn ./shrtn-charts \
  --set server.env.APPWRITE_ENDPOINT="$APPWRITE_ENDPOINT" \
  --set server.env.APPWRITE_PROJECT_ID="$APPWRITE_PROJECT_ID" \
  --set server.env.APPWRITE_API_KEY="$APPWRITE_API_KEY" \
  --set server.env.APPWRITE_COLLECTION_ID="$APPWRITE_COLLECTION_ID" \
  --set server.env.APPWRITE_DATABASE_ID="$APPWRITE_DATABASE_ID"

# Get the Minikube IP
MINIKUBE_IP=$(minikube ip)
echo "Add the following entry to your /etc/hosts file:"
echo "$MINIKUBE_IP shrtn.local"

# Wait for the pods to be ready
echo "Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=shrtn --timeout=180s

echo "Deployment completed!"
echo "Access your application at: http://shrtn.local"