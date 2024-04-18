# Use the official Go image which includes necessary tools to build Go projects.
# Make sure to match the Go version with what you need for your project.
FROM golang:1.19-buster as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o my-go-service ./cmd/server

# Start a new stage from scratch
FROM debian:buster-slim

WORKDIR /

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/my-go-service .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./my-go-service"]
