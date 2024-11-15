# Use a minimal base image with Go
FROM golang:1.20-alpine as builder

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy all local files to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal image for runtime
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the configuration file
COPY config.yml .

# Command to run the executable
CMD ["./main"]
