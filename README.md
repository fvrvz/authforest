# Auth Service Go

A robust authentication service and **OpenID Connect Identity Provider (OIDC IDP)** built with Go and the Gin web framework. Provides both traditional JWT-based authentication and a standards-compliant OIDC Authorization Code Flow with PKCE, compatible with any OIDC client library (angular-oauth2-oidc, oidc-client-ts, next-auth, etc.).

## 🚀 Features

- **OIDC Identity Provider**: Full OpenID Connect Core 1.0 compliant IDP
- **Authorization Code Flow + PKCE**: Secure OAuth 2.0 flow with S256 code challenge (RFC 7636)
- **RS256 Token Signing**: ID tokens and OIDC access tokens signed with RSA-SHA256
- **OIDC Discovery**: Auto-discoverable configuration at `/.well-known/openid-configuration`
- **JWKS Endpoint**: Public key endpoint for token verification at `/.well-known/jwks.json`
- **UserInfo Endpoint**: Scope-aware claims (profile, email) per OIDC Core §5.4
- **Dynamic Client Registration**: Register OAuth clients via API
- **Refresh Token Rotation**: Secure refresh token reuse detection
- **JWT Authentication**: Legacy HS256 token-based auth for internal API routes
- **User Management**: User registration, login, logout, and profile management
- **Database Integration**: PostgreSQL with GORM ORM
- **CORS Support**: Configurable Cross-Origin Resource Sharing
- **Docker Support**: Containerized deployment with Docker and Docker Compose

## 🛠 Tech Stack

- **Language**: Go 1.24+
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **OIDC/OAuth2**: OpenID Connect Core 1.0, RFC 7636 (PKCE)
- **Token Signing**: RS256 (OIDC) + HS256 (legacy JWT) via golang-jwt/jwt/v5
- **Containerization**: Docker & Docker Compose
- **Logger**: Custom logger (gologger)

## 📁 Project Structure

```
auth-service-go/
├── config/                 # Configuration management
│   ├── config.go
│   └── files/
│       └── config.yml
├── constants/              # Application constants
├── controllers/            # HTTP request handlers
│   ├── auth.controller.go
│   ├── users.controller.go
│   └── oauth.controller.go  # OIDC route definitions
├── db/                     # Database connection and setup
├── dto/                    # Data Transfer Objects
│   ├── auth.dto.go
│   ├── common.dto.go
│   ├── config.dto.go
│   ├── user.dto.go
│   └── oauth.dto.go         # OIDC request/response DTOs
├── helpers/                # Utility functions
│   ├── jwt.go               # Legacy HS256 JWT helpers
│   ├── oidc_jwt.go          # RS256 OIDC token generation
│   ├── rsa.go               # RSA key management
│   └── normalize-date.go
├── middlewares/            # HTTP middlewares
├── models/                 # Database models
│   ├── auth.go
│   ├── user.go
│   └── oauth.go             # OAuthClient, AuthorizationCode
├── services/               # Business logic layer
│   ├── auth.service.go
│   ├── users.service.go
│   ├── oidc_discovery.go    # Discovery + JWKS
│   ├── oidc_authorize.go    # Authorization endpoint + login
│   ├── oidc_token.go        # Token exchange (auth code + refresh)
│   ├── oidc_userinfo.go     # UserInfo endpoint
│   ├── oidc_client.go       # Client registration
│   └── oidc_helpers.go      # PKCE, redirect helpers
├── server/                 # Server setup and routing
├── templates/              # HTML templates for OIDC login flow
│   ├── login.html
│   └── error.html
├── postman-collections/    # API testing collections
├── docker-compose.yml
├── Dockerfile
└── main.go                 # Application entry point
```

## 🐳 Docker Setup (Recommended)

The easiest way to run the application is using Docker Compose, which sets up both the auth service and PostgreSQL database.

### Prerequisites

- Docker
- Docker Compose

### Quick Start

1. **Clone the repository**

   ```bash
   git clone https://github.com/fvrvz/auth-service-go.git
   cd auth-service-go
   ```

2. **Configure environment variables**

   Edit the `docker-compose.yml` file and update the following values:

   ```yaml
   environment:
     DB_USER: your_db_user # Change from 'test'
     DB_PASSWORD: your_db_password # Change from 'test'
     JWT_SECRET: "your-secret-key" # Change to a secure random string
     OIDC_ISSUER: "http://localhost:8080" # Your externally reachable base URL
   ```

3. **Start the services**

   ```bash
   docker-compose up -d
   ```

4. **Verify the service is running**
   ```bash
   curl http://localhost:8080/health
   ```

