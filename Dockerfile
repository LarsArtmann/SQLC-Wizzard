# Multi-stage Docker build for SQLC-Wizard

# Build stage
FROM golang:1.25-alpine AS builder

# Install git for version info
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${GIT_VERSION:-dev} -X main.Commit=${GIT_COMMIT:-unknown} -X main.BuildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
    -o sqlc-wizard \
    cmd/sqlc-wizard/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/sqlc-wizard .

# Copy license and readme
COPY --from=builder /app/LICENSE .
COPY --from=builder /app/README.md .

# Set timezone
ENV TZ=UTC

# Create non-root user
RUN addgroup -S sqlcwizard && adduser -S sqlcwizard -G sqlcwizard
USER sqlcwizard

# Expose volume for configs
VOLUME ["/configs"]

# Set entrypoint
ENTRYPOINT ["./sqlc-wizard"]

# Default command shows help
CMD ["--help"]

# Labels for metadata
LABEL maintainer="Lars Artmann <lars@artmann.email>"
LABEL description="SQLC-Wizard - Interactive CLI wizard for generating sqlc configurations"
LABEL version="${GIT_VERSION:-dev}"
LABEL org.opencontainers.image.title="SQLC-Wizard"
LABEL org.opencontainers.image.description="Interactive CLI wizard for sqlc configurations"
LABEL org.opencontainers.image.vendor="Lars Artmann"
LABEL org.opencontainers.image.licenses="MIT"