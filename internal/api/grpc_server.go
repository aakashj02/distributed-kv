package api

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/aakashj02/distributed-kv/internal/hashing"
	"github.com/aakashj02/distributed-kv/internal/store"
	"github.com/aakashj02/distributed-kv/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server implements the gRPC interface and holds cluster state
type Server struct {
	kv.UnimplementedKeyValueStoreServer
	Storage store.Store
	NodeID  string            // Who am I? (e.g., "localhost:5001")
	Ring    *hashing.HashRing // The cluster map
}

func (s *Server) Put(ctx context.Context, req *kv.PutRequest) (*kv.PutResponse, error) {
	// 1. Ask the ring who owns this data
	targetNode := s.Ring.GetNode(req.GetKey())

	// 2. If it belongs to us, save it locally
	if targetNode == s.NodeID {
		fmt.Printf("[Local] Saving key: '%s'\n", req.GetKey())
		err := s.Storage.Set(req.GetKey(), req.GetValue())
		if err != nil {
			return &kv.PutResponse{Success: false}, err
		}
		return &kv.PutResponse{Success: true}, nil
	}

	// 3. If it belongs to someone else, forward it!
	fmt.Printf("[Routing] Key '%s' belongs to %s. Forwarding...\n", req.GetKey(), targetNode)
	return s.forwardPut(ctx, targetNode, req)
}

func (s *Server) Get(ctx context.Context, req *kv.GetRequest) (*kv.GetResponse, error) {
	targetNode := s.Ring.GetNode(req.GetKey())

	if targetNode == s.NodeID {
		fmt.Printf("[Local] Retrieving key: '%s'\n", req.GetKey())
		val, err := s.Storage.Get(req.GetKey())
		if err != nil {
			return &kv.GetResponse{Value: "", Found: false}, nil
		}
		return &kv.GetResponse{Value: val, Found: true}, nil
	}

	fmt.Printf("[Routing] Key '%s' is on %s. Fetching...\n", req.GetKey(), targetNode)
	return s.forwardGet(ctx, targetNode, req)
}

// --- Helper Functions to act as a Client to other nodes ---

func (s *Server) forwardPut(ctx context.Context, targetAddress string, req *kv.PutRequest) (*kv.PutResponse, error) {
	conn, err := grpc.Dial(targetAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := kv.NewKeyValueStoreClient(conn)
	forwardCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return client.Put(forwardCtx, req)
}

func (s *Server) forwardGet(ctx context.Context, targetAddress string, req *kv.GetRequest) (*kv.GetResponse, error) {
	conn, err := grpc.Dial(targetAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := kv.NewKeyValueStoreClient(conn)
	forwardCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return client.Get(forwardCtx, req)
}

// StartGRPCServer now accepts the NodeID and HashRing
func StartGRPCServer(port string, storage store.Store, nodeID string, ring *hashing.HashRing) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	kv.RegisterKeyValueStoreServer(grpcServer, &Server{
		Storage: storage,
		NodeID:  nodeID,
		Ring:    ring,
	})

	return grpcServer.Serve(listener)
}