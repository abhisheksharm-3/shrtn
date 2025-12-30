# shrtn - URL Shortener

A modern URL shortening service with QR code generation and link previews.

## Features

- 🔗 **URL Shortening** - Transform long URLs into clean, memorable short links
- 📱 **QR Codes** - Auto-generated, branded QR codes for every link
- 🖼️ **Link Previews** - Open Graph metadata fetching for destination URLs
- 🌓 **Dark/Light Mode** - Theme switcher with system preference detection
- 🔒 **Security** - Rate limiting, API key auth, SSRF protection

## Quick Start

### Option 1: Pull from DockerHub (Recommended)

```bash
# Create .env file with your credentials
cat > .env << EOF
APPWRITE_PROJECT_ID=your-project-id
APPWRITE_API_KEY=your-api-key
APPWRITE_DATABASE_ID=your-database-id
APPWRITE_COLLECTION_ID=your-collection-id
EOF

# Run the container
docker run -d -p 80:80 --env-file .env abhisheksan/shrtn:latest
```

### Option 2: Using docker-compose.deploy.yml

```bash
# Copy deploy compose file to your server
curl -O https://raw.githubusercontent.com/abhisheksharm-3/shrtn/main/docker-compose.deploy.yml

# Create .env and run
docker compose -f docker-compose.deploy.yml up -d
```

### Option 3: Build Locally

```bash
git clone https://github.com/abhisheksharm-3/shrtn.git
cd shrtn
docker build -t shrtn .
docker run -p 80:80 --env-file .env shrtn
```

## CI/CD Pipeline

The project uses GitHub Actions to automatically build and push Docker images:

1. **On push to `main`**: Builds and pushes with `latest` tag
2. **On tag `v*`**: Builds and pushes with version tag (e.g., `v1.0.0`)

### Setting Up CI/CD

Add these secrets to your GitHub repository:
- `DOCKERHUB_USERNAME`: Your DockerHub username
- `DOCKERHUB_TOKEN`: DockerHub access token (create at hub.docker.com → Account Settings → Security)

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Docker Container                     │
│  ┌─────────────────────────────────────────────────┐    │
│  │              Nginx (Port 80)                    │    │
│  │  • Serves frontend static files                 │    │
│  │  • Proxies /api/* to backend                    │    │
│  │  • Handles short URL redirects                  │    │
│  └─────────────────────────────────────────────────┘    │
│                         │                               │
│                         ▼                               │
│  ┌─────────────────────────────────────────────────┐    │
│  │           Go Backend (Port 8080)                │    │
│  │  • API endpoints                                │    │
│  │  • NOT exposed externally                       │    │
│  └─────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `APPWRITE_PROJECT_ID` | Appwrite project ID | Yes |
| `APPWRITE_API_KEY` | Appwrite API key | Yes |
| `APPWRITE_DATABASE_ID` | Appwrite database ID | Yes |
| `APPWRITE_COLLECTION_ID` | Appwrite collection ID | Yes |
| `API_KEY` | API key for authenticated endpoints | No |
| `RATE_LIMIT_PER_MINUTE` | Rate limit (default: 60) | No |
| `RATE_LIMIT_BURST` | Burst limit (default: 10) | No |

## Tech Stack

**Frontend**: React 19, TypeScript, TanStack Query, Tailwind CSS  
**Backend**: Go, Gin, Appwrite  
**Infrastructure**: Docker, Nginx
