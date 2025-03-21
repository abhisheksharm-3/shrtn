# Development stage with Air for hot-reload
FROM golang:1.24-alpine AS development

WORKDIR /app

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Create temp directory with proper permissions
RUN mkdir -p tmp && chmod 777 tmp

# Copy go mod files
COPY server/go.mod server/go.sum ./

# Download dependencies
RUN go mod download

# Copy Air configuration
COPY server/.air.toml ./.air.toml

# Expose the port
EXPOSE 8080

# Show debug info and run air
CMD ["sh", "-c", "echo 'Starting with Air' && pwd && ls -la && air -c .air.toml"]

# Builder stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY server/go.mod server/go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY server/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final production stage
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

CMD ["./main"]