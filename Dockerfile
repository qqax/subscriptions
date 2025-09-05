# --- Build stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/server

# --- Final stage (minimal runtime image) ---
FROM scratch

# Copy binary and necessary files from builder
COPY --from=builder /app/main /
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Use a non-root user (optional)
USER 1000:1000

# Expose the application port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/main"]
