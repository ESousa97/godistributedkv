# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install protoc if needed for CI (though we committed .pb.go)
# RUN apk add --no-cache protobuf

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o godistributedkv ./cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/godistributedkv .

# Expose gRPC port
EXPOSE 50051

# Default command
ENTRYPOINT ["./godistributedkv"]
