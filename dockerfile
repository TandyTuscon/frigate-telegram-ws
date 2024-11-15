# Use a lightweight Go image to build the binary
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifest files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application code to the working directory
COPY . .

# Build the Go application
RUN go build -o /app/main .

# Use a minimal image for runtime
FROM alpine:latest

# Set the working directory inside the runtime container
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the configuration file
COPY config.yml .

# Expose necessary ports (if any)
# EXPOSE 8080

# Set the entrypoint for the application
ENTRYPOINT ["./main"]
