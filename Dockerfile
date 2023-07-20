# Stage 1: Build Golang backend
FROM golang:1.20 AS backend-builder

WORKDIR /app/backend

# Copy the backend source code into the container
COPY backend .

# Build the Golang backend
RUN go build -o /server-bin

# Stage 2: Build React frontend
FROM node:14 AS frontend-builder

WORKDIR /app/frontend

# Copy the frontend source code into the container
COPY frontend .

# Install dependencies and build the React frontend
RUN npm install
RUN npm run build

# Stage 3: Create final image for backend
FROM golang:1.20 AS backend

WORKDIR /app

# Copy the Golang backend executable from the previous stage
COPY --from=backend-builder /server-bin .

# Expose the port for the backend
EXPOSE 8080

# Start the backend
CMD ["./server-bin"]

# Stage 4: Create final image for frontend
FROM nginx:alpine AS frontend

# Copy the built React frontend from the previous stage
COPY --from=frontend-builder /app/frontend/build /usr/share/nginx/html

# Replace the default Nginx configuration with a custom one that listens on port 3000
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose the port for the frontend
EXPOSE 3000

# Start the Nginx server to serve the frontend
CMD ["nginx", "-g", "daemon off;"]
