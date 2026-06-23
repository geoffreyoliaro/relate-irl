# DevOps & CI/CD Guide

This document outlines the CI/CD pipeline, deployment processes, and operational best practices for this project.

## Overview

The project uses GitHub Actions for continuous integration and deployment, with automated testing, security scanning, and deployment to Google Cloud Run.

### CI/CD Pipeline Architecture

```
Push to GitHub
    ↓
┌─────────────────────────────────────────┐
│  GitHub Actions Workflows               │
├─────────────────────────────────────────┤
│ • Frontend CI (Node.js)                 │
│ • Backend Testing (Go)                  │
│ • Security Scans (SAST, Dependencies)   │
│ • Code Quality Analysis                 │
│ • PR Checks & Validation                │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│  Build & Push (Master branch)           │
├─────────────────────────────────────────┤
│ • Docker image build                    │
│ • Push to GCP Artifact Registry         │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│  Deployment (Master branch)             │
├─────────────────────────────────────────┤
│ • Deploy to Cloud Run                   │
│ • Health checks                         │
│ • Performance monitoring                │
└─────────────────────────────────────────┘
```

## Workflows

### 1. Frontend CI (`frontend-ci.yml`)
**Triggers:** Push/PR to frontend files

**Steps:**
- TypeScript type checking
- ESLint linting
- Build verification
- Bundle size analysis
- Test execution

**Artifacts:** `.next/` build folder

### 2. Backend Testing (Built into `deploy.yml`)
**Triggers:** Push/PR to backend files

**Steps:**
- Go linter (golangci-lint)
- Unit tests with coverage
- Coverage threshold validation (50%)
- Docker image build

### 3. Security Scanning (`security.yml`)
**Triggers:** Daily (2 AM UTC), push, PR

**Scans:**
- NPM/PNPM dependency audit
- Go security (gosec)
- Container scanning (Trivy)
- CodeQL SAST analysis

### 4. Code Quality Analysis (`code-quality.yml`)
**Triggers:** Push/PR

**Checks:**
- Go linting
- TypeScript type checking
- ESLint rules
- Code duplication detection
- Cyclomatic complexity analysis

### 5. PR Checks (`pr-checks.yml`)
**Triggers:** Pull requests

**Validations:**
- Commit message format
- File changes summary
- PR description check
- Version bump detection

### 6. Release Management (`release.yml`)
**Triggers:** Version bump in `package.json`, manual dispatch

**Actions:**
- Automatic changelog generation
- GitHub Release creation
- Version tagging

### 7. Post-Deployment Monitoring (`monitoring.yml`)
**Triggers:** Deployment success, hourly schedule

**Monitors:**
- API health checks
- Performance baselines
- Uptime tracking
- Metrics collection
- Log aggregation

## Environment Variables & Secrets

### Required GitHub Secrets

For GCP deployment, configure these secrets in your repository:

```
GCP_WORKLOAD_IDENTITY_PROVIDER
GCP_SERVICE_ACCOUNT
PRODUCTION_API_URL (optional, for health checks)
CODECOV_TOKEN (optional, for coverage reports)
```

### Required Repository Variables

Configure these in repository settings:

```
GCP_PROJECT_ID: stable-glass-256211
REGISTRY: us-central1-docker.pkg.dev
SERVICE_NAME: relate-irl-api
REGION: us-central1
```

## Deployment Process

### To GCP Cloud Run (Master branch only)

1. **Automatic Deployment**
   - Push code to `master`
   - All tests must pass
   - Docker image is built and pushed
   - Service is automatically deployed to Cloud Run
   - Health checks verify deployment success

2. **Manual Deployment**
   - Trigger `deploy.yml` workflow manually
   - Select version and environment

### Health Checks

The deployment includes automatic health checks:
- Waits 10 seconds for service startup
- Attempts connection to `/health` endpoint
- 5 retry attempts with 5-second intervals
- Deployment fails if health check doesn't pass

## Rollback Procedure

### Automatic Rollback
If health checks fail after deployment, manually redeploy the previous version:

```bash
gcloud run deploy relate-irl-api \
  --image=us-central1-docker.pkg.dev/stable-glass-256211/relate-irl-api:master-<PREVIOUS_SHA>
```

### Manual Rollback
Use the GCP Console to revert to a previous Cloud Run revision.

## Monitoring & Logging

### Application Logs
Logs are automatically collected by Cloud Run:
- Access via: GCP Console → Cloud Run → Service → Logs

### Health Checks
- Endpoint: `<SERVICE_URL>/health`
- Frequency: Tested on every deployment
- Expected Response: `200 OK`

### Performance Metrics
- Memory usage monitoring
- CPU usage tracking
- Response time analysis
- Error rate tracking

## Best Practices

### Branch Protection Rules

For `master` and `develop` branches, enforce:

1. **Require status checks to pass before merging:**
   - `Frontend CI`
   - `Backend Testing`
   - `Code Quality`
   - `Security Scanning`

2. **Require code reviews:** Minimum 1 approval

3. **Dismiss stale PR approvals:** On new commits

4. **Require up-to-date branches:** Before merging

5. **Require conversation resolution:** Before merging

### Commit Message Format

Follow Conventional Commits:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`

**Example:**
```
feat(auth): add GitHub OAuth support

- Integrated GitHub OAuth provider
- Added user mapping for GitHub profiles
- Updated authentication flow

Closes #123
```

### Code Review Checklist

Before merging, verify:

- ✅ All CI checks pass
- ✅ Security scans show no critical issues
- ✅ Code coverage maintained or improved
- ✅ No console.log statements in production code
- ✅ Tests added for new features
- ✅ Documentation updated
- ✅ No breaking changes without discussion

## Troubleshooting

### Failed Frontend CI
```bash
# Run locally to debug
pnpm lint
pnpm tsc --noEmit
pnpm build
```

### Failed Backend Tests
```bash
# Run tests locally
go test -v ./...
go test -race -coverprofile=coverage.out ./...
```

### Failed Deployment
1. Check logs in GitHub Actions
2. Verify GCP credentials and permissions
3. Check health endpoint is responding
4. Review Cloud Run revision logs

### Security Scan Failures
1. Review SARIF reports in GitHub Security tab
2. Run `gosec` locally: `gosec ./...`
3. Check npm audit: `pnpm audit`

## Emergency Procedures

### Critical Bug in Production

1. **Immediate Action:**
   - Revert to previous stable version using Cloud Run console
   - Notify team

2. **Investigation:**
   - Review logs and error reports
   - Create hotfix branch from `master`

3. **Fix & Deploy:**
   - Create PR for hotfix
   - Merge after quick review
   - Deploy via `master` push (automatic)

### Security Incident

1. **Immediate Actions:**
   - Disable the affected service if needed
   - Notify security team
   - Check if secrets were exposed

2. **Remediation:**
   - Create patch in dedicated branch
   - Run security scans locally
   - Deploy after verification

3. **Post-Incident:**
   - Document incident and resolution
   - Update security policies if needed
   - Add regression tests

## Contacts & Escalation

- **DevOps Lead:** [Your contact]
- **Security Team:** [Your contact]
- **Deployment Issues:** Create issue in repository

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Google Cloud Run Documentation](https://cloud.google.com/run/docs)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Security Best Practices](./SECURITY.md)
