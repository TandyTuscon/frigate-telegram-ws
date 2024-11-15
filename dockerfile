# Use a minimal base image with Go for building the application
FROM golang:1.20-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set working directory inside the container
WORKDIR /app

# Install Git for cloning the repository
RUN apk add --no-cache git

# Clone the repository
RUN git clone https://github.com/TandyTuscon/frigate-telegram-ws.git .

# Download dependencies and build the Go application
RUN go mod tidy && go build -o main main.go

# Use a minimal base image for the runtime environment
FROM alpine:latest

# Set working directory inside the runtime container
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the configuration file from the builder stage
COPY --from=builder /app/config.yml .

# Expose the port if needed (optional, adjust as necessary)
EXPOSE 8080

# Set the default command to run the application
CMD ["./main"]
