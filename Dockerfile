# Build stage
FROM golang:1.16-alpine AS builder

# Install git and SSL certificates
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate SQL code
RUN go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
RUN sqlc generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

# Install SSL certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]