The application will be available at `http://localhost:8080` and PostgreSQL at `localhost:5432`.

### Docker Compose Services

- **auth-service**: The main Go application (port 8080)
- **db**: PostgreSQL database (port 5432)
- **postgres_data**: Persistent volume for database data

## 🔧 Manual Setup

If you prefer to run the application without Docker:

### Prerequisites

- Go 1.24 or higher
- PostgreSQL 12+
- Git

### Installation Steps

1. **Clone the repository**

   ```bash
   git clone https://github.com/fvrvz/auth-service-go.git
   cd auth-service-go
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**

   Create a database named `auth-service-go` or update the configuration accordingly.

4. **Configure environment variables**

   Create a `.env` file in the root directory:

   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=auth-service-go
   DB_SSLMODE=disable
   DB_CHANNEL_BINDING=disable
   JWT_SECRET=your-jwt-secret-key
   OIDC_ISSUER=http://localhost:8080
   ```

   > **Note:** An RSA private key (`rsa_private.pem`) is auto-generated on first startup if it doesn't exist. To use your own key, place it at the path specified in `config.yml` under `oidc.rsa_key_path`.

5. **Update configuration (optional)**

   Modify `config/files/config.yml` if needed for custom settings like CORS origins or JWT expiry times.

6. **Build and run the application**

   ```bash
   # Build the binary
   go build -o auth-service-go .

   # Run the compiled binary
   ./auth-service-go

   # Build a Windows executable
   go build -o auth-service-go.exe

   # Or, for development, run directly
   go run .
   ```

The application will start on port 8080 (configurable in `config.yml`).

## 📚 API Endpoints

### OIDC / OAuth 2.0 Endpoints

| Endpoint                            | Method   | Description                                                    |
| ----------------------------------- | -------- | -------------------------------------------------------------- |
| `/.well-known/openid-configuration` | GET      | OIDC Discovery document                                        |
| `/.well-known/jwks.json`            | GET      | JSON Web Key Set (public signing keys)                         |
| `/oauth2/authorize`                 | GET      | Authorization endpoint (renders login page)                    |
| `/oauth2/authorize`                 | POST     | Login form submission (issues authorization code)              |
| `/oauth2/token`                     | POST     | Token endpoint (exchanges code/refresh token for tokens)       |
| `/oauth2/userinfo`                  | GET/POST | UserInfo endpoint (returns claims based on access token scope) |

### Auth Endpoints (Legacy JWT)

| Endpoint               | Method | Auth | Description            |
| ---------------------- | ------ | ---- | ---------------------- |
| `/api/v1/auth/login`   | POST   | No   | User login (HS256 JWT) |
| `/api/v1/auth/logout`  | GET    | Yes  | User logout            |
| `/api/v1/auth/refresh` | POST   | Yes  | Refresh JWT token      |

### User Management Endpoints

| Endpoint                 | Method | Auth | Description       |
| ------------------------ | ------ | ---- | ----------------- |
| `/api/v1/users/register` | POST   | No   | User registration |
| `/api/v1/users`          | GET    | Yes  | Get all users     |
| `/api/v1/users/:userId`  | GET    | Yes  | Get specific user |
| `/api/v1/users/:userId`  | PATCH  | Yes  | Update user       |
| `/api/v1/users/:userId`  | DELETE | Yes  | Delete user       |

### Protected OIDC Management Endpoints

| Endpoint                  | Method | Auth | Description                 |
| ------------------------- | ------ | ---- | --------------------------- |
| `/api/v1/oauth2/register` | POST   | Yes  | Register a new OAuth client |

### API Testing

Import the Postman collection from `postman-collections/AuthServiceGo.postman_collection.json` for easy API testing.

## 🔑 OIDC Integration Guide

### 1. Register an OAuth Client

```bash
# First, obtain a JWT token via the legacy login endpoint
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_id": "admin", "password": "password"}' | jq -r '.token')

# Register a public client (SPA)
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

Response:

```json
{
  "client_id": "generated-uuid",
  "client_name": "My SPA",
  "client_type": "public",
  "redirect_uris": ["http://localhost:4200/callback"],
  "scopes": "openid profile email offline_access",
  "grant_types": "authorization_code"
}
```

### 2. Authorization Code Flow with PKCE

**Step 1:** Redirect user to the authorization endpoint:

```
GET /oauth2/authorize?
  response_type=code&
  client_id=YOUR_CLIENT_ID&
  redirect_uri=http://localhost:4200/callback&
  scope=openid profile email&
  state=random_state_value&
  nonce=random_nonce_value&
  code_challenge=BASE64URL_SHA256_OF_VERIFIER&
  code_challenge_method=S256
