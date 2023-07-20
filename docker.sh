# Build the backend image
docker build -t backend-image:latest --target backend .

# Build the frontend image
docker build -t frontend-image:latest --target frontend .

# Run the backend container and publish port 8080
docker run -d -p 8080:8080 backend-image:latest

# Run the frontend container and publish port 3000
docker run -d -p 3000:3000 frontend-image:latest
