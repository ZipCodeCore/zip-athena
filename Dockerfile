# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o athena .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache sqlite-libs ca-certificates

# Create app directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/athena .

# Copy static files
COPY --from=builder /app/static ./static

# Copy learning resources
COPY --from=builder /app/learning ./learning

# Create data directory for SQLite database
RUN mkdir -p /data

# Set environment variable for database path
ENV DB_PATH=/data/athena.db

# Expose port
EXPOSE 9206

# Run the application
CMD ["./athena"]
