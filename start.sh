#!/bin/bash

# Function to handle termination of the script
trap 'kill $(jobs -p)' EXIT
#export DB_DIALECT=sqlite3
#export DB_DATASOURCE=$(pwd)/backend/db/database.db
#export DB_MIGRATIONS_DIR=$(pwd)/backend/db/migrations
# Start the backend
cd backend
echo "Starting backend server..."
go run . &
cd ..

# Start the frontend
cd frontend
echo "Starting frontend React app..."
npm start & # Using npx to run react-scripts
cd ..

# Wait for all background processes to finish (or be killed)
wait