# Use Golang base image
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY services/auth-service/go.mod services/auth-service/go.sum shared-libs/go.mod shared-libs/go.sum ./

# Download dependencies
RUN go mod download

# Copy the full application source
COPY services/auth-service /app/auth-service
COPY shared-libs /app/shared-libs

# Change working directory to auth-service
WORKDIR /app/auth-service

# Update go.mod replace path for shared-libs inside Docker
RUN go mod edit -replace=github.com/Aditya-PS-05/NeetChamp/shared-libs=/app/shared-libs
RUN go mod tidy

# Build the Go binary
RUN go build -o auth-service .

# Create a smaller final image
FROM golang:1.22

WORKDIR /app

# Copy the compiled binary from builder stage
COPY --from=builder /app/auth-service/auth-service .

# Run the service
CMD ["/app/auth-service"]
