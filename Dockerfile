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

# Enable CGO and build the Go application
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o myapp cmd/main.go

# Use a minimal base image for the final container
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the built Go application from the builder stage
COPY --from=builder /app/myapp .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]
