# Use the official Golang image as the base image with version 1.20.3
FROM golang:1.20.3 as backend-builder

# Set the working directory inside the container to /app
WORKDIR /app

# Copy only the go.mod and go.sum files to the container
# and download the Go dependencies. This step utilizes
# Docker's layer caching to optimize the build process.
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire current directory (including the source code)
# into the container's working directory /app.
COPY . .

# Build the Go backend executable named "main".
RUN go build -o main .

# Expose port 8080 to allow external access to the backend service.
# Note that this only exposes the port within the container, and
# you need to map it to a host port when running the container.
EXPOSE 8080

# Set the default command to run when the container starts.
CMD ["./main"]
