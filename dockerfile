# Use the official Go image as the base image
FROM golang:1.23-bullseye AS build-stage

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download and cache the dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/main ./cmd/server/main.go

# Use a smaller base image for the final stage
FROM golang:1.23-alpine AS release-stage

# Install Air for live-reloading
RUN go install github.com/air-verse/air@latest

# Install curl for health checks
RUN apk update && apk add --no-cache curl

# Set the current working directory inside the container
WORKDIR /app

# Copy the source code into the final image
COPY . .

# Set the command to run the application
CMD ["air"]
