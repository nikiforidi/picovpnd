# Use official Golang image as build stage
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Copy go.mod and go.sum if present, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main main.go

# Use a minimal base image for running
FROM debian:bookworm-slim

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/main .

# Run the binary as CMD
# CMD ["./main"]