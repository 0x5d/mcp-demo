package mcp

import (
	"mcp-server/internal/k8s"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type NamespaceParams struct {
	Namespace string `json:"namespace,omitempty" jsonschema:"description:Kubernetes namespace"`
}

func NewServer(kubeClient *k8s.Client) *mcp.Server {

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-server",
		Version: "1.0.0",
	}, nil)

	addK8sTools(server, kubeClient)

	return server
}

func NewStdioTransport() *mcp.StdioTransport {
	return mcp.NewStdioTransport()
}
