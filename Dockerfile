# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o www2rss ./cmd/www2rss

# Final stage
FROM alpine:latest

# Install ca-certificates and sqlite
RUN apk --no-cache add ca-certificates sqlite

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/www2rss .

# Create data directory for database
RUN mkdir -p /root/data

# Expose port
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
CMD ["./www2rss"]
