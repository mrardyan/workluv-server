# Single stage build for better DigitalOcean compatibility
FROM golang:1.23-alpine

# Install git, ca-certificates, and wget
RUN apk add --no-cache git ca-certificates wget

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Build the migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate

# Ensure directories exist and have proper permissions
RUN mkdir -p ./src ./migration

# Verify files were copied correctly
RUN echo "=== Container Structure ===" && \
    ls -la && \
    echo "=== Src Directory ===" && \
    ls -la ./src/ && \
    echo "=== Migration Directory ===" && \
    ls -la ./migration/

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health/simple || exit 1

# Run the application directly
CMD ["./main"]