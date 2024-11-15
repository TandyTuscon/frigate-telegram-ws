# Use a Go image for building
FROM golang:1.20-alpine as builder

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set working directory
WORKDIR /app

# Install Git to pull repository
RUN apk add --no-cache git

# Clone the repository
RUN git clone https://github.com/TandyTuscon/frigate-telegram-ws.git .

# Download dependencies
RUN go mod download

# Build the Go application
RUN go build -o frigate-telegram-ws main.go

# Use a minimal runtime image
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the built application and configuration
COPY --from=builder /app/frigate-telegram-ws .
COPY --from=builder /app/config.yml .

# Set permissions for the binary (optional)
RUN chmod +x ./frigate-telegram-ws

# Run the application
CMD ["./frigate-telegram-ws"]
