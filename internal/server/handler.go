package server

import (
	"context"

	pb "github.com/esousa97/godistributedkv/api/proto"
	"github.com/esousa97/godistributedkv/internal/cluster"
	"github.com/esousa97/godistributedkv/internal/storage"
)

// Server implements the gRPC KeyValue service.
type Server struct {
	pb.UnimplementedKeyValueServer
	store   *storage.Store
	cluster *cluster.Manager
	nodeID  string
}

// NewServer creates and returns a new Server instance.
func NewServer(store *storage.Store, cluster *cluster.Manager, nodeID string) *Server {
	return &Server{
		store:   store,
		cluster: cluster,
		nodeID:  nodeID,
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

// Set handles the gRPC Set request. Only the leader can perform Set.
func (s *Server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	if !s.cluster.IsLeader() {
		return &pb.SetResponse{
			Success:    false,
			LeaderHint: s.cluster.GetLeader(),
		}, nil
	}

	s.store.Set(req.GetKey(), req.GetValue())
	return &pb.SetResponse{
		Success: true,
	}, nil
}

// Delete handles the gRPC Delete request.
func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if !s.cluster.IsLeader() {
		return &pb.DeleteResponse{
			Success: false,
		}, nil
	}

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

// RequestVote handles election votes.
func (s *Server) RequestVote(ctx context.Context, req *pb.VoteRequest) (*pb.VoteResponse, error) {
	return s.cluster.HandleRequestVote(req), nil
}

// Heartbeat handles leader heartbeats.
func (s *Server) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	return s.cluster.HandleHeartbeat(req), nil
}
