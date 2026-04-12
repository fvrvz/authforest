---
title: Docker Setup
---

# Docker Setup

The fastest way to run AuthForest is with **Docker Compose**. This starts both the Go backend and PostgreSQL database.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose installed

## Quick Start

Create a `docker-compose.yml`:

```yaml
services:
  authforest:
    image: authforest/authforest:latest
    ports:
      - '8080:8080'
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=authforest
      - DB_PASSWORD=authforest
      - DB_NAME=authforest
      - JWT_SECRET=change-me-to-a-secure-random-string
      - ADMIN_EMAIL=admin@example.com
      - ADMIN_PASSWORD=Admin@123
    volumes:
      - authforest-data:/data
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=authforest
      - POSTGRES_PASSWORD=authforest
      - POSTGRES_DB=authforest
    volumes:
      - pg-data:/var/lib/postgresql/data
    healthcheck:
      test: ['CMD-SHELL', 'pg_isready -U authforest']
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  authforest-data:
  pg-data:
```

Then run:

```bash
docker compose up -d
```

AuthForest will be available at `http://localhost:8080`.

## Configuration

| Variable         | Description                 | Default      |
| ---------------- | --------------------------- | ------------ |
| `DB_HOST`        | PostgreSQL host             | `localhost`  |
| `DB_PORT`        | PostgreSQL port             | `5432`       |
| `DB_USER`        | Database user               | `authforest` |
| `DB_PASSWORD`    | Database password           | —            |
| `DB_NAME`        | Database name               | `authforest` |
| `JWT_SECRET`     | Secret for HS256 API tokens | —            |
| `ADMIN_EMAIL`    | Initial admin user email    | —            |
| `ADMIN_PASSWORD` | Initial admin user password | —            |
| `SERVER_PORT`    | Port to listen on           | `8080`       |

## RSA Key Persistence

AuthForest generates an RSA key pair on first startup for signing OIDC tokens (RS256). The key is stored in the `/data` directory inside the container. Mount this as a volume to persist keys across restarts:

```yaml
volumes:
  - authforest-data:/data
```

> If you lose the RSA key, all previously issued OIDC tokens will become invalid.

## Verify It's Running

```bash
curl http://localhost:8080/.well-known/openid-configuration
```

You should see the OIDC discovery document with all supported endpoints.
