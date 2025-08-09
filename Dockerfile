# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and SSL certificates
RUN apk add --no-cache ca-certificates git tzdata

# Create non-root user for security
RUN adduser -D -g '' appuser

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o auth-service-go .

# Stage 2: Create a minimal runtime image
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy binary from builder
COPY --from=builder /app/auth-service-go /auth-service-go

# Use non-root user
USER appuser

# Health check (adjust the command based on your app's health endpoint)
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#     CMD ["wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health", "||", "exit", "1"]

# Expose port (adjust if your app uses a different port)
EXPOSE 8080

# Command to run the app
ENTRYPOINT ["/auth-service-go"]
