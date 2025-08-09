# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and SSL certificates
RUN apk add --no-cache git ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o auth-service-go .

# Stage 2: Create a minimal runtime image
FROM alpine:3.22

WORKDIR /app

# Install CA certificates for HTTPS support
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/auth-service-go .

# Command to run the app
ENTRYPOINT ["./auth-service-go"]
