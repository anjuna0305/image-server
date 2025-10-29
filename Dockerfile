# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o image-server main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS support (if needed)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/image-server .

# Create uploads directory
RUN mkdir -p uploads && chmod 755 uploads

# Expose port
EXPOSE 8000

# Environment variables (can be overridden at runtime)
ENV UPLOAD_DIR_PATH=/app/uploads
ENV SERVER_PORT=:8000

# SECRET_KEY must be provided at runtime via docker run -e or docker-compose

# Run the application
CMD ["./image-server"]

