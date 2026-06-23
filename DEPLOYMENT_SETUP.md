# Deployment Setup Guide for Relate IRL

## Quick Start

This guide walks through deploying Relate IRL to GCP Cloud Run with GitHub Actions CI/CD automation.

### Prerequisites

- GitHub account with repository at: `https://github.com/geoffreyoliaro/relate-irl`
- GCP Project: `stable-glass-256211`
- FalkorDB Cloud instance: `deploy_instance_1` (already configured)
- Xano instance with authentication endpoint

### Architecture

```
GitHub Push
    ↓
GitHub Actions Workflow
    ├─ Run Tests (Go unit tests)
    ├─ Build Docker Image
    ├─ Push to GCP Artifact Registry
    └─ Deploy to Cloud Run
    
Cloud Run Service
    ├─ Go API Backend
    ├─ FalkorDB Cloud Connection
    ├─ Xano Auth Integration
    └─ Health Checks
```

## Step 1: Setup GCP Workload Identity Federation

This allows GitHub Actions to authenticate securely without storing long-lived credentials.

### 1a. Run the setup script

```bash
cd /path/to/relate-irl
chmod +x scripts/setup-gcp.sh
./scripts/setup-gcp.sh
```

This script:
- Creates a service account: `relate-irl-github@stable-glass-256211.iam.gserviceaccount.com`
- Sets up Workload Identity Federation
- Configures trust between GitHub and GCP
- Outputs the secrets needed for GitHub

### 1b. Verify the setup

After running the script, you should see output like:

```
Add these as GitHub Secrets:
1. GCP_WORKLOAD_IDENTITY_PROVIDER:
   projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/github-pool/providers/github-provider

2. GCP_SERVICE_ACCOUNT:
   relate-irl-github@stable-glass-256211.iam.gserviceaccount.com
```

## Step 2: Configure GitHub Secrets

Go to: **GitHub Repository → Settings → Secrets and variables → Actions**

Click "New repository secret" and add these:

### Deployment Secrets (from GCP setup script output)

```
GCP_WORKLOAD_IDENTITY_PROVIDER: <from step 1b>
GCP_SERVICE_ACCOUNT: relate-irl-github@stable-glass-256211.iam.gserviceaccount.com
```

### FalkorDB Configuration

```
FALKORDB_HOST: deploy_instance_1.falkordb.io
FALKORDB_PORT: 6379
FALKORDB_USER: oliaro
FALKORDB_PASSWORD: Pasnipop
```

### JWT Secret

Generate a new secure secret:

```bash
openssl rand -hex 32
# Output: ad7f81a603abc09b48547866a0e178a039f2b5db871e7809b910a05fb60ce8f8
```

Add as secret:

```
JWT_SECRET: ad7f81a603abc09b48547866a0e178a039f2b5db871e7809b910a05fb60ce8f8
```

### Xano Configuration

```
XANO_API_BASE: https://amara.xano.io/api:v1
```

## Step 3: Test the CI/CD Pipeline

### 3a. Make a test push

```bash
git add .
git commit -m "Deploy: Initial setup"
git push origin master
```

### 3b. Monitor the workflow

Go to: **GitHub Repository → Actions**

Watch the workflow run:
1. **Test** — Go tests execute
2. **Build and Push** — Docker image builds and pushes to Artifact Registry
3. **Deploy** — Cloud Run deployment
4. **Health Check** — Verifies service is running

### 3c. Check deployment status

```bash
gcloud run services describe relate-irl-api \
  --region=us-central1 \
  --format='value(status.url)'
```

This outputs your live service URL, e.g.:
```
https://relate-irl-api-xxxxx.run.app
```

## Step 4: Verify End-to-End

### Health Check

```bash
curl https://relate-irl-api-xxxxx.run.app/healthz
```

Expected response:
```json
{"status":"ok"}
```

### Login Test

```bash
curl -X POST https://relate-irl-api-xxxxx.run.app/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"amara@horizonvc.com","password":"demo123"}'
```

Expected response:
```json
{
  "token": "eyJhbGc...",
  "expires_in": 3600,
  "user": {
    "id": 123,
    "name": "Amara",
    "email": "amara@horizonvc.com"
  }
}
```

