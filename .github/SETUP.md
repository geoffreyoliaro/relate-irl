# GitHub Actions Setup Guide

This guide walks you through setting up GitHub Actions for automated CI/CD deployment.

## Prerequisites

- GitHub repository at: https://github.com/geoffreyoliaro/relate-irl.git
- GCP Project: `stable-glass-256211`
- FalkorDB Cloud instance: `deploy_instance_1`
- Xano instance configured

## Step 1: Create GitHub Secrets

Go to: **GitHub Repository → Settings → Secrets and variables → Actions**

Add these secrets:

```
GCP_WORKLOAD_IDENTITY_PROVIDER: (see GCP setup below)
GCP_SERVICE_ACCOUNT: (see GCP setup below)
FALKORDB_HOST: deploy_instance_1.falkordb.io
FALKORDB_USER: oliaro
FALKORDB_PASSWORD: Pasnipop
JWT_SECRET: ad7f81a603abc09b48547866a0e178a039f2b5db871e7809b910a05fb60ce8f8
XANO_API_BASE: https://amara.xano.io/api:v1_abcd1234
```

## Step 2: GCP Workload Identity Federation Setup

This allows GitHub Actions to authenticate securely without long-lived service account keys.

### 2a. Create Service Account

```bash
gcloud iam service-accounts create relate-irl-github \
  --display-name "GitHub Actions for relate-irl" \
  --project=stable-glass-256211
```

### 2b. Grant Cloud Run deployment permissions

```bash
gcloud projects add-iam-policy-binding stable-glass-256211 \
  --member="serviceAccount:relate-irl-github@stable-glass-256211.iam.gserviceaccount.com" \
  --role="roles/run.developer"

gcloud projects add-iam-policy-binding stable-glass-256211 \
  --member="serviceAccount:relate-irl-github@stable-glass-256211.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"
```

### 2c. Create Workload Identity Provider

```bash
gcloud iam workload-identity-pools create "github-pool" \
  --project="stable-glass-256211" \
  --location="global" \
  --display-name="GitHub Actions Pool"
```

### 2d. Create Workload Identity Provider credential

```bash
gcloud iam workload-identity-pools providers create-oidc "github-provider" \
  --project="stable-glass-256211" \
  --location="global" \
  --workload-identity-pool="github-pool" \
  --display-name="GitHub Provider" \
  --attribute-mapping="google.subject=assertion.sub,assertion.aud=assertion.aud,assertion.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com"
```

### 2e. Create service account binding

```bash
gcloud iam service-accounts add-iam-policy-binding \
  relate-irl-github@stable-glass-256211.iam.gserviceaccount.com \
  --project="stable-glass-256211" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/stable-glass-256211/locations/global/workforcePools/github-pool/providers/github-provider/attributes.repository/geoffreyoliaro/relate-irl"
```

Wait, the last command uses workforce pools. Let me use the correct attribute condition:

```bash
gcloud iam service-accounts add-iam-policy-binding \
  relate-irl-github@stable-glass-256211.iam.gserviceaccount.com \
  --project="stable-glass-256211" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/stable-glass-256211/locations/global/workloadIdentityPools/github-pool/providers/github-provider/attributes.repository/geoffreyoliaro/relate-irl"
```

### 2f. Get Workload Identity Provider resource name

```bash
gcloud iam workload-identity-pools providers describe github-provider \
  --project="stable-glass-256211" \
  --location="global" \
  --workload-identity-pool="github-pool" \
  --format="value(name)"
```

This output should be set as `GCP_WORKLOAD_IDENTITY_PROVIDER` secret in GitHub.

### 2g. Get Service Account email

```bash
echo "relate-irl-github@stable-glass-256211.iam.gserviceaccount.com"
```

This should be set as `GCP_SERVICE_ACCOUNT` secret in GitHub.

## Step 3: Verify Setup

1. Go to GitHub Repository → Actions
2. Push a commit to `develop` or `master` branch
3. Watch the workflow run in the Actions tab
4. Check deployment status

## Workflow Behavior

- **PR to master/develop**: Run tests only
- **Push to develop**: Run tests → Build → Push to Artifact Registry
- **Push to master**: Run tests → Build → Push → Deploy to Cloud Run

## Monitoring

View logs: GitHub Actions tab → Select workflow → Click run

View deployments: `gcloud run services describe relate-irl-api --region=us-central1`

View Cloud Run logs: `gcloud run services logs read relate-irl-api --region=us-central1 --limit=50`

## Troubleshooting

**Workflow not triggering**: Push to master or develop branch

**Authentication failed**: Verify GCP_WORKLOAD_IDENTITY_PROVIDER and GCP_SERVICE_ACCOUNT secrets

**Build failed**: Check Go tests pass locally: `go test ./...`

**Deployment failed**: Check Cloud Run quotas and service account permissions
