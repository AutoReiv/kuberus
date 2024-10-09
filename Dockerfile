# Stage 1: Build the Next.js frontend
FROM node:18 AS frontend-builder

WORKDIR /app/frontend

COPY site/package*.json ./
RUN npm install

# Set environment variables
ENV NEXT_PUBLIC_API_URL=

COPY site/ .
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.23 AS backend-builder

WORKDIR /app/backend
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/main.go

# Final stage
FROM node:18-alpine

# Install Nginx
RUN apk add --no-cache nginx

# Copy the built backend
COPY --from=backend-builder /app/backend/server /usr/bin/server

# Copy the built frontend
COPY --from=frontend-builder /app/frontend/.next /app/frontend/.next
COPY --from=frontend-builder /app/frontend/node_modules /app/frontend/node_modules
COPY --from=frontend-builder /app/frontend/package.json /app/frontend/package.json

# Copy the Nginx configuration file
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80 for Nginx
EXPOSE 80

# Start the backend server, Next.js, and Nginx
CMD ["sh", "-c", "if [ ! -f /root/.kube/config ]; then echo 'Kubeconfig not found. Exiting.' >&2; exit 1; else /usr/bin/server & cd /app/frontend && npm start & nginx -g 'daemon off;'; fi"]
