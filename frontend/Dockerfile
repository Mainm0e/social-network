# Build React frontend
FROM node:14 AS frontend-builder

WORKDIR /app

# Copy the frontend source code into the container
COPY . .
# Install dependencies and build the React frontend
RUN npm install
RUN npm run build

# Copy the built React frontend from the previous stage
COPY --from=frontend-builder /app/frontend/build /usr/share/nginx/html

# Replace the default Nginx configuration with a custom one that listens on port 3000
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose the port for the frontend
EXPOSE 3000

# Start the Nginx server to serve the frontend
CMD ["nginx", "-g", "daemon off;"]
