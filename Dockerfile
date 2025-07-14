# Build stage
FROM golang:1.19-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the applications
RUN go build -o fedramp-server cmd/server/main.go
RUN go build -o gocomply_fedramp cli/gocomply_fedramp/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 fedramp && \
    adduser -D -u 1000 -G fedramp fedramp

# Set working directory
WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/fedramp-server /app/
COPY --from=builder /app/gocomply_fedramp /app/

# Copy web assets
COPY --from=builder /app/web /app/web

# Copy bundled resources
COPY --from=builder /app/bundled /app/bundled

# Create directories for data
RUN mkdir -p /app/data /app/logs && \
    chown -R fedramp:fedramp /app

# Switch to non-root user
USER fedramp

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Default command
CMD ["/app/fedramp-server"] 