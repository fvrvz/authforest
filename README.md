<p align="center">
  <img src="frontend/static/favicon.svg" width="64" height="64" alt="AuthForest" />
</p>

<h1 align="center">AuthForest</h1>

<p align="center">
  A comprehensive authentication and identity platform — Go backend with OpenID Connect support + modern SvelteKit frontend, managed as a pnpm monorepo.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/SvelteKit-2.x-FF3E00?logo=svelte&logoColor=white" alt="SvelteKit" />
  <img src="https://img.shields.io/badge/PostgreSQL-17-4169E1?logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/License-MIT-green" alt="License" />
</p>

---

## Project Structure

```
authforest/
├── backend/                 # Go authentication service
│   ├── config/              # Configuration management
│   ├── controllers/         # HTTP route handlers
│   ├── services/            # Business logic (OIDC, auth, users)
│   ├── models/              # Database models (GORM)
│   ├── helpers/             # Utilities (JWT, RSA, OIDC)
│   ├── middlewares/         # HTTP middlewares
│   ├── dto/                 # Data Transfer Objects
│   ├── db/                  # Database setup & seeding
│   ├── templates/           # HTML templates (OIDC login flow)
│   ├── Dockerfile           # Multi-stage production image
│   └── main.go              # Entry point
│
├── frontend/                # SvelteKit Identity Provider UI
│   ├── src/                 # Components, routes, state, services
│   ├── static/              # Static assets
│   ├── Dockerfile           # Multi-stage (build → nginx)
│   └── svelte.config.js     # SvelteKit config
│
├── docker-compose.yml       # Full-stack orchestration
├── go.work                  # Go workspace (gopls support)
├── package.json             # Monorepo root config
├── pnpm-workspace.yaml      # pnpm workspace definition
└── .github/workflows/       # CI/CD pipelines
```

## Quick Start

### Prerequisites

