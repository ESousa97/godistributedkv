package server

import (
	"context"

	pb "github.com/esousa97/godistributedkv/api/proto"
	"github.com/esousa97/godistributedkv/internal/storage"
)

// Server implements the gRPC KeyValue service.
type Server struct {
	pb.UnimplementedKeyValueServer
	store  *storage.Store
	nodeID string
}

// NewServer creates and returns a new Server instance.
func NewServer(store *storage.Store, nodeID string) *Server {
	return &Server{
		store:  store,
		nodeID: nodeID,
	}
}

// Get handles the gRPC Get request.
func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	val, found := s.store.Get(req.GetKey())
	return &pb.GetResponse{
		Value: val,
		Found: found,
	}, nil
}

// Set handles the gRPC Set request.
func (s *Server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	s.store.Set(req.GetKey(), req.GetValue())
	return &pb.SetResponse{
		Success: true,
	}, nil
}

// Delete handles the gRPC Delete request.
func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.store.Delete(req.GetKey())
	return &pb.DeleteResponse{
		Success: true,
	}, nil
}

// Ping handles the gRPC Ping request for cluster health checks.
func (s *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		NodeId:  s.nodeID,
		Healthy: true,
	}, nil
}
