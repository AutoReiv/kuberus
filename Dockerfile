# Use the official Golang image as the base image for building the application
FROM golang:1.23-alpine as builder

# Install necessary build dependencies
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN GOOS=linux GOARCH=amd64 go build -o myapp cmd/main.go

# Use a minimal base image for the final container
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Create a non-root user
RUN adduser -S myappuser

# Set the working directory inside the container
WORKDIR /root/

# Copy the built Go application from the builder stage
COPY --from=builder /app/myapp .

# Change ownership of the application binary
RUN chown myappuser myapp

# Switch to the non-root user
USER myappuser

# Expose the port the app runs on
EXPOSE 8080

# Add a health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 CMD wget --spider http://localhost:8080/health || exit 1

# Command to run the executable
CMD ["./myapp"]