- [Go 1.24+](https://go.dev/dl/) (for backend)
- [Node.js 20+](https://nodejs.org/) & [pnpm](https://pnpm.io/) (for frontend)
- [Docker](https://www.docker.com/) (for database / full-stack)

### 1. Clone & Install

```bash
git clone https://github.com/fvrvz/authforest.git
cd authforest
pnpm install
```

### 2. Set Up Environment

```bash
# Backend
cp backend/.env-example backend/.env
# Edit backend/.env with your DB credentials and secrets

# Frontend
cp frontend/.env-example frontend/.env
# Edit frontend/.env with your API URL
```

### 3. Start Development

#### Option A: Docker for everything

```bash
docker compose up --build
```

This starts all three services — backend (`:8080`), frontend (`:5173`), and PostgreSQL.

#### Option B: Docker for DB only (recommended for development)

```bash
# Start only the database
docker compose up db

# In another terminal — run the backend with hot reload
cd backend
go run .

# In another terminal — run the frontend dev server
cd frontend
pnpm dev
```

#### Option C: Pick what you need

```bash
docker compose up db backend    # DB + Backend (run frontend locally)
docker compose up db             # DB only (run both services locally)
```

### 4. Verify

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **OIDC Discovery**: http://localhost:8080/.well-known/openid-configuration

## Tech Stack

| Layer        | Technology                                                    |
| ------------ | ------------------------------------------------------------- |
| **Backend**  | Go 1.24+, Gin, GORM, PostgreSQL                               |
| **Frontend** | SvelteKit 2.x, Svelte 5, TypeScript, Tailwind CSS 4, Flowbite |
| **Auth**     | OIDC Core 1.0, OAuth 2.0 + PKCE, RS256/HS256 JWT              |
| **Infra**    | Docker, Docker Compose, GitHub Actions, GitHub Pages          |

## Features

### OIDC / OAuth 2.0

- Full OpenID Connect Core 1.0 compliant Identity Provider
- Authorization Code Flow + PKCE (S256)
- RS256 token signing with auto-generated RSA keys
- OIDC Discovery (`/.well-known/openid-configuration`)
- JWKS endpoint (`/.well-known/jwks.json`)
- Scope-aware UserInfo endpoint (profile, email)
- Dynamic OAuth client registration
- Refresh token rotation with reuse detection

### Security

- bcrypt password hashing
- Single-use authorization codes
- Strict redirect URI validation (RFC 6749 §3.1.2)
- `Cache-Control: no-store` on token responses
- Configurable CORS

### Frontend

- Modern identity provider UI
- Dark/light theme support
- Responsive design
- Authentication flows (login, register)
- User & role management dashboard
- OIDC client management

## API Endpoints

### OIDC / OAuth 2.0

| Endpoint                            | Method   | Description             |
| ----------------------------------- | -------- | ----------------------- |
| `/.well-known/openid-configuration` | GET      | OIDC Discovery document |
| `/.well-known/jwks.json`            | GET      | JSON Web Key Set        |
| `/oauth2/authorize`                 | GET/POST | Authorization endpoint  |
| `/oauth2/token`                     | POST     | Token exchange          |
| `/oauth2/userinfo`                  | GET/POST | UserInfo claims         |

### Auth (Legacy JWT)

| Endpoint               | Method | Auth | Description       |
| ---------------------- | ------ | ---- | ----------------- |
| `/api/v1/auth/login`   | POST   | No   | Login (HS256 JWT) |
| `/api/v1/auth/logout`  | GET    | Yes  | Logout            |
| `/api/v1/auth/refresh` | POST   | Yes  | Refresh token     |

### User Management

| Endpoint                 | Method | Auth | Description |
| ------------------------ | ------ | ---- | ----------- |
| `/api/v1/users/register` | POST   | No   | Register    |
| `/api/v1/users`          | GET    | Yes  | List users  |
| `/api/v1/users/:userId`  | GET    | Yes  | Get user    |
| `/api/v1/users/:userId`  | PATCH  | Yes  | Update user |
| `/api/v1/users/:userId`  | DELETE | Yes  | Delete user |

### OAuth Client Management

| Endpoint                  | Method | Auth | Description           |
| ------------------------- | ------ | ---- | --------------------- |
| `/api/v1/oauth2/register` | POST   | Yes  | Register OAuth client |

## OIDC Integration

### Register a Client

```bash
curl -X POST http://localhost:8080/api/v1/oauth2/register \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "client_name": "My SPA",
    "client_type": "public",
    "redirect_uris": ["http://localhost:4200/callback"],
    "scopes": "openid profile email offline_access",
    "grant_types": "authorization_code"
  }'
```

### Client Library Example (oidc-client-ts)

```typescript
const config: UserManagerSettings = {
  authority: "http://localhost:8080",
  client_id: "YOUR_CLIENT_ID",
  redirect_uri: "http://localhost:3000/callback",
  response_type: "code",
  scope: "openid profile email",
};
```

### Supported Scopes

| Scope            | Claims                                                                  |
| ---------------- | ----------------------------------------------------------------------- |
| `openid`         | `sub`, `iss`, `aud`, `exp`, `iat`, `auth_time`, `nonce`                 |
| `profile`        | `name`, `given_name`, `family_name`, `preferred_username`, `updated_at` |
| `email`          | `email`, `email_verified`                                               |
| `offline_access` | Issues a refresh token                                                  |

## Configuration

The backend uses `config/files/config.yml` with environment variable substitution:

```yaml
server:
  port: 8080
  cors:
    allowOrigins: ["http://localhost:5173"]

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  username: ${DB_USER}
  password: ${DB_PASSWORD}
  db: ${DB_NAME}

jwt:
  expiry_minutes: 15
  jwt_secret: ${JWT_SECRET}

oidc:
  issuer: ${OIDC_ISSUER}
  rsa_key_path: rsa_private.pem # auto-generated if missing
```

### Environment Variables

| Variable              | Description                | Default      | Required |
| --------------------- | -------------------------- | ------------ | -------- |
| `DB_HOST`             | Database host              | `localhost`  | Yes      |
| `DB_PORT`             | Database port              | `5432`       | Yes      |
| `DB_USER`             | Database username          | —            | Yes      |
| `DB_PASSWORD`         | Database password          | —            | Yes      |
| `DB_NAME`             | Database name              | `authforest` | Yes      |
| `DB_SSLMODE`          | PostgreSQL SSL mode        | `disable`    | No       |
| `JWT_SECRET`          | HS256 JWT signing secret   | —            | Yes      |
| `OIDC_ISSUER`         | OIDC Issuer URL            | —            | Yes      |
| `VITE_BASE_URL`       | Backend API URL (frontend) | —            | Yes      |
| `VITE_OIDC_CLIENT_ID` | OIDC client ID (frontend)  | —            | Yes      |

## Scripts

```bash
# From monorepo root
pnpm dev              # Start frontend dev server
pnpm build            # Build frontend for production
pnpm test             # Run frontend tests
pnpm lint             # Lint frontend

# Backend (from backend/)
go run .              # Run locally
go test ./...         # Run tests
go build -o authforest .  # Build binary
```

## CI/CD

GitHub Actions workflows with monorepo path filters:

| Workflow           | Trigger               | Purpose                         |
| ------------------ | --------------------- | ------------------------------- |
| `backend-ci.yml`   | `backend/**` changes  | Go lint, test, build            |
| `frontend-ci.yml`  | `frontend/**` changes | Svelte check, lint, test, build |
| `deploy-pages.yml` | Push to `develop`     | Deploy frontend to GitHub Pages |
| `docker-build.yml` | Push to `main`        | Build & push Docker images      |
| `security.yml`     | Weekly + push         | Security scanning               |
| `release.yml`      | Tag push              | Create GitHub release           |

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License — see [LICENSE](LICENSE) for details.

---

<p align="center">
  <img src="frontend/static/fayso-logo.svg" height="20" alt="FaySo" />
</p>
