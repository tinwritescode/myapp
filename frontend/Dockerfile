# Multi-stage build for React frontend
FROM node:18-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files
COPY package.json bun.lockb ./

# Install dependencies using bun (faster than npm)
RUN npm install -g bun
RUN bun install

# Copy source code
COPY . .

# Set environment variables for production build
ENV VITE_API_URL=https://myapp-1757744589.fly.dev/api/v1

# Build the application
RUN bun run build

# Production stage with nginx
FROM nginx:alpine

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
