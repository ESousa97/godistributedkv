package cluster

import (
	"context"
	"log"
	"sync"
	"time"

	pb "github.com/esousa97/godistributedkv/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Peer represents a remote node in the cluster.
type Peer struct {
	Addr    string
	Client  pb.KeyValueClient
	Conn    *grpc.ClientConn
	Healthy bool
}

// Manager manages the cluster peers and their health status.
type Manager struct {
	mu     sync.RWMutex
	peers  map[string]*Peer
	nodeID string
}

// NewManager initializes a new cluster Manager.
func NewManager(nodeID string, peerAddrs []string) *Manager {
	m := &Manager{
		peers:  make(map[string]*Peer),
		nodeID: nodeID,
	}

	for _, addr := range peerAddrs {
		m.peers[addr] = &Peer{
			Addr:    addr,
			Healthy: false,
		}
	}

	return m
}

// Start initiates background health checks for all peers.
func (m *Manager) Start() {
	for _, p := range m.peers {
		peer := p // Local copy for closure safety
		go m.healthCheckLoop(peer)
	}
}

func (m *Manager) healthCheckLoop(peer *Peer) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// First ping attempt immediately
	m.pingPeer(peer)

	for {
		<-ticker.C
		m.pingPeer(peer)
	}
}

func (m *Manager) pingPeer(peer *Peer) {
	if peer.Conn == nil {
		// Using grpc.Dial for better compatibility with older versions or WSL quirks
		conn, err := grpc.NewClient(peer.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("ERROR: Could not create gRPC client for %s: %v", peer.Addr, err)
			return
		}
		peer.Conn = conn
		peer.Client = pb.NewKeyValueClient(conn)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := peer.Client.Ping(ctx, &pb.PingRequest{NodeId: m.nodeID})

	m.mu.Lock()
	defer m.mu.Unlock()

	if err != nil {
		if peer.Healthy {
			log.Printf("ALERT: Peer %s is now UNHEALTHY: %v", peer.Addr, err)
		}
		peer.Healthy = false
		return
	}

	if !peer.Healthy {
		log.Printf("OK: Peer %s is now HEALTHY (ID: %s)", peer.Addr, resp.GetNodeId())
	}
	peer.Healthy = resp.GetHealthy()
}

// GetPeers returns the current status of all peers.
func (m *Manager) GetPeers() map[string]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := make(map[string]bool)
	for addr, p := range m.peers {
		status[addr] = p.Healthy
	}
	return status
}
