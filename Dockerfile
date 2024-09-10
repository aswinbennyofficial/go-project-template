# Build stage
FROM golang:1.22.6-alpine AS builder
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code and config
COPY . .

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:3.14
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy the config directory
COPY config ./config

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]