# URL Shortener Cleanup Script (Windows)
# Created: 2025-04-02
# Author: abhisheksharm-3

# Stop on any error
$ErrorActionPreference = "Stop"

Write-Host "🧹 Starting Minikube cleanup..." -ForegroundColor Cyan

# Stop Minikube
Write-Host "🛑 Stopping Minikube..." -ForegroundColor Yellow
minikube stop

# Delete Minikube
Write-Host "🗑️ Deleting Minikube..." -ForegroundColor Yellow
minikube delete

Write-Host "✅ Minikube cleanup complete!" -ForegroundColor Green
