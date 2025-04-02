# URL Shortener Deployment Script (Windows)
# Created: 2025-04-02
# Author: abhisheksharm-3

# Stop on any error
$ErrorActionPreference = "Stop"

Write-Host "🚀 Starting URL Shortener deployment..." -ForegroundColor Cyan

# 1. Start Minikube if not running
try {
    minikube status | Out-Null
    Write-Host "✅ Minikube is already running" -ForegroundColor Green
}
catch {
    Write-Host "📦 Starting Minikube cluster..." -ForegroundColor Yellow
    minikube start
}

# 2. Enable Ingress addon
Write-Host "🔌 Enabling Ingress addon..." -ForegroundColor Yellow
minikube addons enable ingress

# 3. Point Docker to Minikube's Docker daemon
Write-Host "🔄 Pointing Docker CLI to Minikube's Docker daemon..." -ForegroundColor Yellow
& minikube -p minikube docker-env --shell powershell | Invoke-Expression
Write-Host "👉 Now using Docker inside Minikube" -ForegroundColor Green

# 4. Build Docker images
Write-Host "🏗️ Building server Docker image..." -ForegroundColor Yellow
docker build -t shrtn-server:v1 --target production -f server/Dockerfile .

Write-Host "🏗️ Building client Docker image..." -ForegroundColor Yellow
docker build -t shrtn-client:v1 --target production -f client/Dockerfile .

# 5. Create namespace (using existing namespace.yaml)
Write-Host "🌐 Creating namespace..." -ForegroundColor Yellow
kubectl apply -f kubernetes/namespace.yaml

# 6. Switch to the namespace
kubectl config set-context --current --namespace=shrtn

# 7. Apply ConfigMap
Write-Host "⚙️ Creating ConfigMap..." -ForegroundColor Yellow
kubectl apply -f kubernetes/configmap.yaml

# 8. Create Secret from .env file
Write-Host "🔐 Creating Secret from .env file..." -ForegroundColor Yellow
if (Test-Path .env) {
    $secretExists = kubectl get secret appwrite-credentials -n shrtn --ignore-not-found
    if ($secretExists) {
        Write-Host "Secret 'appwrite-credentials' already exists. Deleting it first..." -ForegroundColor Yellow
        kubectl delete secret appwrite-credentials -n shrtn
    }
    kubectl create secret generic appwrite-credentials --from-env-file=.env -n shrtn
    Write-Host "✅ Secret created successfully" -ForegroundColor Green
} else {
    Write-Host "❌ .env file not found! Please create a .env file with your Appwrite credentials" -ForegroundColor Red
    exit 1
}

# 9. Deploy server
Write-Host "🔙 Deploying server..." -ForegroundColor Yellow
kubectl apply -f kubernetes/backend-deployment.yaml

# 10. Deploy client
Write-Host "🖥️ Deploying client..." -ForegroundColor Yellow
kubectl apply -f kubernetes/frontend-deployment.yaml

# 11. Wait for deployments to be ready
Write-Host "⏳ Waiting for deployments to be ready..." -ForegroundColor Yellow
kubectl rollout status deployment/shrtn-server -n shrtn
kubectl rollout status deployment/shrtn-client -n shrtn

# 12. Create Ingress
Write-Host "🌍 Creating Ingress..." -ForegroundColor Yellow
kubectl apply -f kubernetes/ingress.yaml

# 13. Apply HPA if it exists
if (Test-Path kubernetes/hpa.yaml) {
    Write-Host "⚖️ Applying Horizontal Pod Autoscaler..." -ForegroundColor Yellow
    kubectl apply -f kubernetes/backend-hpa.yaml
}

# 14. Set up local hosts file (requires running PowerShell as Administrator)
$minikubeIp = minikube ip
$hostsFile = "$env:windir\System32\drivers\etc\hosts"
$hostsContent = Get-Content $hostsFile
if ($hostsContent -notcontains "$minikubeIp shrtn.local") {
    try {
        Write-Host "📝 Updating hosts file..." -ForegroundColor Yellow
        Add-Content -Path $hostsFile -Value "`n$minikubeIp shrtn.local" -ErrorAction Stop
        Write-Host "✅ Hosts file updated successfully" -ForegroundColor Green
    } 
    catch {
        Write-Host "❗ Could not update hosts file. Please run PowerShell as Administrator or manually add this line to $hostsFile:" -ForegroundColor Red
        Write-Host "$minikubeIp shrtn.local" -ForegroundColor Yellow
    }
} else {
    Write-Host "✅ Hosts file entry already exists" -ForegroundColor Green
}

# 15. Print final instructions
Write-Host "`n✅ Deployment complete!" -ForegroundColor Green
Write-Host "🌐 Access your URL shortener at http://shrtn.local" -ForegroundColor Cyan
Write-Host "📊 Check resources with: kubectl get all -n shrtn" -ForegroundColor Cyan
Write-Host "📝 View server logs with: kubectl logs -f -l app=shrtn-server -n shrtn" -ForegroundColor Cyan
Write-Host "📝 View client logs with: kubectl logs -f -l app=shrtn-client -n shrtn" -ForegroundColor Cyan