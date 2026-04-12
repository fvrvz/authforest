# GitHub Workflows Documentation

This document describes the improved GitHub Actions workflows for the AuthForest project.

## Workflows Overview

### 1. CI Workflow (`ci.yml`)

**Triggers:** Push and PR to `develop` and `main` branches

**Features:**

- **Testing:** Runs unit tests with PostgreSQL service
- **Linting:** Uses golangci-lint with comprehensive rules
- **Build Testing:** Multi-platform Docker build test
- **Coverage:** Reports to Codecov

### 2. Docker Build & Push (`docker-build.yml`)

**Triggers:** Push to `develop`/`main`, tags `v*`, PRs, manual dispatch

**Features:**

- **Multi-platform builds:** linux/amd64, linux/arm64
- **Smart tagging strategy:**
  - Branch names for development
  - Semantic versioning for releases
  - SHA-based tags for traceability
  - `latest` tag for default branch
- **Security:** Build attestations, SBOM generation, Trivy scanning
- **Caching:** GitHub Actions cache for faster builds
- **Conditional push:** Only pushes on non-PR events

**Image Tags Generated:**

- `develop` → `ghcr.io/fvrvz/authforest:develop`
- `v1.2.3` → `ghcr.io/fvrvz/authforest:v1.2.3`, `1.2`, `1`, `latest`
- Commit → `ghcr.io/fvrvz/authforest:develop-sha123abc`

### 3. Release Workflow (`release.yml`)

**Triggers:** Git tags `v*`

**Features:**

- **Automatic changelog generation**
- **GitHub release creation**
- **Links to Docker images**
- **Prerelease detection** (for tags with `-` like `v1.0.0-beta`)

### 4. Security & Dependencies (`security.yml`)

**Triggers:** Weekly schedule, manual dispatch, push to main branches

**Features:**

- **Security scanning:** gosec, Trivy
- **Dependency review:** On PRs
- **Automated dependency updates:** Creates PRs weekly
- **SARIF uploads:** Security issues in GitHub Security tab

## Tagging Strategy

### Semantic Versioning

Use semantic versioning tags for releases:

- `v1.0.0` - Major release
- `v1.1.0` - Minor release
- `v1.1.1` - Patch release
- `v1.0.0-beta.1` - Prerelease

### Creating a Release

1. Create and push a tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
2. The release workflow will automatically:
   - Create a GitHub release with changelog
   - Docker images will be built and tagged appropriately

## Docker Images

Images are available at: `ghcr.io/fvrvz/authforest`

### Available Tags

- `latest` - Latest stable release
- `v1.2.3` - Specific version
- `1.2` - Major.minor version
- `1` - Major version
- `develop` - Latest development build
- `develop-sha123abc` - Specific commit

### Pulling Images

```bash
# Latest stable
docker pull ghcr.io/fvrvz/authforest:latest

# Specific version
docker pull ghcr.io/fvrvz/authforest:v1.0.0

# Development
docker pull ghcr.io/fvrvz/authforest:develop
```

## Security Features

1. **Vulnerability Scanning:** Trivy scans for OS and dependency vulnerabilities
2. **Code Security:** gosec analyzes Go code for security issues
3. **Supply Chain Security:** Build attestations and SBOM generation
4. **Dependency Monitoring:** Automated dependency review and updates

## Configuration Files

- `.golangci.yml` - Linter configuration with comprehensive rules
- `.dockerignore` - Optimizes Docker build context
- Updated `Dockerfile` - Multi-stage build with security best practices

## Best Practices Implemented

1. **Security First:** Non-root user, minimal base image, security scanning
2. **Multi-platform:** Supports both amd64 and arm64 architectures
3. **Efficient Caching:** GitHub Actions cache and Docker layer caching
4. **Comprehensive Testing:** Unit tests, linting, security checks
5. **Automated Releases:** Semantic versioning with automatic changelog
6. **Supply Chain Security:** Attestations and vulnerability monitoring
