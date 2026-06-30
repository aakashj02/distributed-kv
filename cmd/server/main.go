package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aakashj02/distributed-kv/internal/api"
	"github.com/aakashj02/distributed-kv/internal/hashing"
	"github.com/aakashj02/distributed-kv/internal/store"
)

func main() {
	port := flag.String("port", "5001", "The port for the gRPC server to listen on")
	flag.Parse()

	nodeID := "localhost:" + *port

	// 1. Initialize the Hash Ring (3 virtual replicas per node)
	ring := hashing.NewHashRing(3)
	
	// 2. Register our two servers in the cluster map
	ring.AddNode("localhost:5001")
	ring.AddNode("localhost:5002")

	fmt.Printf("Starting Distributed Node Engine -> %s...\n", nodeID)
	
	localDB := store.NewMemStore()

	// 3. Pass the ID and Ring into the server
	err := api.StartGRPCServer(*port, localDB, nodeID, ring)
	if err != nil {
		log.Fatalf("Failed to bind gRPC server: %v", err)
	}
}