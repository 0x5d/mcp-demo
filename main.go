package main

import (
	"context"
	"log"

	"mcp-server/internal/k8s"
	"mcp-server/internal/mcp"
)

func main() {
	// Initialize Kubernetes client
	kubeClient, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	server := mcp.NewServer(kubeClient)

	log.Println("MCP Server started")
	server.Run(context.Background(), mcp.NewStdioTransport())
}
