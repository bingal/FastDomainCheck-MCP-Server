# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Set GOPROXY for better download speed in China
# ENV GOPROXY=https://goproxy.cn,direct

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o FastDomainCheck-MCP-Server

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates whois bind-tools

# Copy binary from builder
COPY --from=builder /app/FastDomainCheck-MCP-Server .

# Run the server (skip health check in container)
CMD ["./FastDomainCheck-MCP-Server"] 