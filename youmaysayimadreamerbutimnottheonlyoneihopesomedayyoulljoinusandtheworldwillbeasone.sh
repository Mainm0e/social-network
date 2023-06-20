#!/bin/bash

# Function to handle termination of the script
trap 'kill $(jobs -p)' EXIT

# Start the backend
cd backend
echo "Starting backend server..."
go run . &
cd ..

# Start the frontend
cd frontend
echo "Starting frontend React app..."
npx react-scripts start & # Using npx to run react-scripts
cd ..

# Wait for all background processes to finish (or be killed)
wait