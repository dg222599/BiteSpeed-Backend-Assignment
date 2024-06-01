# # Use an official Golang runtime as a parent image
# FROM golang:1.18-alpine AS builder

# # Set the Current Working Directory inside the container
# WORKDIR /app

# # Copy go mod and sum files
# COPY go.mod go.sum ./

# # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# # Copy the source code into the container
# COPY . .

# # Build the Go app
# RUN go build -o main .

# # Start a new stage from scratch
# FROM alpine:latest

# WORKDIR /root/

# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /app/main .

# # Expose port
# EXPOSE 3000

# # Command to run the executable
# CMD ["./main"]

#Testing Dockerfile for deploying on render
# Dockerfile for the api service
FROM golang:1.18-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 3000

ENV DATABASE_URL=postgres://dg22:DRPXdr2GNFhVAjVGhnakA2qEVHOrdJDe@dpg-cpdlqevsc6pc7393u96g-a.oregon-postgres.render.com/bitespeed_z7ix
# Command to run the executable
CMD ["./main"]
