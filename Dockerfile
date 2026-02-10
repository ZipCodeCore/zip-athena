# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=1 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags="-w -s" \
    -o athena .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    sqlite-libs \
    ca-certificates \
    tzdata \
    && addgroup -g 1000 athena \
    && adduser -D -u 1000 -G athena athena

# Create app directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/athena .

# Copy static files
COPY --from=builder /app/static ./static

# Copy learning resources
COPY --from=builder /app/learning ./learning

# Copy database migrations
COPY --from=builder /app/database/migrations ./database/migrations

# Create data directory for SQLite database
RUN mkdir -p /data && chown -R athena:athena /data /app

# Switch to non-root user
USER athena

# Default environment variables (can be overridden)
ENV PORT=9206 \
    DB_TYPE=sqlite \
    DB_PATH=/data/athena.db \
    LOG_LEVEL=info \
    LOG_FORMAT=json

# Expose port
EXPOSE 9206

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Run the application
CMD ["./athena"]
