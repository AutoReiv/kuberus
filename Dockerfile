FROM golang:1.23-alpine AS go-builder

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
RUN GOOS=linux GOARCH=amd64 go build -o k-rbac cmd/main.go

# Stage 2: Build the React frontend
FROM node:18-alpine AS react-builder

# Set the working directory inside the container
WORKDIR /app

# Copy the package.json and package-lock.json files
COPY site/package.json site/package-lock.json ./

# Install the dependencies
RUN npm install

# Copy the rest of the application code
COPY site ./

# Build the React application
RUN npm run build

# Stage 3: Create the final image
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates nginx

# Create a non-root user
RUN adduser -S myappuser

# Set the working directory inside the container
WORKDIR /root/

# Copy the built Go application from the go-builder stage
COPY --from=go-builder /app/k-rbac .

# Copy the built React app from the react-builder stage
COPY --from=react-builder /app/.next /usr/share/nginx/html

# Change ownership of the application binary
RUN chown myappuser k-rbac

# Switch to the non-root user
USER myappuser

# Expose the ports the apps run on
EXPOSE 80

# Add a health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 CMD wget --spider http://localhost:8080/health || exit 1

# Command to run both the Go backend and nginx
CMD ["sh", "-c", "./k-rbac & nginx -g 'daemon off;'"]