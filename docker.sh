#!/bin/bash

# Build the backend image
docker build -t backend-image:latest -f backend/Dockerfile backend

# Build the frontend image
docker build -t frontend-image:latest -f frontend/Dockerfile frontend

# Run the backend container and publish port 8080
docker run -d -p 8080:8080 --name backend-container backend-image:latest

# Run the frontend container and publish port 3000
docker run -d -p 3000:3000 --name frontend-container  frontend-image:latest