```

The user sees a login page, authenticates, and is redirected back with an authorization code.

**Step 2:** Exchange the code for tokens:

```bash
curl -X POST http://localhost:8080/oauth2/token \
  -d "grant_type=authorization_code" \
  -d "code=AUTHORIZATION_CODE" \
  -d "redirect_uri=http://localhost:4200/callback" \
  -d "client_id=YOUR_CLIENT_ID" \
  -d "code_verifier=YOUR_ORIGINAL_CODE_VERIFIER"
```

Response:

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 900,
  "id_token": "eyJhbGciOiJSUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJSUzI1NiIs...",
  "scope": "openid profile email"
}
```

### 3. Using with Client Libraries

#### angular-oauth2-oidc

```typescript
export const authConfig: AuthConfig = {
  issuer: "http://localhost:8080",
  clientId: "YOUR_CLIENT_ID",
  redirectUri: window.location.origin + "/callback",
  responseType: "code",
  scope: "openid profile email",
  useSilentRefresh: true,
};
```

#### oidc-client-ts / React

```typescript
const config: UserManagerSettings = {
  authority: "http://localhost:8080",
  client_id: "YOUR_CLIENT_ID",
  redirect_uri: "http://localhost:3000/callback",
  response_type: "code",
  scope: "openid profile email",
};
```

The client library will auto-discover all endpoints via `/.well-known/openid-configuration`.

### Supported Scopes

| Scope            | Claims Returned                                                         |
| ---------------- | ----------------------------------------------------------------------- |
| `openid`         | `sub`, `iss`, `aud`, `exp`, `iat`, `auth_time`, `nonce`                 |
| `profile`        | `name`, `given_name`, `family_name`, `preferred_username`, `updated_at` |
| `email`          | `email`, `email_verified`                                               |
| `offline_access` | Issues a refresh token                                                  |

## ⚙️ Configuration

The application uses a YAML configuration file (`config/files/config.yml`) with environment variable substitution:

```yaml
server:
  port: 8080
  cors:
    allowOrigins: ["http://localhost:5173"]
    allowCredentials: true

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  username: ${DB_USER}
  password: ${DB_PASSWORD}
  db: ${DB_NAME}

jwt:
  expiry_minutes: 15
  jwt_secret: ${JWT_SECRET}
  refresh_token_expiry_hours: 2

oidc:
  issuer: ${OIDC_ISSUER} # Base URL of the IDP (e.g. http://localhost:8080)
  rsa_key_path: rsa_private.pem # Path to RSA private key (auto-generated if missing)
```

## 🔐 Security Features

- **RS256 Token Signing**: OIDC tokens signed with RSA-SHA256 (2048-bit key, auto-generated)
- **PKCE (S256)**: Proof Key for Code Exchange prevents authorization code interception attacks
- **Authorization Code Replay Detection**: Codes are single-use, marked as used on exchange
- **Refresh Token Rotation**: Old refresh tokens are revoked when a new one is issued
- **Scope-Gated Claims**: UserInfo and ID token claims filtered by granted scopes
- **Redirect URI Validation**: Strict exact-match comparison per RFC 6749 §3.1.2
- **Cache-Control: no-store**: All token endpoint responses prevent caching (OIDC Core §3.1.3.3)
- **Password Hashing**: bcrypt-based secure password storage
- **CORS Protection**: Configurable cross-origin request handling
- **Input Validation**: Request validation using Gin's built-in validators

## 🚦 Health Check

The service includes health check endpoints for monitoring:

- Basic health status endpoint
- Database connectivity verification

## 🧪 Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
CGO_ENABLED=0 GOOS=linux go build -o auth-service-go .
```

## 📝 Environment Variables

| Variable             | Description                       | Default         | Required |
| -------------------- | --------------------------------- | --------------- | -------- |
| `DB_HOST`            | Database host                     | localhost       | Yes      |
| `DB_PORT`            | Database port                     | 5432            | Yes      |
| `DB_USER`            | Database username                 | -               | Yes      |
| `DB_PASSWORD`        | Database password                 | -               | Yes      |
| `DB_NAME`            | Database name                     | auth-service-go | Yes      |
| `DB_SSLMODE`         | PostgreSQL SSL mode               | disable         | No       |
| `DB_CHANNEL_BINDING` | PostgreSQL channel binding        | disable         | No       |
| `JWT_SECRET`         | HS256 JWT signing secret (legacy) | -               | Yes      |
| `OIDC_ISSUER`        | OIDC Issuer URL (your base URL)   | -               | Yes      |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues or have questions, please create an issue in the GitHub repository.
