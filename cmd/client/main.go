package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aakashj02/distributed-kv/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Attempting to connect to Node A on port 5001...")

	// 1. Dial the server over the local network
	conn, err := grpc.Dial("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// 2. Create a new client using the generated Protocol Buffer code
	client := kv.NewKeyValueStoreClient(conn)

	// Create a timeout context so we don't wait forever
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 3. Send a PUT request over the network
	fmt.Println("\nSending PUT request: Key='my_resume', Value='is looking great'")
	putRes, err := client.Put(ctx, &kv.PutRequest{
		Key:   "my_resume",
		Value: "is looking great",
	})
	if err != nil {
		log.Fatalf("Could not put data: %v", err)
	}
	fmt.Printf("Server responded -> Success: %v\n", putRes.GetSuccess())

	// 4. Send a GET request over the network
	fmt.Println("\nSending GET request for Key='my_resume'")
	getRes, err := client.Get(ctx, &kv.GetRequest{
		Key: "my_resume",
	})
	if err != nil {
		log.Fatalf("Could not get data: %v", err)
	}

	if getRes.GetFound() {
		fmt.Printf("Server responded -> Value: '%s'\n", getRes.GetValue())
	} else {
		fmt.Println("Server responded -> Key not found.")
	}
}