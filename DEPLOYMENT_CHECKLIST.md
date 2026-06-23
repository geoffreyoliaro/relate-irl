# Deployment Checklist

## Information You Need to Provide

### GitHub
- [ ] **GitHub username** — e.g., `octocat`
- [ ] **GitHub repository name** — e.g., `mininexus`
- [ ] **GitHub repository URL** — e.g., `https://github.com/octocat/mininexus`

### GCP
- [ ] **GCP Project ID** — e.g., `mininexus-prod` (must be globally unique)
- [ ] **GCP Project Number** — (auto-generated, shown in GCP Console)
- [ ] **GCP Region** — e.g., `us-central1`, `us-east1`, `europe-west1`
- [ ] **Billing Account ID** — (for enabling APIs)

### Xano (Authentication)
- [ ] **Xano API Base URL** — e.g., `https://your-instance.xano.io/api:12345`
- [ ] **Xano Instance Name** — for reference

### Security
- [ ] **JWT_SECRET** — Generate with: `openssl rand -hex 32` (save this securely!)

### FalkorDB (Already have ✅)
- [x] **FalkorDB Host** — `deploy_instance_1.falkordb.io`
- [x] **FalkorDB User** — `oliaro`
- [x] **FalkorDB Password** — `Pasnipop`
- [x] **FalkorDB Instance ID** — `instance-kd24k3j41`

---

## Setup Steps

### Step 1: Prepare Your Environment
- [ ] Fork/clone this repository to your GitHub account
- [ ] Generate JWT_SECRET: `openssl rand -hex 32`
- [ ] Have your Xano API Base URL ready

### Step 2: GCP Project Setup (v0 can help automate)
- [ ] Create GCP Project (or provide existing project ID)
- [ ] Enable billing on GCP Project
- [ ] Create Artifact Registry repository
- [ ] Set up Workload Identity Federation
- [ ] Create service account for GitHub Actions

### Step 3: GitHub Secrets Configuration
Add these 8 secrets to GitHub repo → Settings → Secrets and variables:

| Secret | Example Value |
|--------|---|
| `GCP_PROJECT_ID` | `mininexus-prod` |
| `GCP_WORKLOAD_IDENTITY_PROVIDER` | `projects/123456789/locations/global/workloadIdentityPools/github-pool/providers/github-provider` |
| `GCP_SERVICE_ACCOUNT` | `github-actions-deployer@mininexus-prod.iam.gserviceaccount.com` |
| `FALKORDB_HOST` | `deploy_instance_1.falkordb.io` |
| `FALKORDB_PASSWORD` | `Pasnipop` |
| `JWT_SECRET` | (your generated 64-char hex string) |
| `XANO_API_BASE` | `https://your-instance.xano.io/api:xxxx` |
| `GCP_REGION` | `us-central1` |

### Step 4: Deploy
- [ ] Push code to GitHub: `git push origin main`
- [ ] Watch GitHub Actions workflow run
- [ ] Verify Cloud Run deployment
- [ ] Test API endpoints

---

## Quick Deploy Command (Manual)

Once all secrets are configured, you can manually deploy with:

```bash
gcloud run deploy mininexus-api \
  --image us-central1-docker.pkg.dev/YOUR_PROJECT_ID/mininexus/api:latest \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars "FALKORDB_HOST=deploy_instance_1.falkordb.io,FALKORDB_USER=oliaro,FALKORDB_PASSWORD=Pasnipop,JWT_SECRET=YOUR_JWT_SECRET,XANO_API_BASE=YOUR_XANO_URL,GIN_MODE=release"
```

---

## Support

For detailed instructions, see **DEPLOYMENT_GUIDE.md**

Need help? Common issues in the guide:
- Workload Identity setup
- Docker image push failures
- FalkorDB connection timeouts
- API authentication errors
