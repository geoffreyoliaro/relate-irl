# GitHub Configuration & CI/CD Setup

Welcome to the comprehensive GitHub Actions CI/CD pipeline for this project. This directory contains all the configurations, workflows, and documentation for automated testing, security scanning, and deployment.

## 📋 Quick Navigation

### Documentation
- **[CI/CD Summary](CI_CD_SUMMARY.md)** ← Start here for overview
- **[DevOps Guide](DEVOPS.md)** ← Complete operational guide
- **[Security Policy](SECURITY.md)** ← Security practices
- **[GitHub Setup](GITHUB_SETUP.md)** ← Step-by-step setup instructions

### Configuration Files
- **[CODEOWNERS](CODEOWNERS)** - Code ownership rules
- **[Pull Request Template](pull_request_template.md)** - PR format guide

### Issue & PR Templates
- **[Bug Report Template](ISSUE_TEMPLATE/bug_report.md)**
- **[Feature Request Template](ISSUE_TEMPLATE/feature_request.md)**

## 🚀 CI/CD Workflows

### Automated Workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| **Frontend CI** | PR/Push to frontend files | Test, lint, build Next.js |
| **Backend Tests** | PR/Push to backend files | Test, lint Go code, coverage |
| **Security Scanning** | Daily + all events | SAST, dependency, container scanning |
| **Code Quality** | All push/PR | Complexity, duplication, maintainability |
| **PR Checks** | Pull requests | Commit format, version bumps, descriptions |
| **Build & Deploy** | Master push | Docker build, push, Cloud Run deploy |
| **Release** | Version bump | Changelog generation, GitHub release |
| **Monitoring** | Post-deploy + hourly | Health checks, metrics, uptime |

### Workflow Status

To see current workflow status:
1. Navigate to **Actions** tab in GitHub
2. Select workflow to view recent runs
3. Click run to see logs and details

## ✅ Implementation Checklist

Before using this pipeline:

- [ ] Read [CI/CD Summary](CI_CD_SUMMARY.md)
- [ ] Follow [GitHub Setup](GITHUB_SETUP.md) guide
- [ ] Configure branch protection rules
- [ ] Add required GitHub Secrets
- [ ] Set up code scanning (CodeQL)
- [ ] Enable Dependabot
- [ ] Create test PR to verify workflows
- [ ] Share documentation with team

## 🔧 Configuration

### Required GitHub Secrets

```yaml
GCP_WORKLOAD_IDENTITY_PROVIDER  # GCP workload identity
GCP_SERVICE_ACCOUNT             # GCP service account email
```

### Optional GitHub Secrets

```yaml
CODECOV_TOKEN               # For coverage reporting
PRODUCTION_API_URL         # For health checks
```

### Repository Variables

```yaml
GCP_PROJECT_ID     = stable-glass-256211
REGISTRY           = us-central1-docker.pkg.dev
SERVICE_NAME       = relate-irl-api
REGION             = us-central1
```

## 📊 Workflow Files

```
workflows/
├── deploy.yml              # Backend test, build, deploy (existing)
├── frontend-ci.yml         # Frontend test, lint, build
├── security.yml            # SAST, dependencies, containers
├── pr-checks.yml          # PR validation and checks
├── code-quality.yml       # Code complexity and quality
├── release.yml            # Automated releases
└── monitoring.yml         # Post-deploy monitoring
```

### Workflow Timing

- **Frontend CI:** ~5-10 minutes
- **Backend Tests:** ~5-10 minutes
- **Security Scans:** ~10-15 minutes
- **Full Pipeline:** ~15-20 minutes (parallel execution)
- **Deployment:** ~5-10 minutes after all checks pass

## 🔐 Security Features

✅ **Code Scanning**
- CodeQL SAST analysis
- Go security (gosec)
- JavaScript/TypeScript linting

✅ **Dependency Management**
- NPM/PNPM audit
- Go vulnerability checks
- Dependabot alerts

✅ **Container Security**
- Trivy image scanning
- Layer vulnerability analysis

✅ **Secret Management**
- GitHub Secrets for sensitive data
- No credentials in code
- Audit logging

## 📈 Monitoring & Alerts

### GitHub Notifications
- Failed workflow runs
- PR check status
- Security alerts
- Deployment notifications

### Workflow Insights
- Run duration trends
- Success/failure rates
- Performance metrics

