FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application (static binary)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o www2rss ./cmd/www2rss

# Install golang-migrate CLI
RUN go install -tags "sqlite3" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite curl

# Set working directory
WORKDIR /root/

# Copy app binary from builder
COPY --from=builder /app/www2rss .

# Copy migrate binary from builder
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Copy templates and migrations
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/db/migrations ./db/migrations

# Copy entrypoint script
COPY scripts/entrypoint.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

# Create data directory for SQLite database
RUN mkdir -p /root/data

# Expose port
EXPOSE 8080

# Healthcheck using curl
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run migrations and start the application
CMD ["./entrypoint.sh"]
