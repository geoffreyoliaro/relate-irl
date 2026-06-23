# GitHub Repository Setup Guide

This guide walks you through setting up the GitHub repository with all CI/CD workflows and best practices.

## Prerequisites

- GitHub account with repository created
- Repository cloned locally
- Admin access to repository settings

## Step 1: Repository Configuration

### Enable Branch Protection Rules

1. Go to **Settings → Branches**
2. Add a branch protection rule for `master`:
   - ✅ Require a pull request before merging
   - ✅ Require status checks to pass:
     - `Frontend CI / Lint & Type Check`
     - `Frontend CI / Build Frontend`
     - `Frontend CI / Run Tests`
     - `Test & Build` (from deploy.yml)
     - `Code Quality Analysis / Code Quality Checks`
     - `Security & Dependency Scanning / Dependency Scanning`
   - ✅ Require code reviews (minimum: 1)
   - ✅ Dismiss stale pull request approvals on push
   - ✅ Require branches to be up to date before merging
   - ✅ Require conversation resolution before merging

### For `develop` branch:
Same as master, but can be less strict (e.g., 0 required reviewers if team is small)

## Step 2: GitHub Secrets Configuration

### GCP Secrets (Required for Deployment)

1. Go to **Settings → Secrets and variables → Actions**
2. Add these repository secrets:

```
GCP_WORKLOAD_IDENTITY_PROVIDER
  → From: gcloud iam workload-identity-pools providers describe <PROVIDER> \
           --project=<PROJECT> \
           --location=global \
           --workload-identity-pool=github-pool

GCP_SERVICE_ACCOUNT
  → Email of your GCP service account
```

### Optional Secrets

```
CODECOV_TOKEN
  → From: codecov.io (for coverage reports)

PRODUCTION_API_URL
  → Your production API endpoint for health checks
```

## Step 3: GitHub Actions Permissions

1. Go to **Settings → Actions → General**
2. Configure:
   - **Actions permissions:** Allow all actions and reusable workflows
   - **Fork pull request workflows:** Choose your security level
   - **Workflow permissions:** 
     - ✅ Read and write permissions
     - ✅ Allow GitHub Actions to create and approve pull requests

## Step 4: Code Scanning Setup

### Enable CodeQL

1. Go to **Security → Code scanning → Set up code scanning**
2. Select **CodeQL analysis**
3. Choose **Default** configuration
4. Click **Enable CodeQL**

### Enable Dependabot

1. Go to **Settings → Code security and analysis**
2. Enable:
   - ✅ Dependabot alerts
   - ✅ Dependabot security updates
   - ✅ Grouped security updates (optional)

## Step 5: Configure Issue & PR Templates

The templates are automatically configured:
- `.github/ISSUE_TEMPLATE/bug_report.md`
- `.github/ISSUE_TEMPLATE/feature_request.md`
- `.github/pull_request_template.md`

These will automatically appear when creating issues/PRs.

## Step 6: Configure CODEOWNERS

Update `.github/CODEOWNERS` with actual GitHub usernames:

```
# Replace @github_username with actual usernames
* @username1 @username2
```

## Step 7: Enable Issue Templates

1. Go to **Settings → General**
2. Enable:
   - ✅ Issues
   - ✅ Discussions (optional)
   - ✅ Wikis

3. Go to **Settings → Issue templates**
4. Verify custom issue templates appear

## Step 8: Set Up Environments (Optional)

1. Go to **Settings → Environments**
2. Create environments for:
   - `development`
   - `staging`
   - `production`

3. Configure protection rules if needed (requires reviewers, etc.)

## Step 9: Configure Notifications

1. Go to **Settings → Notifications**
2. Choose notification preferences for:
   - Workflow runs
   - Pull request reviews
   - Issues
   - Discussions

## Step 10: Verify Workflows

1. Go to **Actions** tab
2. Verify all workflows are present:
   - ✅ Backend tests & deployment (`deploy.yml`)
   - ✅ Frontend CI (`frontend-ci.yml`)
   - ✅ Security & Dependency Scanning (`security.yml`)
   - ✅ PR Checks (`pr-checks.yml`)
   - ✅ Code Quality Analysis (`code-quality.yml`)
   - ✅ Release Management (`release.yml`)
   - ✅ Post-Deployment Monitoring (`monitoring.yml`)

## First Deployment

1. **Create a test branch:**
   ```bash
   git checkout -b feature/test-workflow
   ```

2. **Make a small change:**
   ```bash
   echo "# Test" >> README.md
   git add README.md
   git commit -m "feat: test workflow"
   git push -u origin feature/test-workflow
   ```

3. **Create a Pull Request** and verify:
   - ✅ All workflows run successfully
   - ✅ Status checks pass
   - ✅ Required reviewers are notified

4. **Merge the PR** and verify deployment workflow runs

## Troubleshooting

### Workflows Not Running

1. Check that workflows are in `.github/workflows/`
2. Verify branch is not protected from actions
3. Check repository permissions: **Settings → Actions → General**
4. Review workflow syntax: [GitHub Actions Syntax](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)

### Status Checks Not Appearing

1. Verify workflow file syntax (YAML format)
2. Check that job names match branch protection rule names
3. Re-run workflow after fixing syntax errors
4. Update branch protection rules if job names changed

### Secrets Not Available

1. Verify secrets are in **Settings → Secrets and variables → Actions**
2. Confirm secrets are referenced correctly in workflows (`${{ secrets.SECRET_NAME }}`)
3. Note: Secrets are not available in fork PRs for security reasons

### Deployment Failures

1. Check GCP credentials are valid
2. Verify service account has required permissions
3. Review Cloud Run deployment limits and quotas
4. Check application logs in GCP Console

## Maintenance

### Weekly Tasks
- Review failed workflow runs
- Check for Dependabot alerts

### Monthly Tasks
- Update workflow dependencies
- Review security scan results
- Check for deprecation warnings

### Quarterly Tasks
- Audit GitHub secrets and access
- Review branch protection rules
- Update CI/CD configuration

## Documentation

Key documents in this repository:
- `.github/DEVOPS.md` - Complete CI/CD pipeline documentation
- `.github/SECURITY.md` - Security policies and best practices
- `.github/CODEOWNERS` - Code ownership configuration
- `.github/workflows/*` - Individual workflow definitions

## Support

For issues with:
- **GitHub Actions:** [GitHub Actions Documentation](https://docs.github.com/en/actions)
- **GCP Deployment:** [Google Cloud Run Docs](https://cloud.google.com/run/docs)
- **Repository Setup:** Create an issue in this repository

---

**Next Steps:**
1. Follow all steps above
2. Create a test PR to verify workflows
3. Share this document with your team
4. Reference `.github/DEVOPS.md` for operational procedures

**Last Updated:** June 23, 2026
