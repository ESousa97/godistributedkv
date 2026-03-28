package config

import (
	"flag"
	"strings"
)

// Config holds the application configuration.
type Config struct {
	Addr    string
	Peers   []string
	WALPath string
}

// Load parses flags and returns the configuration.
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
