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

The server will log health check status and leader election transitions (TERM updates, Leader elected).

## Implemented Features

- **Set(key, value)**: Only allowed on the **Leader** node. Propagated to followers via quorum.
- **Get(key)**: Retrieves the value for a given key from any node.
- **Delete(key)**: Only allowed on the **Leader** node. Propagated via quorum.
- **Majority Consensus Replication (Quorum)**: Operations are only committed after at least N/2 + 1 nodes acknowledge the update.
- **Leader Election**: Automated election using terms and votes when the leader fails.
- **Leader Redirection**: Followers return the address of the current leader for write operations.
- **Heartbeats**: Leader maintains authority via periodic heartbeats.
- **Thread-safety**: Protected via `sync.RWMutex`.
- **gRPC Reflection**: Enabled to facilitate testing with tools like `grpcurl`.
