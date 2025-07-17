package mcp

import (
	"context"
	"mcp-server/internal/k8s"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func addK8sTools(server *mcp.Server, kubeClient *k8s.Client) {

	// Add tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "k8s_get_nodes",
		Description: "Get all Kubernetes nodes",
	}, func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[struct{}]) (*mcp.CallToolResultFor[any], error) {
		nodes, err := kubeClient.GetNodes(ctx)
		if err != nil {
			return nil, err
		}
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{
				Text: k8s.FormatNodes(nodes),
			}},
		}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "k8s_get_pods",
		Description: "Get Kubernetes pods",
	}, func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[NamespaceParams]) (*mcp.CallToolResultFor[any], error) {
		namespace := params.Arguments.Namespace
		if namespace == "" {
			namespace = "default"
		}
		pods, err := kubeClient.GetPods(ctx, namespace)
		if err != nil {
			return nil, err
		}
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{
				Text: k8s.FormatPods(pods, namespace),
			}},
		}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "k8s_get_services",
		Description: "Get Kubernetes services",
	}, func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[NamespaceParams]) (*mcp.CallToolResultFor[any], error) {
		namespace := params.Arguments.Namespace
		if namespace == "" {
			namespace = "default"
		}
		services, err := kubeClient.GetServices(ctx, namespace)
		if err != nil {
			return nil, err
		}
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{
				Text: k8s.FormatServices(services, namespace),
			}},
		}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "k8s_get_deployments",
		Description: "Get Kubernetes deployments",
	}, func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[NamespaceParams]) (*mcp.CallToolResultFor[any], error) {
		namespace := params.Arguments.Namespace
		if namespace == "" {
			namespace = "default"
		}
		deployments, err := kubeClient.GetDeployments(ctx, namespace)
		if err != nil {
			return nil, err
		}
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{
				Text: k8s.FormatDeployments(deployments, namespace),
			}},
		}, nil
	})
}
