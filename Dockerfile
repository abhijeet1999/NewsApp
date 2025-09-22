# Build stage
FROM golang:1.25-alpine AS builder

# Install git, ca-certificates, and gcc for SQLite compilation
RUN apk add --no-cache git ca-certificates gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -tags sqlite3 -a -installsuffix cgo -o main ./cmd/main

# Final stage
FROM alpine:latest

# Install ca-certificates and sqlite for runtime
RUN apk --no-cache add ca-certificates sqlite

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy input file
COPY --from=builder /app/cmd/main/input.txt .

# Create output and data directories
RUN mkdir -p output data

# Expose port (if needed for future web interface)
EXPOSE 8080

# Set environment variable for SQLite
ENV NEWS_API_KEY=""

# Run the application
CMD ["./main"]
