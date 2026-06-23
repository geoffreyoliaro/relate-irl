# MiniNexus Deployment Guide (GitHub + GCP)

This guide walks you through deploying the MiniNexus API to **GCP Cloud Run** with **GitHub Actions** CI/CD.

---

## Prerequisites

You already have:
- ✅ FalkorDB Cloud instance: `deploy_instance_1` (oliaro/Pasnipop)
- ✅ Go API code in GitHub

You need to provide:
1. **GitHub repository** (must be connected to this project)
2. **GCP Project** (must have billing enabled)
3. **Xano API Base URL** (for authentication)
4. **A generated JWT_SECRET** (security token)

---

## Step 1: Prepare GitHub Repository

### 1.1 Create GitHub repository (if not already done)
```bash
cd /path/to/mininexus
git init
git add .
git commit -m "Initial commit: MiniNexus API"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/mininexus.git
git push -u origin main
```

### 1.2 Generate JWT_SECRET
Run this command and save the output — you'll need it in Step 3:
```bash
openssl rand -hex 32
# Example output: a7f9e2c1d4b6c8e3f5a9d2b1c4e7f0a3d6b9c2e5f8a1d4c7b0e3f6a9d2c5
```

---

## Step 2: Set Up GCP Project

### 2.1 Create GCP Project (if you don't have one)
```bash
# Set your project ID
export GCP_PROJECT_ID="mininexus-prod"

# Create project
gcloud projects create $GCP_PROJECT_ID --name="MiniNexus API"

# Set as default
gcloud config set project $GCP_PROJECT_ID

# Enable billing (REQUIRED — get billing account ID from GCP Console)
gcloud billing projects link $GCP_PROJECT_ID \
  --billing-account=BILLING_ACCOUNT_ID
```

### 2.2 Enable required APIs
```bash
gcloud services enable \
  artifactregistry.googleapis.com \
  cloudrun.googleapis.com \
  iam.googleapis.com \
  iamcredentials.googleapis.com
```

### 2.3 Create Artifact Registry repository
```bash
gcloud artifacts repositories create mininexus \
  --repository-format=docker \
  --location=us-central1 \
  --description="MiniNexus API Docker images"
```

### 2.4 Set up Workload Identity Federation
This allows GitHub Actions to securely deploy without storing long-lived keys.

```bash
export GCP_PROJECT_ID="mininexus-prod"
export WORKLOAD_IDENTITY_PROVIDER="projects/$(gcloud projects describe $GCP_PROJECT_ID --format='value(projectNumber)')/locations/global/workloadIdentityPools/github-pool/providers/github-provider"

# Create the identity pool (run once)
gcloud iam workload-identity-pools create "github-pool" \
  --project="$GCP_PROJECT_ID" \
  --location="global" \
  --display-name="GitHub Actions Pool"

# Create the identity provider (run once)
gcloud iam workload-identity-pools providers create-oidc "github-provider" \
  --project="$GCP_PROJECT_ID" \
  --location="global" \
  --workload-identity-pool="github-pool" \
  --display-name="GitHub Actions Provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com"
```

### 2.5 Create service account for Cloud Run deployment
```bash
export GCP_PROJECT_ID="mininexus-prod"
export SERVICE_ACCOUNT="github-actions-deployer"

# Create service account
gcloud iam service-accounts create $SERVICE_ACCOUNT \
  --project=$GCP_PROJECT_ID \
  --display-name="GitHub Actions Cloud Run Deployer"

# Grant necessary roles
gcloud projects add-iam-policy-binding $GCP_PROJECT_ID \
  --member="serviceAccount:${SERVICE_ACCOUNT}@${GCP_PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding $GCP_PROJECT_ID \
  --member="serviceAccount:${SERVICE_ACCOUNT}@${GCP_PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/run.admin"

gcloud iam service-accounts add-iam-policy-binding \
  "${SERVICE_ACCOUNT}@${GCP_PROJECT_ID}.iam.gserviceaccount.com" \
  --project="$GCP_PROJECT_ID" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/$(gcloud projects describe $GCP_PROJECT_ID --format='value(projectNumber)')/locations/global/workloadIdentityPools/github-pool/attribute.repository/YOUR_GITHUB_USERNAME/mininexus"
```

---

## Step 3: Add GitHub Secrets

Go to your GitHub repository → **Settings** → **Secrets and variables** → **Actions** → **New repository secret**

Add these 8 secrets:

| Secret Name | Value | Source |
|---|---|---|
| `GCP_PROJECT_ID` | `mininexus-prod` | From Step 2.1 |
| `GCP_WORKLOAD_IDENTITY_PROVIDER` | Full provider URI (from Step 2.4 output) | From Step 2.4 |
| `GCP_SERVICE_ACCOUNT` | `github-actions-deployer@mininexus-prod.iam.gserviceaccount.com` | From Step 2.5 |
| `FALKORDB_HOST` | `deploy_instance_1.falkordb.io` | Already configured |
| `FALKORDB_PASSWORD` | `Pasnipop` | Already configured |
| `JWT_SECRET` | Your generated secret from Step 1.2 | Generate once, reuse |
| `XANO_API_BASE` | `https://your-instance.xano.io/api:xxxx` | Get from Xano dashboard |
| `GCP_REGION` | `us-central1` | (or your preferred region) |

---

## Step 4: Create GitHub Actions Workflow

Create `.github/workflows/deploy.yml`:

