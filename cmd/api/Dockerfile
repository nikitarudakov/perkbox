# Dockerfile
FROM golang:1.24.0

# Set working directory
WORKDIR /app

# Copy go mod and sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN go build -o user-service ./cmd/api/main.go

# Expose API port
EXPOSE 8080

# Run the binary
CMD ["./user-service"]