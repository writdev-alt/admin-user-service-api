# Google Cloud Run Deployment Script (PowerShell)
# Usage: .\deploy.ps1 [-ServiceName <name>] [-Region <region>] [-ProjectId <project-id>] [-EnvVars <env-vars>]
# Example: .\deploy.ps1 -EnvVars "SERVER_SECRET=secret,DATABASE_HOST=host,DATABASE_DBNAME=dbname"

param(
    [string]$ServiceName = "admin-user-service",
    [string]$Region = "asia-southeast1",
    [string]$ProjectId = "ascendant-quest-475706-g6",
    [string]$EnvVars = ""
)

# Set error action preference
$ErrorActionPreference = "Stop"

# Get project ID from gcloud if not provided
if ([string]::IsNullOrEmpty($ProjectId)) {
    try {
        $ProjectId = gcloud config get-value project 2>$null
        if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrEmpty($ProjectId)) {
            Write-Host "Error: PROJECT_ID is required. Set it using -ProjectId parameter or configure gcloud default project." -ForegroundColor Red
            Write-Host "Usage: .\deploy.ps1 [-ServiceName <name>] [-Region <region>] [-ProjectId <project-id>]" -ForegroundColor Yellow
            exit 1
        }
    }
    catch {
        Write-Host "Error: Failed to get project ID from gcloud. Please set it using -ProjectId parameter." -ForegroundColor Red
        exit 1
    }
}

Write-Host "Deploying to Google Cloud Run..." -ForegroundColor Cyan
Write-Host "Project: $ProjectId" -ForegroundColor Green
Write-Host "Service: $ServiceName" -ForegroundColor Green
Write-Host "Region: $Region" -ForegroundColor Green
Write-Host ""

# Generate a unique tag (using timestamp)
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
# Artifact Registry format: REGION-docker.pkg.dev/PROJECT_ID/REPOSITORY/IMAGE
# gcloud artifacts repositories create admin-user-service-api --repository-format=docker --location=asia-southeast1 --description="Docker repository for admin-user-service-api"
$repository = "admin-user-service-api"   # Artifact Registry repo; create with: gcloud artifacts repositories create REPO_NAME --repository-format=docker --location=$Region
$imageBase = "${Region}-docker.pkg.dev/$ProjectId/$repository/$ServiceName"
$imageTagWithTimestamp = "$imageBase`:$timestamp"
$imageTagLatest = "$imageBase`:latest"

# Step 1: Build Docker image
Write-Host "Building Docker image..." -ForegroundColor Cyan
docker build -t $imageTagWithTimestamp -t $imageTagLatest .

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Docker build failed!" -ForegroundColor Red
    exit $LASTEXITCODE
}

# Step 2: Configure Docker to use gcloud as credential helper
Write-Host "Configuring Docker authentication..." -ForegroundColor Cyan
gcloud auth configure-docker --quiet

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Docker authentication failed!" -ForegroundColor Red
    exit $LASTEXITCODE
}

# Step 3: Push Docker image to Container Registry
Write-Host "Pushing Docker image to Container Registry..." -ForegroundColor Cyan
docker push $imageTagWithTimestamp

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Docker push failed!" -ForegroundColor Red
    exit $LASTEXITCODE
}

docker push $imageTagLatest

if ($LASTEXITCODE -ne 0) {
    Write-Host "Warning: Failed to push latest tag, continuing..." -ForegroundColor Yellow
}

# Step 4: Deploy to Cloud Run
Write-Host "Deploying to Cloud Run..." -ForegroundColor Cyan

# Build the deployment command
if (-not [string]::IsNullOrEmpty($EnvVars)) {
    Write-Host "Setting environment variables..." -ForegroundColor Cyan
    gcloud run deploy $ServiceName `
        --image $imageTagWithTimestamp `
        --region $Region `
        --platform managed `
        --allow-unauthenticated `
        --port 8080 `
        --memory 512Mi `
        --cpu 1 `
        --min-instances 0 `
        --max-instances 10 `
        --timeout 300 `
        --set-env-vars $EnvVars `
        --project $ProjectId
} else {
    gcloud run deploy $ServiceName `
        --image $imageTagWithTimestamp `
        --region $Region `
        --platform managed `
        --allow-unauthenticated `
        --port 8080 `
        --memory 512Mi `
        --cpu 1 `
        --min-instances 0 `
        --max-instances 10 `
        --timeout 300 `
        --project $ProjectId
}

if ($LASTEXITCODE -ne 0) {
    Write-Host "" -ForegroundColor Red
    Write-Host "Error: Cloud Run deployment failed!" -ForegroundColor Red
    Write-Host "" -ForegroundColor Yellow
    Write-Host "Note: Make sure required environment variables are set:" -ForegroundColor Yellow
    Write-Host "  - SERVER_SECRET (required)" -ForegroundColor Yellow
    Write-Host "  - DATABASE_DRIVER (required)" -ForegroundColor Yellow
    Write-Host "  - DATABASE_DBNAME (required)" -ForegroundColor Yellow
    Write-Host "  - DATABASE_USERNAME (required)" -ForegroundColor Yellow
    Write-Host "  - DATABASE_PASSWORD (required)" -ForegroundColor Yellow
    Write-Host "  - DATABASE_HOST (required)" -ForegroundColor Yellow
    Write-Host "" -ForegroundColor Yellow
    Write-Host "Set them using: -EnvVars 'KEY1=value1,KEY2=value2'" -ForegroundColor Yellow
    Write-Host "Or update them after deployment using Cloud Console or gcloud CLI" -ForegroundColor Yellow
    exit $LASTEXITCODE
}

Write-Host ""
Write-Host "Deployment complete!" -ForegroundColor Green
Write-Host "Service URL: https://$ServiceName-$Region-$ProjectId.a.run.app" -ForegroundColor Cyan
Write-Host ""
Write-Host "View logs: gcloud run services logs read $ServiceName --region=$Region" -ForegroundColor Yellow
Write-Host "View service: gcloud run services describe $ServiceName --region=$Region" -ForegroundColor Yellow

