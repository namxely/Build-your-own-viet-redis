# VietRedis Server Docker Image
# Multi-stage build for optimized production image

#################################################
# Build Stage
#################################################
FROM golang:1.21-alpine AS builder

# Set build arguments
ARG VERSION=2.0.0-saigon
ARG BUILD_TIME
ARG GIT_COMMIT

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimization flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=${VERSION} -X 'main.buildTime=${BUILD_TIME}' -X main.gitCommit=${GIT_COMMIT} -X main.buildBy=Docker" \
    -a -installsuffix cgo \
    -o vietredis \
    .

#################################################
# Production Stage
#################################################
FROM alpine:3.18

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    dumb-init \
    su-exec \
    redis \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 999 vietredis && \
    adduser -D -u 999 -G vietredis -s /bin/sh vietredis

# Set timezone to Vietnam
ENV TZ=Asia/Ho_Chi_Minh
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Create directories
RUN mkdir -p \
    /data \
    /var/log/vietredis \
    /etc/vietredis \
    && chown -R vietredis:vietredis /data /var/log/vietredis

# Copy binary from builder stage
COPY --from=builder /build/vietredis /usr/local/bin/vietredis

# Copy configuration files
COPY vietredis.conf /etc/vietredis/vietredis.conf
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

# Make scripts executable
RUN chmod +x /usr/local/bin/vietredis /usr/local/bin/docker-entrypoint.sh

# Set working directory
WORKDIR /data

# Expose ports
EXPOSE 6379 9121

# Set volume
VOLUME ["/data", "/var/log/vietredis"]

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD redis-cli -p 6379 ping || exit 1

# Use dumb-init for proper signal handling
ENTRYPOINT ["dumb-init", "--"]

# Default command
CMD ["/usr/local/bin/docker-entrypoint.sh"]

# Labels for metadata
LABEL maintainer="Namxely <dev.namxely@gmail.com>" \
      org.opencontainers.image.title="VietRedis Server" \
      org.opencontainers.image.description="🇻🇳 High-Performance Redis Implementation by Namxely" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.url="https://namxely.github.io/viet-redis-server" \
      org.opencontainers.image.source="https://github.com/namxely/viet-redis-server" \
      org.opencontainers.image.documentation="https://github.com/namxely/viet-redis-server/blob/main/README.md" \
      org.opencontainers.image.vendor="Namxely Development" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.authors="Namxely (@namxely)"
