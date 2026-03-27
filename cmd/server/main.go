package main

import (
	"log"
	"net"

	pb "github.com/esousa97/godistributedkv/api/proto"
	"github.com/esousa97/godistributedkv/internal/cluster"
	"github.com/esousa97/godistributedkv/internal/config"
	"github.com/esousa97/godistributedkv/internal/server"
	"github.com/esousa97/godistributedkv/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting godistributedkv...")

	// Load configuration.
	cfg := config.Load()
	log.Printf("Config loaded - Addr: %s, Peers: %v", cfg.Addr, cfg.Peers)

	// Initialize in-memory storage.
	kvStore := storage.NewStore()

	// Initialize cluster manager.
	clusterMgr := cluster.NewManager(cfg.Addr, cfg.Peers)
	clusterMgr.Start()

	// Initialize gRPC server.
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Fatalf("CRITICAL: failed to listen on %s: %v", cfg.Addr, err)
	}

	grpcServer := grpc.NewServer()
	srv := server.NewServer(kvStore, clusterMgr, cfg.Addr)
	pb.RegisterKeyValueServer(grpcServer, srv)

	// Enable reflection for tools like grpcurl.
	reflection.Register(grpcServer)

	log.Printf("gRPC server is listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("CRITICAL: failed to serve: %v", err)
	}
}
