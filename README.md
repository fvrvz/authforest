# Auth Service Go

A robust authentication and user management service built with Go and the Gin web framework. This service provides JWT-based authentication, user registration, and CRUD operations with PostgreSQL integration.

## 🚀 Features

- **JWT Authentication**: Secure token-based authentication with access and refresh tokens
- **User Management**: User registration, login, logout, and profile management
- **Database Integration**: PostgreSQL with GORM ORM
- **CORS Support**: Configurable Cross-Origin Resource Sharing
- **Middleware**: Authentication middleware for protected routes
- **Docker Support**: Containerized deployment with Docker and Docker Compose
- **Configuration Management**: YAML-based configuration with environment variable support

## 🛠 Tech Stack

- **Language**: Go 1.24+
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt/v5)
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
├── db/                     # Database connection and setup
├── dto/                    # Data Transfer Objects
├── helpers/                # Utility functions (JWT, date normalization)
├── middlewares/            # HTTP middlewares
├── models/                 # Database models
├── services/               # Business logic layer
├── server/                 # Server setup and routing
├── postman-collections/    # API testing collections
├── docker-compose.yml      # Multi-container Docker configuration
├── Dockerfile              # Container image definition
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
   JWT_SECRET=your-jwt-secret-key
   ```

5. **Update configuration (optional)**

   Modify `config/files/config.yml` if needed for custom settings like CORS origins or JWT expiry times.

6. **Build and run the application**
   ```bash
   go build -o auth-service-go .
   ./auth-service-go
   ```

The application will start on port 8080 (configurable in `config.yml`).

## 📚 API Endpoints

### Public Endpoints

- `POST /api/v1/users/register` - User registration
- `POST /api/v1/auth/login` - User login

### Protected Endpoints (Require Authentication)

- `GET /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `GET /api/v1/users/` - Get all users
- `GET /api/v1/users/:userId` - Get specific user
- `DELETE /api/v1/users/:userId` - Delete user

### API Testing

Import the Postman collection from `postman-collections/AuthServiceGo.postman_collection.json` for easy API testing.

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
```

## 🔐 Security Features

- **JWT Tokens**: Short-lived access tokens with longer-lived refresh tokens
- **Password Hashing**: Secure password storage
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

| Variable      | Description        | Default         | Required |
| ------------- | ------------------ | --------------- | -------- |
| `DB_HOST`     | Database host      | localhost       | Yes      |
| `DB_PORT`     | Database port      | 5432            | Yes      |
| `DB_USER`     | Database username  | -               | Yes      |
| `DB_PASSWORD` | Database password  | -               | Yes      |
| `DB_NAME`     | Database name      | auth-service-go | Yes      |
| `JWT_SECRET`  | JWT signing secret | -               | Yes      |

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
