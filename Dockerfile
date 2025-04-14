# Build stage
FROM golang:1.23.4-alpine AS builder

# Install required build tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o zodc-service-masterflow ./cmd/main.go

# Final stage
FROM alpine:3.19

# Set TimeZone
ENV TZ=Asia/Ho_Chi_Minh

# Add non root user
RUN adduser -D appuser

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/zodc-service-masterflow .

# Use non root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./zodc-service-masterflow"]