### Access via
1. **GitHub Actions tab** - View all workflows
2. **Security tab** - View code scanning and Dependabot alerts
3. **Environments** - View deployment history
4. **Email notifications** - Workflow alerts

## 🐛 Troubleshooting

### Workflows Not Running
1. Check workflow file syntax (YAML format)
2. Verify branch protection rules allow actions
3. Review workflow logs in GitHub Actions tab
4. Check that paths match for triggering workflows

### Failed Tests
- Run tests locally: `pnpm test`, `go test ./...`
- Check environment variables
- Review workflow logs for specific errors
- Check if dependencies need updating

### Deployment Issues
- Verify GCP credentials in secrets
- Check Cloud Run quotas and limits
- Review application logs in GCP Console
- Ensure health endpoint is responding

### Security Scan Alerts
- Review CodeQL alerts in Security tab
- Update vulnerable dependencies
- Address linting issues
- Check container image vulnerabilities

See [DevOps Guide](DEVOPS.md#troubleshooting) for more details.

## 👥 Team Guidelines

### Commit Message Format

Follow Conventional Commits:
```
type(scope): subject

body

footer
```

**Types:** feat, fix, docs, style, refactor, perf, test, chore

**Example:**
```
feat(api): add user authentication

- Added JWT token generation
- Implemented token validation middleware

Closes #123
```

### PR Process

1. Create feature branch
2. Make changes following code standards
3. Push branch and create PR
4. Wait for all workflows to pass
5. Request review (CODEOWNERS auto-assigned)
6. Address review comments
7. Merge when approved
8. Automatic deployment after merge to master

### Code Review

- Review code quality and security
- Verify tests pass
- Check for performance implications
- Ensure documentation is updated
- Look for breaking changes

## 📚 Documentation

### For New Team Members
1. Start with [CI/CD Summary](CI_CD_SUMMARY.md)
2. Read [GitHub Setup](GITHUB_SETUP.md) for configuration
3. Review [Security Policy](SECURITY.md)
4. Study [DevOps Guide](DEVOPS.md) for operations

### For Operators
- [DevOps Guide](DEVOPS.md) - Complete reference
- [Security Policy](SECURITY.md) - Security procedures
- Workflow comments - Technical implementation details

### For Developers
- [Pull Request Template](pull_request_template.md) - PR format
- [CODEOWNERS](CODEOWNERS) - Who reviews what
- Workflow status - Check GitHub Actions tab

## 🎯 Best Practices

✅ **Do**
- Write clear commit messages
- Add meaningful PR descriptions
- Run tests locally before pushing
- Keep PRs focused on single feature
- Respond to code review comments
- Use GitHub Issues for planning

❌ **Don't**
- Commit secrets or credentials
- Force-push to shared branches
- Skip failing tests
- Ignore security scan alerts
- Merge without reviews (when required)
- Use overly generic commit messages

## 📞 Support

### Issues
- GitHub Issues - Report bugs or request features
- Use provided issue templates
- Include reproduction steps and environment info

### Documentation Questions
- Check [DevOps Guide](DEVOPS.md)
- Review workflow comments in `.yml` files
- Search existing GitHub Issues

### Infrastructure Issues
- GCP Console logs - Application errors
- GitHub Actions logs - Build/test errors
- Cloud Run logs - Deployment issues

## 🔄 Maintenance

### Weekly
- [ ] Review failed workflow runs
- [ ] Check Dependabot alerts
- [ ] Verify deployments succeeded

### Monthly
- [ ] Update workflow dependencies
- [ ] Review security scan results
- [ ] Update documentation as needed

### Quarterly
- [ ] Full security audit
- [ ] Review workflow performance
- [ ] Update best practices documentation

## 📞 Contact

- **DevOps Questions:** Review DEVOPS.md or create issue
- **Security Questions:** See SECURITY.md
- **Setup Help:** Follow GITHUB_SETUP.md

## 🎉 Success Indicators

After setup, you'll see:
- ✅ All workflows running on PRs
- ✅ Status checks required for merge
- ✅ Code reviews auto-assigned
- ✅ Security alerts in dashboard
- ✅ Automatic deployments on master
- ✅ Release notes auto-generated

## 📞 Quick Links

- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Google Cloud Run](https://cloud.google.com/run)
- [CodeQL Analysis](https://codeql.github.com/)

---

**Last Updated:** June 23, 2026  
**Version:** 1.0.0  
**Status:** ✅ Production Ready

For comprehensive documentation, start with [CI/CD Summary](CI_CD_SUMMARY.md).
