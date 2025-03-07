# Build stage
FROM golang:1.24 AS builder

# Set working directory
WORKDIR /build

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the application with a specific output name
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /build/main ./app/cmd/main.go

# Run stage
FROM alpine:latest

# Install CA certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from the build stage with a specific name
COPY --from=builder /build/main .

# Copy the config folder
COPY --from=builder /build/config ./config

# Ensure the binary is executable
RUN chmod +x /app/main

# Expose the port your application runs on
EXPOSE 8080

# Run the application with the correct binary name
CMD ["./main"]
