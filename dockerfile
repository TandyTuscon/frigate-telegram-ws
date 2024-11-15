# Use a minimal base image with Go
FROM golang:1.20-alpine as builder

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set working directory
WORKDIR /app

# Clone the repository
RUN apk add --no-cache git
RUN git clone https://github.com/TandyTuscon/frigate-telegram-ws.git .

# Build the Go application
RUN go build -o main main.go

# Use a minimal image for runtime
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the configuration file
COPY --from=builder /app/config.yml .

# Expose the port if needed
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
