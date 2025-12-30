# shrtn - Backend

Go backend service for the URL shortener.

## Tech Stack

- **Framework**: Gin
- **Database**: Appwrite
- **Language**: Go 1.24+

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/shorten` | Create shortened URL |
| `GET` | `/api/:shortCode` | Get URL info |
| `GET` | `/api/urls` | List all URLs (paginated) |
| `DELETE` | `/api/:shortCode` | Delete URL |
| `GET` | `/api/preview?url=` | Fetch link metadata |
| `GET` | `/:shortCode` | Redirect to original URL |
| `GET` | `/health` | Health check |

## Local Development

```bash
# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Run server
go run cmd/server/main.go
```

## Environment Variables

See `.env.example` for all required variables.

## Directory Structure

```
server/
├── cmd/server/main.go      # Entry point
├── internal/
│   ├── api/                # HTTP handlers & router
│   ├── config/             # Configuration
│   ├── middleware/         # Rate limiting, auth, security
│   ├── model/              # Data models
│   ├── repository/         # Database layer
│   └── service/            # Business logic
├── Dockerfile
├── go.mod
└── .env.example
```

## Docker

```bash
docker build -t shrtn-backend .
docker run -p 8080:8080 --env-file .env shrtn-backend
```
