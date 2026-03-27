# godistributedkv

A basic distributed Key-Value Store written in Go, using gRPC for external communication and thread-safe in-memory storage.

## Project Structure

- `api/proto/`: Contains the gRPC contract definition (`kv.proto`).
- `cmd/server/`: Application entry point.
- `internal/server/`: gRPC handler implementation.
- `internal/storage/`: Thread-safe in-memory storage logic.

## Prerequisites

1. **Go 1.21+**
2. **Protocol Buffers Compiler (protoc)**
3. **Go gRPC Plugins**:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

## Generating gRPC Code

To generate Go files from the `.proto` definition, run:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/proto/kv.proto
```

## How to Run

1. **Unit Tests**:
   ```bash
   go test ./internal/storage/... -v
   ```

2. **Single Node**:
   ```bash
   go run cmd/server/main.go --addr :50051
   ```

3. **Multi-Node Cluster**:
   Open multiple terminals and start nodes pointing to each other:
   
   **Node 1**:
   ```bash
   go run cmd/server/main.go --addr :50051 --peers :50052,:50053
   ```
   
   **Node 2**:
   ```bash
   go run cmd/server/main.go --addr :50052 --peers :50051,:50053
   ```

The server will log health check status (HEALTHY/UNHEALTHY) for each configured peer every 5 seconds.

## Implemented Features

- **Set(key, value)**: Stores a key and its value.
- **Get(key)**: Retrieves the value for a given key.
- **Delete(key)**: Removes a key from the storage.
- **Ping**: Health check endpoint for cluster members.
- **Cluster Management**: Background health checks for a static list of peers.
- **Thread-safety**: Protected via `sync.RWMutex`.
- **gRPC Reflection**: Enabled to facilitate testing with tools like `grpcurl`.
