package main

import (
	"log"
	"net"

	pb "github.com/esousa97/godistributedkv/api/proto"
	"github.com/esousa97/godistributedkv/internal/server"
	"github.com/esousa97/godistributedkv/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize in-memory storage.
	kvStore := storage.NewStore()

	// Initialize gRPC server.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	srv := server.NewServer(kvStore)
	pb.RegisterKeyValueServer(grpcServer, srv)

	// Enable reflection for tools like grpcurl.
	reflection.Register(grpcServer)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
