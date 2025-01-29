package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"github.com/panjek26/k8s-workload-analyzer/pkg/analyzer"
	"github.com/panjek26/k8s-workload-analyzer/pkg/api"
)

func main() {
	// Get DeepSeek API key from environment
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		log.Fatal("DEEPSEEK_API_KEY environment variable is required")
	}

	// Create gRPC server
	server := grpc.NewServer()
	analyzerService := analyzer.NewAnalyzerService(apiKey)
	api.RegisterWorkloadAnalyzerServer(server, analyzerService)

	// Start listening
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Starting gRPC server on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}