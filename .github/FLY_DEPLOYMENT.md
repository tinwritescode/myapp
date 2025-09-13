# Fly.io GitHub Actions CI/CD Setup

This document explains how to set up and use the GitHub Actions CI/CD pipeline for deploying to Fly.io.

## Required GitHub Secrets

Before the deployment can work, you need to set up the following secret in your GitHub repository:

### FLY_API_TOKEN

1. **Generate the token:**
   ```bash
   flyctl auth token
   ```

2. **Add to GitHub Secrets:**
   - Go to your GitHub repository
   - Navigate to Settings → Secrets and variables → Actions
   - Click "New repository secret"
   - Name: `FLY_API_TOKEN`
   - Value: Paste the token from step 1

## How It Works

### On Pull Requests
- **Backend**: Runs Go tests and builds the application
- **Frontend**: Runs linting and builds the React application
- **No deployment**: Only testing and validation

### On Push to Main Branch
- **Backend**: Deploys to `myapp-1757744589` using the root `fly.toml`
- **Frontend**: Deploys to `myapp-frontend` using `frontend/fly.toml`
- **Parallel deployment**: Both services deploy simultaneously

## Applications

- **Backend**: `myapp-1757744589` (Go API server)
- **Frontend**: `myapp-frontend` (React SPA)

## Manual Deployment

If you need to deploy manually:

```bash
# Deploy backend
flyctl deploy --config fly.toml

# Deploy frontend
flyctl deploy --config fly.toml --config-dir frontend
```

## Monitoring

- Check deployment status in the GitHub Actions tab
- Monitor application health in the Fly.io dashboard
- Backend health check: `/api/v1/ping`
- Frontend health check: `/`

## Troubleshooting

1. **Authentication errors**: Verify `FLY_API_TOKEN` secret is set correctly
2. **Build failures**: Check the Actions logs for specific error messages
3. **Deployment failures**: Ensure both `fly.toml` files are properly configured
4. **Health check failures**: Verify the health check endpoints are accessible