```yaml
name: Build & Deploy to Cloud Run

on:
  push:
    branches: [main]
  workflow_dispatch:

env:
  REGISTRY: us-central1-docker.pkg.dev
  IMAGE_NAME: mininexus

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    permissions:
      contents: read
      id-token: write

    steps:
      - uses: actions/checkout@v4

      - name: Authenticate to Google Cloud
        uses: google-github-actions/auth@v2
        with:
          workload_identity_provider: ${{ secrets.GCP_WORKLOAD_IDENTITY_PROVIDER }}
          service_account: ${{ secrets.GCP_SERVICE_ACCOUNT }}

      - name: Set up Cloud SDK
        uses: google-github-actions/setup-gcloud@v2

      - name: Configure Docker for Artifact Registry
        run: |
          gcloud auth configure-docker ${{ env.REGISTRY }}

      - name: Build & push Docker image
        run: |
          IMAGE_TAG="${{ env.REGISTRY }}/${{ secrets.GCP_PROJECT_ID }}/${{ env.IMAGE_NAME }}/api:${{ github.sha }}"
          docker build -t $IMAGE_TAG .
          docker push $IMAGE_TAG
          echo "IMAGE_TAG=${IMAGE_TAG}" >> $GITHUB_ENV

      - name: Deploy to Cloud Run
        run: |
          gcloud run deploy mininexus-api \
            --image ${{ env.IMAGE_TAG }} \
            --region ${{ secrets.GCP_REGION }} \
            --allow-unauthenticated \
            --set-env-vars "FALKORDB_HOST=${{ secrets.FALKORDB_HOST }},FALKORDB_USER=oliaro,FALKORDB_PASSWORD=${{ secrets.FALKORDB_PASSWORD }},JWT_SECRET=${{ secrets.JWT_SECRET }},XANO_API_BASE=${{ secrets.XANO_API_BASE }},GIN_MODE=release" \
            --memory 512Mi \
            --cpu 1 \
            --timeout 60 \
            --project ${{ secrets.GCP_PROJECT_ID }}

      - name: Smoke test
        run: |
          SERVICE_URL=$(gcloud run services describe mininexus-api \
            --region ${{ secrets.GCP_REGION }} \
            --project ${{ secrets.GCP_PROJECT_ID }} \
            --format 'value(status.url)')
          
          echo "Deployment URL: $SERVICE_URL"
          curl -f "$SERVICE_URL/health" || echo "Health check will be available once API is fully deployed"
```

---

## Step 5: Deploy

### Option A: Push to GitHub (Automatic)
```bash
git push origin main
# GitHub Actions will automatically build and deploy
```

**Monitor deployment:**
- Go to GitHub repo → **Actions** tab
- Click the latest workflow run
- Watch the build → push → deploy steps

### Option B: Manual Deploy (Testing)
```bash
export GCP_PROJECT_ID="mininexus-prod"
export GCP_REGION="us-central1"

gcloud run deploy mininexus-api \
  --image us-central1-docker.pkg.dev/${GCP_PROJECT_ID}/mininexus/api:latest \
  --region ${GCP_REGION} \
  --allow-unauthenticated \
  --set-env-vars "FALKORDB_HOST=deploy_instance_1.falkordb.io,FALKORDB_USER=oliaro,FALKORDB_PASSWORD=Pasnipop,JWT_SECRET=$(openssl rand -hex 32),XANO_API_BASE=YOUR_XANO_URL,GIN_MODE=release" \
  --memory 512Mi \
  --cpu 1 \
  --project ${GCP_PROJECT_ID}
```

---

## Step 6: Verify Deployment

Once deployed, you'll get a service URL like:
```
https://mininexus-api-xxxxxx-uc.a.run.app
```

Test it:
```bash
SERVICE_URL="https://mininexus-api-xxxxxx-uc.a.run.app"

# Login
curl -X POST "$SERVICE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"amara@horizonvc.com","password":"demo123"}'

# List people (with JWT token from login)
curl -X GET "$SERVICE_URL/api/v1/people" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Summary: What You Need to Provide

To complete the deployment, please provide:

1. ✅ **FalkorDB credentials** — Already have (oliaro/Pasnipop/deploy_instance_1)
2. ⏳ **GitHub repository URL** — Your GitHub username and repo name
3. ⏳ **GCP Project ID** — Your GCP project (or let me create one guide)
4. ⏳ **Xano API Base URL** — From your Xano dashboard
5. ⏳ **JWT_SECRET** — Generate using `openssl rand -hex 32`

Once you provide these, I'll update the secrets in your repository and you can push to GitHub to trigger automatic deployment!

---

## Troubleshooting

### "Cloud Run deployment failed: Permission denied"
- Verify service account has `roles/run.admin` and `roles/artifactregistry.writer`
- Check Workload Identity Federation provider configuration

### "Docker image push failed"
- Ensure `gcloud auth configure-docker` was run
- Verify `roles/artifactregistry.writer` is assigned

### "FalkorDB connection timeout"
- Check that `FALKORDB_HOST` is correct: `deploy_instance_1.falkordb.io`
- Verify Cloud Run has network access (usually automatic)
- Test locally first: `redis-cli -h deploy_instance_1.falkordb.io -a Pasnipop`

### "API returns 401 Unauthorized"
- Ensure `JWT_SECRET` matches between login and API calls
- Check that Xano URL is correct in `XANO_API_BASE`
