#!/bin/bash

# GCP Deployment Setup Script
# This script sets up Workload Identity Federation for GitHub Actions CI/CD
# Usage: ./scripts/setup-gcp.sh

set -e

# Configuration
GCP_PROJECT_ID="stable-glass-256211"
GCP_REGION="us-central1"
SERVICE_ACCOUNT_NAME="relate-irl-github"
WORKLOAD_POOL_NAME="github-pool"
WORKLOAD_PROVIDER_NAME="github-provider"
GITHUB_REPO="geoffreyoliaro/relate-irl"
SERVICE_NAME="relate-irl-api"

echo "================================================"
echo "GCP Deployment Setup"
echo "================================================"
echo "Project ID: $GCP_PROJECT_ID"
echo "Region: $GCP_REGION"
echo "GitHub Repo: $GITHUB_REPO"
echo ""

# Step 1: Set project
echo "[1/7] Setting GCP project..."
gcloud config set project $GCP_PROJECT_ID
gcloud config set compute/region $GCP_REGION

# Step 2: Create service account
echo "[2/7] Creating service account..."
if gcloud iam service-accounts describe "$SERVICE_ACCOUNT_NAME@$GCP_PROJECT_ID.iam.gserviceaccount.com" &>/dev/null; then
  echo "Service account already exists, skipping creation"
else
  gcloud iam service-accounts create $SERVICE_ACCOUNT_NAME \
    --display-name="GitHub Actions for relate-irl" \
    --project=$GCP_PROJECT_ID
fi

SERVICE_ACCOUNT_EMAIL="$SERVICE_ACCOUNT_NAME@$GCP_PROJECT_ID.iam.gserviceaccount.com"
echo "Service Account: $SERVICE_ACCOUNT_EMAIL"

# Step 3: Grant permissions
echo "[3/7] Granting IAM permissions..."
gcloud projects add-iam-policy-binding $GCP_PROJECT_ID \
  --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
  --role="roles/run.developer" \
  --condition=None

gcloud projects add-iam-policy-binding $GCP_PROJECT_ID \
  --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
  --role="roles/artifactregistry.writer" \
  --condition=None

# Step 4: Create Workload Identity Pool
echo "[4/7] Creating Workload Identity Pool..."
if gcloud iam workload-identity-pools describe $WORKLOAD_POOL_NAME \
  --project=$GCP_PROJECT_ID \
  --location=global &>/dev/null; then
  echo "Workload Identity Pool already exists, skipping creation"
else
  gcloud iam workload-identity-pools create $WORKLOAD_POOL_NAME \
    --project=$GCP_PROJECT_ID \
    --location=global \
    --display-name="GitHub Actions Pool"
fi

# Step 5: Create Workload Identity Provider
echo "[5/7] Creating Workload Identity Provider..."
if gcloud iam workload-identity-pools providers describe $WORKLOAD_PROVIDER_NAME \
  --project=$GCP_PROJECT_ID \
  --location=global \
  --workload-identity-pool=$WORKLOAD_POOL_NAME &>/dev/null; then
  echo "Workload Identity Provider already exists, skipping creation"
else
  gcloud iam workload-identity-pools providers create-oidc $WORKLOAD_PROVIDER_NAME \
    --project=$GCP_PROJECT_ID \
    --location=global \
    --workload-identity-pool=$WORKLOAD_POOL_NAME \
    --display-name="GitHub Provider" \
    --attribute-mapping="google.subject=assertion.sub,assertion.aud=assertion.aud,assertion.repository=assertion.repository" \
    --issuer-uri="https://token.actions.githubusercontent.com" \
    --attribute-condition="assertion.repository_owner == 'geoffreyoliaro'"
fi

# Step 6: Configure service account impersonation
echo "[6/7] Configuring service account impersonation..."
gcloud iam service-accounts add-iam-policy-binding $SERVICE_ACCOUNT_EMAIL \
  --project=$GCP_PROJECT_ID \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/$GCP_PROJECT_ID/locations/global/workloadIdentityPools/$WORKLOAD_POOL_NAME/providers/$WORKLOAD_PROVIDER_NAME/attributes.repository/$GITHUB_REPO"

# Step 7: Get Workload Identity Provider resource name
echo "[7/7] Getting Workload Identity Provider resource name..."
WORKLOAD_IDENTITY_PROVIDER=$(gcloud iam workload-identity-pools providers describe $WORKLOAD_PROVIDER_NAME \
  --project=$GCP_PROJECT_ID \
  --location=global \
  --workload-identity-pool=$WORKLOAD_POOL_NAME \
  --format="value(name)")

echo ""
echo "================================================"
echo "Setup Complete!"
echo "================================================"
echo ""
echo "Add these as GitHub Secrets:"
echo "1. GCP_WORKLOAD_IDENTITY_PROVIDER:"
echo "   $WORKLOAD_IDENTITY_PROVIDER"
echo ""
echo "2. GCP_SERVICE_ACCOUNT:"
echo "   $SERVICE_ACCOUNT_EMAIL"
echo ""
echo "Steps to add secrets:"
echo "1. Go to: https://github.com/geoffreyoliaro/relate-irl/settings/secrets/actions"
echo "2. Click 'New repository secret'"
echo "3. Add each secret above"
echo ""
echo "Then, push to master or develop branch to trigger deployment:"
echo "  git push origin master"
echo ""
