# Stage 1: Build the Next.js frontend
FROM node:18 AS frontend-builder

# Set the working directory
WORKDIR /app/frontend

# Copy the frontend package.json and package-lock.json
COPY site/package*.json ./

# Install frontend dependencies
RUN npm install

# Copy the rest of the frontend source code
COPY site/ .

# Build the frontend
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.23 AS backend-builder

WORKDIR /app/backend

# Copy the backend go.mod and go.sum
COPY go.mod go.sum ./

# Download backend dependencies
RUN go mod download

# Copy the rest of the backend source code
COPY . .

# Build the backend
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/main.go

# Stage 3: Create the final image
FROM nginx:alpine

# Copy the built backend from the backend-builder stage
COPY --from=backend-builder /app/backend/server /usr/bin/server

# Copy the built frontend from the frontend-builder stage
COPY --from=frontend-builder /app/frontend/.next /usr/share/nginx/html

# Copy the Nginx configuration file
COPY nginx.conf /etc/nginx/nginx.conf

# Expose the port Nginx is running on
EXPOSE 80

# Start the backend server and Nginx
CMD ["sh", "-c", "/usr/bin/server & nginx -g 'daemon off;'"]