## Workflow Details

### Branch Strategy

- **Pushes to `master`**: Full deployment (test → build → push → deploy)
- **Pushes to `develop`**: Build only (test → build → push, no deployment)
- **Pull Requests**: Tests only (no build or deployment)

### Environment Variables

Environment variables are injected at deployment time via `gcloud run deploy`:

```bash
--set-env-vars=FALKORDB_HOST=${{ secrets.FALKORDB_HOST }},\
                FALKORDB_PORT=6379,\
                FALKORDB_USER=${{ secrets.FALKORDB_USER }},\
                FALKORDB_PASSWORD=${{ secrets.FALKORDB_PASSWORD }},\
                JWT_SECRET=${{ secrets.JWT_SECRET }},\
                XANO_API_BASE=${{ secrets.XANO_API_BASE }}
```

### Auto-Scaling Configuration

The Cloud Run service is configured for:
- **CPU**: 1 vCPU
- **Memory**: 1 Gi
- **Timeout**: 300 seconds
- **Max Instances**: 10 (auto-scales based on traffic)
- **Min Instances**: 1 (always warm)

### Health Checks

Cloud Run automatically performs HTTP health checks every 10 seconds to:
- `/healthz` endpoint (defined in Go backend)
- Expects HTTP 200 response
- Restarts service if unhealthy

## Troubleshooting

### Workflow Fails at Auth Step

**Problem**: "Workload Identity authentication failed"

**Solution**:
1. Verify `GCP_WORKLOAD_IDENTITY_PROVIDER` secret is correct
2. Verify `GCP_SERVICE_ACCOUNT` secret is correct
3. Re-run setup script: `./scripts/setup-gcp.sh`

### Deployment Fails: "Service account not found"

**Problem**: "Service account 'relate-irl-github@...' not found"

**Solution**:
1. Verify service account exists: `gcloud iam service-accounts list --filter="email:relate-irl-github"`
2. If missing, re-run setup script

### Health Check Fails

**Problem**: "Health check failed after 5 attempts"

**Solution**:
1. Check service logs: `gcloud run services logs read relate-irl-api --limit=50`
2. Verify environment variables: `gcloud run services describe relate-irl-api --format=json | jq '.spec.template.spec.containers[0].env'`
3. Verify FalkorDB connection: Check firewall rules, credentials

### Image Push Failed

**Problem**: "Failed to push image to Artifact Registry"

**Solution**:
1. Verify Artifact Registry exists: `gcloud artifacts repositories list`
2. Create if missing: `gcloud artifacts repositories create relate-irl --repository-format=docker --location=us-central1`

## Monitoring & Logs

### View Live Logs

```bash
gcloud run services logs read relate-irl-api \
  --region=us-central1 \
  --limit=50 \
  --follow
```

### View Deployment History

```bash
gcloud run services revisions list relate-irl-api \
  --region=us-central1 \
  --format=table
```

### Check Service Status

```bash
gcloud run services describe relate-irl-api \
  --region=us-central1 \
  --format=json
```

## Rollback

If deployment has issues, rollback to previous version:

```bash
gcloud run services update-traffic relate-irl-api \
  --to-revisions REVISION_NAME=100 \
  --region=us-central1
```

Where `REVISION_NAME` is from the revisions list above.

## Cost Optimization

### Current Configuration

- Cloud Run: ~$5-10/month (first 2M requests free per month)
- Artifact Registry: ~$0.10/GB storage
- FalkorDB: Varies by plan (external service)

### Cost Reduction Tips

1. Set `--max-instances=5` instead of 10
2. Use `--min-instances=0` to scale to zero when idle (5-10 min idle time)
3. Monitor traffic and adjust CPU allocation

## Next Steps

1. Deploy by pushing to `master`: `git push origin master`
2. Monitor workflow at: GitHub Actions tab
3. Access your API at: Cloud Run service URL
4. Set up custom domain (optional)
5. Enable Cloud Armor for DDoS protection (optional)

## References

- GCP Cloud Run: https://cloud.google.com/run/docs
- GitHub Actions: https://docs.github.com/en/actions
- Workload Identity Federation: https://cloud.google.com/iam/docs/workload-identity-federation
- API Documentation: See `/docs/API.md`
