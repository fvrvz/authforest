---
title: Docker Setup
---

# Docker Setup

The fastest way to run AuthForest is with **Docker Compose**. The monorepo includes a `docker-compose.yml` that orchestrates the Go backend, SvelteKit frontend, and PostgreSQL database.

## Prerequisites

- <a href="https://docs.docker.com/get-docker/" target="_blank" rel="noopener noreferrer">Docker</a> and Docker Compose installed
- <a href="https://git-scm.com/" target="_blank" rel="noopener noreferrer">Git</a>

## Quick Start

```bash
git clone https://github.com/fvrvz/authforest.git
cd authforest
```

### Full Stack

Start everything — backend, frontend, and database:

```bash
docker compose up --build
```

| Service        | URL                                                      |
| -------------- | -------------------------------------------------------- |
| Frontend       | `http://localhost:5173`                                  |
| Backend API    | `http://localhost:8080`                                  |
| OIDC Discovery | `http://localhost:8080/.well-known/openid-configuration` |

### Development Mode (Recommended)

For local development with hot reload, start only the database via Docker and run the services locally:

```bash
# Start PostgreSQL
docker compose up db
```

Then in separate terminals:

```bash
# Terminal 1 — Backend
cd backend
cp .env-example .env   # edit with your values
go run .

# Terminal 2 — Frontend
cd frontend
cp .env-example .env   # edit with your values
pnpm install
pnpm dev
```

### Mix and Match

```bash
docker compose up db backend    # DB + Backend (run frontend locally)
docker compose up db             # DB only (run both locally)
docker compose up --build        # Everything
```

## Configuration

Environment variables are set in `docker-compose.yml` for the backend service:

| Variable      | Description                        | Default      |
| ------------- | ---------------------------------- | ------------ |
| `DB_HOST`     | PostgreSQL host                    | `db`         |
| `DB_PORT`     | PostgreSQL port                    | `5432`       |
| `DB_USER`     | Database user                      | —            |
| `DB_PASSWORD` | Database password                  | —            |
| `DB_NAME`     | Database name                      | `authforest` |
| `JWT_SECRET`  | Secret for HS256 API tokens        | —            |
| `OIDC_ISSUER` | Your externally reachable base URL | —            |

For local development without Docker, copy `.env-example` to `.env` in the `backend/` directory and set your values there.

## RSA Key Persistence

AuthForest generates an RSA key pair on first startup for signing OIDC tokens (RS256). In Docker, the key is stored in the `app_data` volume mounted at `/data`:

```yaml
volumes:
  - app_data:/data
```

> If you lose the RSA key, all previously issued OIDC tokens will become invalid.

## Verify It's Running

```bash
curl http://localhost:8080/.well-known/openid-configuration
```

You should see the OIDC discovery document with all supported endpoints.
