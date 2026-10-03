# Stage 1: Build the main Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git, SSL certificates, and tzdata
RUN apk add --no-cache ca-certificates git tzdata

# Create non-root user for security
RUN adduser -D -g '' appuser

# Create writable data directory owned by appuser
RUN mkdir -p /data && chown appuser:appuser /data

# Create /etc/nsswitch.conf — required by Go's DNS resolver in scratch images.
# Without this the runtime crashes with SIGSEGV when resolving hostnames.
RUN echo 'hosts: files dns' > /etc/nsswitch.conf

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the main binary with optimizations
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o authforest .

# Stage 2: Build the healthcheck binary (static, no CGO).
# A dedicated binary is needed because the scratch runtime has no shell or
# dynamically-linked tools (wget/curl require musl libc which isn't in scratch).
# Only uses stdlib (net/http, os, fmt) — no external deps needed.
FROM golang:1.24-alpine AS healthcheck-builder

WORKDIR /app

COPY --from=builder /app/go.mod /app/go.sum ./
COPY cmd/healthcheck/ ./cmd/healthcheck/

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o /healthcheck \
    ./cmd/healthcheck/

# Stage 3: Create a minimal runtime image
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/nsswitch.conf /etc/nsswitch.conf

# Copy the main application binary
COPY --from=builder /app/authforest /authforest

# Copy the static healthcheck binary
COPY --from=healthcheck-builder /healthcheck /healthcheck

# Copy config files and templates
COPY --from=builder /app/config/files/ /config/files/
COPY --from=builder /app/templates/ /templates/

# Copy writable data directory for RSA key persistence
COPY --from=builder --chown=appuser:appuser /data /data

# Use non-root user
USER appuser

# Probe /health every 30s; allow 15s start-up grace period before first probe.
# Returns 200 OK when the app and DB are ready, 503 otherwise.
HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
    CMD ["/healthcheck"]

# Expose port
EXPOSE 8080

# Command to run the app
ENTRYPOINT ["/authforest"]
