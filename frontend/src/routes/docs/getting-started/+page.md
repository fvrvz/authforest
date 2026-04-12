---
title: Getting Started
---

# Getting Started

This guide walks you through setting up AuthForest and integrating it into your first application.

## What is AuthForest?

AuthForest is a self-hosted **OpenID Connect (OIDC) Identity Provider** built with Go. It gives you:

- A standards-compliant OIDC/OAuth 2.0 server
- A built-in admin panel to manage apps, users, and roles
- Docker-ready deployment with PostgreSQL

## Setup Flow

1. **Deploy** AuthForest using [Docker Compose](/docs/docker-setup)
2. **Register an application** — create an OAuth client in the admin panel
3. **Create users** — via self-registration or from the admin panel
4. **Integrate** — use any OIDC client library to connect your app

## 1. Deploy with Docker

See the [Docker Setup](/docs/docker-setup) guide for the full walkthrough.

Quick version:

```bash
git clone https://github.com/fvrvz/authforest.git
cd authforest
```

**Full stack** (backend + frontend + database):

```bash
docker compose up --build
```

**Database only** (recommended for development):

```bash
docker compose up db
```

Then run the backend and frontend locally:

```bash
# Terminal 1 — Backend
cd backend
go run .

# Terminal 2 — Frontend
cd frontend
pnpm dev
```

Backend will be available at `http://localhost:8080` and the frontend at `http://localhost:5173`.

## 2. Access the Admin Panel

Navigate to `http://localhost:5173` and sign in with the credentials you configured.

From the dashboard you can:

- Register OAuth applications
- Manage users (create, edit, assign roles)
- Create and manage roles
- Configure token expiry per application

## 3. Register Your First App

Go to **Applications → New Application** and fill in:

| Field                | Example Value                    |
| -------------------- | -------------------------------- |
| **Application Name** | My Web App                       |
| **Redirect URIs**    | `http://localhost:3000/callback` |
| **Allowed Origins**  | `http://localhost:3000`          |

After saving, you'll receive a **Client ID** and **Client Secret**. See [Register an Application](/docs/register-app) for more details.

## 4. Integrate with Your App

Use any OIDC client library. AuthForest supports Authorization Code Flow with PKCE.

Your OIDC client needs:

| Setting                | Value                          |
| ---------------------- | ------------------------------ |
| **Authority / Issuer** | `http://localhost:8080`        |
| **Client ID**          | From step 3                    |
| **Redirect URI**       | Must match what you registered |
| **Scopes**             | `openid profile email`         |

See the [OIDC Integration](/docs/oidc-integration) section for framework-specific examples (JavaScript, React, Python, Go, etc.).

## OIDC Endpoints

AuthForest exposes standard OIDC endpoints, all discoverable via:

```
GET /.well-known/openid-configuration
```

| Endpoint      | Path                     |
| ------------- | ------------------------ |
| Authorization | `/oauth2/authorize`      |
| Token         | `/oauth2/token`          |
| UserInfo      | `/oauth2/userinfo`       |
| JWKS          | `/.well-known/jwks.json` |
