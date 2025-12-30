# shrtn - Frontend

React frontend for the URL shortener.

## Tech Stack

- **Framework**: React 19
- **Build**: Vite 7
- **Styling**: Tailwind CSS 4
- **State**: TanStack Query
- **Language**: TypeScript

## Features

- URL shortening with custom codes
- QR code generation with branding
- Link preview cards (Open Graph)
- Dark/Light/System theme
- Responsive design

## Local Development

```bash
# Install dependencies
npm install

# Copy environment file
cp .env.example .env

# Run dev server
npm run dev
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `VITE_API_URL` | Backend API URL | Yes (dev) |
| `VITE_API_KEY` | API key for auth | No |

> **Note**: In Docker, `VITE_API_URL` is empty because nginx proxies `/api/*` to the backend.

## Directory Structure

```
client/
├── src/
│   ├── api/                # API client
│   ├── components/
│   │   ├── features/       # Feature components
│   │   ├── layout/         # Layout components
│   │   ├── providers/      # Context providers
│   │   └── ui/             # UI primitives
│   ├── hooks/              # Custom hooks
│   ├── pages/              # Route pages
│   ├── types/              # TypeScript types
│   └── lib/                # Utilities
├── public/
├── Dockerfile
├── nginx.conf
└── .env.example
```

## Build

```bash
npm run build
```

Output is in `dist/` directory.

## Docker

```bash
docker build -t shrtn-frontend .
docker run -p 80:80 shrtn-frontend
```
