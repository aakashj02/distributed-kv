package main

import (
	"fmt"
	"log"

	"github.com/aakashj02/distributed-kv/internal/hashing"
	"github.com/aakashj02/distributed-kv/internal/store"
)

func main() {
	fmt.Println("Initializing Local Storage Engine...")
	db := store.NewMemStore()

	// Test a local set/get sequence
	err := db.Set("user_session_1", "authenticated_token_xyz")
	if err != nil {
		log.Fatalf("Failed to write to store: %v", err)
	}

	val, _ := db.Get("user_session_1")
	fmt.Printf("[Local Store Success] Retrieved key 'user_session_1': %s\n\n", val)

	fmt.Println("Initializing Distributed Hash Ring Ring...")
	ring := hashing.NewHashRing(3) 
	ring.AddNode("Server-Node-A")
	ring.AddNode("Server-Node-B")
	ring.AddNode("Server-Node-C")

	// Map arbitrary sample data to see which cluster machine manages it
	sampleKeys := []string{"profile_pic", "user_metadata", "auth_token", "cached_query_1"}
	for _, key := range sampleKeys {
		assignedNode := ring.GetNode(key)
		fmt.Printf("Data Key '%s' is routed to -> %s\n", key, assignedNode)
	}
}