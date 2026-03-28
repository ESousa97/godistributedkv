// Package config handles the application configuration via command-line flags.
package config

import (
	"flag"
	"strings"
)

// Config holds the application configuration parameters.
type Config struct {
	// Addr is the network address (host:port) where the server will listen.
	Addr string
	// Peers is a list of network addresses for other nodes in the cluster.
	Peers []string
	// WALPath is the filesystem path to the Write-Ahead Log file.
	WALPath string
}

// Load parses command-line flags and returns a populated [Config] instance.
// It supports flags for address, peer list (comma-separated), and WAL file path.
func Load() *Config {
	addr := flag.String("addr", ":50051", "Server listen address")
	peersStr := flag.String("peers", "", "Comma-separated list of peer addresses")
	walPath := flag.String("wal-path", "data/kv.log", "Path to the Write-Ahead Log file")
	flag.Parse()

	var peers []string
	if *peersStr != "" {
		peers = strings.Split(*peersStr, ",")
	}

	return &Config{
		Addr:    *addr,
		Peers:   peers,
		WALPath: *walPath,
	}
}
