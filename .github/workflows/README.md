# GitHub Actions Workflows

## Docker Multi-Platform Build

This workflow automatically builds and publishes Docker images for multiple platforms to both DockerHub and GitHub Container Registry (GHCR).

### Supported Platforms

- linux/amd64
- linux/arm64
- linux/arm/v7

### Registries

- **DockerHub**: `sstreichan/ztdns`
- **GHCR**: `ghcr.io/sstreichan/ztdns`

### Triggers

The workflow runs on:
- Push to `master` or `main` branches
- Push of tags matching `v*` (e.g., `v1.0.0`)
- Pull requests to `master` or `main` branches (build only, no push)
- Manual trigger via workflow_dispatch

### Image Tags

The workflow automatically generates tags based on:
- Branch name (e.g., `master`, `main`)
- Pull request number (e.g., `pr-123`)
- Semantic version tags (e.g., `v1.2.3` → `1.2.3`, `1.2`, `1`)
- `latest` tag for the default branch

### Required Secrets

To use this workflow, configure the following secrets in your repository settings:

#### DockerHub
- `DOCKER_HUB_USERNAME`: Your DockerHub username
- `DOCKER_HUB_TOKEN`: DockerHub access token (create at https://hub.docker.com/settings/security)

#### GitHub Container Registry
- `GITHUB_TOKEN`: Automatically provided by GitHub Actions (no configuration needed)

### Setup Instructions

1. **Create DockerHub Access Token**:
   - Go to https://hub.docker.com/settings/security
   - Click "New Access Token"
   - Give it a description (e.g., "GitHub Actions")
   - Copy the generated token

2. **Add Repository Secrets**:
   - Go to your repository settings
   - Navigate to "Secrets and variables" → "Actions"
   - Click "New repository secret"
   - Add `DOCKER_HUB_USERNAME` with your DockerHub username
   - Add `DOCKER_HUB_TOKEN` with the token from step 1

3. **Verify GHCR Permissions**:
   - The workflow uses `GITHUB_TOKEN` which is automatically provided
   - Ensure the workflow has `packages: write` permission (already configured)

### Usage

Once configured, the workflow will automatically run on pushes and pull requests. Images will be pushed to both registries on successful builds (except for pull requests, which only build without pushing).

To manually trigger the workflow:
1. Go to the "Actions" tab in your repository
2. Select "Docker Multi-Platform Build"
3. Click "Run workflow"

### Example Docker Pull Commands

```bash
# Pull from DockerHub
docker pull sstreichan/ztdns:latest
docker pull sstreichan/ztdns:v1.0.0

# Pull from GHCR
docker pull ghcr.io/sstreichan/ztdns:latest
docker pull ghcr.io/sstreichan/ztdns:v1.0.0
```
