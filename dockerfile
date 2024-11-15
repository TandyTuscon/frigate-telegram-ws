# Use a minimal base image with Go
FROM golang:1.20-alpine as builder

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go application
RUN go build -o main main.go

# Use a minimal image for runtime
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the configuration file
COPY config.yml .

# Expose the port if needed (replace with your app's port)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
