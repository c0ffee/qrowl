# Start from a Go base image
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a lightweight base image for the final image
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .
COPY logo.png .
# Set the command to run when the container starts
CMD ["./main"]
