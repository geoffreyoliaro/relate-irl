# Deployment Checklist - Relate IRL

Follow these steps in order to deploy Relate IRL to production.

## Pre-Deployment (Done in This Session)

- [x] FalkorDB Cloud credentials configured
- [x] Go backend with Xano integration implemented
- [x] Frontend dashboard UI created (login, directory, network, intelligence)
- [x] GitHub Actions CI/CD workflow configured
- [x] API documentation complete
- [x] Integration tests added
- [x] Dark theme and professional styling applied

## Step 1: Setup GCP Workload Identity Federation

Time: 5 minutes

```bash
# 1a. Make setup script executable
chmod +x scripts/setup-gcp.sh

# 1b. Run setup script
./scripts/setup-gcp.sh

# 1c. Note the output values for GitHub Secrets:
# - GCP_WORKLOAD_IDENTITY_PROVIDER
# - GCP_SERVICE_ACCOUNT
```

## Step 2: Configure GitHub Repository

Time: 5 minutes

### 2a. Verify repository URL

- [ ] Repository: https://github.com/geoffreyoliaro/relate-irl

### 2b. Add GitHub Secrets

Go to: **Settings → Secrets and variables → Actions → New repository secret**

Add all these secrets:

```
GCP_WORKLOAD_IDENTITY_PROVIDER: (from setup script output)
GCP_SERVICE_ACCOUNT: relate-irl-github@stable-glass-256211.iam.gserviceaccount.com
FALKORDB_HOST: deploy_instance_1.falkordb.io
FALKORDB_USER: oliaro
FALKORDB_PASSWORD: Pasnipop
JWT_SECRET: ad7f81a603abc09b48547866a0e178a039f2b5db871e7809b910a05fb60ce8f8
XANO_API_BASE: https://amara.xano.io/api:v1
```

Verify by going to **Settings → Secrets and variables → Actions** and confirming all 7 secrets are present.

## Step 3: Deploy via GitHub Actions

Time: 10 minutes

### 3a. Commit all changes

```bash
git add -A
git commit -m "feat: Deploy Relate IRL with CI/CD pipeline"
```

### 3b. Push to master

```bash
git push origin master
```

### 3c. Monitor deployment

Go to: **GitHub → Actions tab**

Watch the workflow:
1. **test** job — Go tests run
2. **build-and-push** job — Docker builds and pushes to Artifact Registry
3. **deploy** job — Deploys to Cloud Run
4. **notify** job — Sends final status

All jobs should complete in 3-5 minutes.

### 3d. Check Cloud Run deployment

```bash
gcloud run services describe relate-irl-api \
  --region=us-central1 \
  --format='value(status.url)'
```

Example output:
```
https://relate-irl-api-abc123.run.app
```

## Step 4: Verify Deployment

Time: 5 minutes

### 4a. Health check

Replace `{SERVICE_URL}` with your URL from step 3d:

```bash
curl {SERVICE_URL}/healthz
```

Expected: `{"status":"ok"}`

### 4b. Login test

```bash
curl -X POST {SERVICE_URL}/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"amara@horizonvc.com","password":"demo123"}'
```

Expected: JWT token + user data

### 4c. View logs

```bash
gcloud run services logs read relate-irl-api \
  --region=us-central1 \
  --limit=20
```

Should show successful initialization and connections.

### 4d. Test frontend (if separate)

Frontend is currently part of Next.js app in same repo. Deploy frontend separately if needed.

## Post-Deployment

### Enable additional monitoring

```bash
# View all revisions
gcloud run services revisions list relate-irl-api --region=us-central1

# Check current traffic split
gcloud run services describe relate-irl-api --region=us-central1 --format=json | jq '.status.traffic'

# View real-time logs with follow
gcloud run services logs read relate-irl-api --follow --region=us-central1
```

### Optional: Setup custom domain

```bash
gcloud run services update relate-irl-api \
  --region=us-central1 \
  --update-env-vars="CUSTOM_DOMAIN=api.relate-irl.com"

# Then add DNS CNAME record
# api.relate-irl.com CNAME ghs.googlehosted.com
```

### Optional: Enable Cloud Armor

Add DDoS protection and security policies:

```bash
# Create security policy
gcloud compute security-policies create relate-irl-policy \
  --description="DDoS and abuse protection"

# Apply to Cloud Run
gcloud beta run services update relate-irl-api \
  --security-policy=relate-irl-policy \
  --region=us-central1
```

## Troubleshooting

### Workflow fails at "Authenticate to Google Cloud"

**Cause**: GitHub Secrets not set correctly

**Fix**:
1. Verify all 7 secrets are set in GitHub
2. Re-run setup script: `./scripts/setup-gcp.sh`
3. Update GitHub secrets with output

### Deployment fails: "Service account not found"

**Cause**: Service account created but not propagated yet

**Fix**: Wait 30 seconds and retry push, or manually run:

```bash
git commit --allow-empty -m "Retry deployment"
git push origin master
```

### Health check fails

**Cause**: Backend not starting or FalkorDB unreachable

**Fix**:
1. Check logs: `gcloud run services logs read relate-irl-api --limit=50`
2. Verify env vars: `gcloud run services describe relate-irl-api --format=json | jq '.spec.template.spec.containers[0].env'`
3. Test FalkorDB connectivity locally first

### Frontend shows "Login failed"

**Cause**: API URL not correct or CORS issue

**Fix**:
1. Update `.env.production`: `NEXT_PUBLIC_API_URL={SERVICE_URL}`
2. Verify backend health: `curl {SERVICE_URL}/healthz`
3. Check browser console for actual error

## Rollback Procedure

If issues arise, revert to previous version:

```bash
# List revisions
gcloud run services revisions list relate-irl-api --region=us-central1

# Rollback to previous revision
gcloud run services update-traffic relate-irl-api \
  --to-revisions PREVIOUS_REVISION_NAME=100 \
  --region=us-central1
```

## Success Criteria

Deployment is successful when:

- [ ] GitHub Actions workflow completes without errors
- [ ] Cloud Run service shows "Active" status
- [ ] Health check returns 200 OK
- [ ] Login endpoint works with demo credentials
- [ ] API returns JWT token
- [ ] No errors in Cloud Run logs
- [ ] Service is publicly accessible

## Performance Targets

Expected performance after deployment:

- Health check latency: <100ms
- Login latency: <500ms (includes Xano call)
- API list endpoint: <200ms
- Intelligence queries: <500ms

Monitor at: `gcloud run services describe relate-irl-api --format=json | jq '.status.latestReadyRevision'`

## Next Steps After Deployment

1. **Monitor**: Set up Cloud Monitoring dashboards
2. **Alerts**: Configure alerts for errors and latency
3. **Scale**: Adjust min/max instances based on traffic
4. **Domain**: Add custom domain if needed
5. **Security**: Enable Cloud Armor and API keys
6. **Analytics**: Setup Application Insights or similar

## Support

For issues, check:
- Cloud Run logs: `gcloud run services logs read relate-irl-api`
- GitHub Actions: GitHub Actions tab
- GCP Console: https://console.cloud.google.com
- API Documentation: `/docs/API.md`

---

**Estimated Total Time**: 25 minutes

**Status**: Ready to deploy
