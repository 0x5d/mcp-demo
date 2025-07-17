package k8s

import (
	"encoding/json"
	"fmt"
)

func FormatNodes(nodes []NodeInfo) string {
	if len(nodes) == 0 {
		return "No nodes found"
	}

	output := fmt.Sprintf("Found %d node(s):\n\n", len(nodes))
	for _, node := range nodes {
		output += fmt.Sprintf("**%s** (%s)\n", node.Name, node.Status)
		output += fmt.Sprintf("  - Version: %s\n", node.Node.Status.NodeInfo.KubeletVersion)
		output += fmt.Sprintf("  - OS: %s\n", node.Node.Status.NodeInfo.OSImage)
		output += fmt.Sprintf("  - Kernel: %s\n", node.Node.Status.NodeInfo.KernelVersion)
		output += fmt.Sprintf("  - Created: %s\n", node.CreationTimestamp.Format("2006-01-02 15:04:05"))

		if len(node.Node.Status.Addresses) > 0 {
			output += "  - Addresses:\n"
			for _, addr := range node.Node.Status.Addresses {
				output += fmt.Sprintf("    - %s: %s\n", addr.Type, addr.Address)
			}
		}

		// Show capacity
		if node.Node.Status.Capacity != nil {
			if cpu := node.Node.Status.Capacity["cpu"]; !cpu.IsZero() {
				output += fmt.Sprintf("  - CPU: %s\n", cpu.String())
			}
			if mem := node.Node.Status.Capacity["memory"]; !mem.IsZero() {
				output += fmt.Sprintf("  - Memory: %s\n", mem.String())
			}
		}
		output += "\n"
	}

	return output
}

func FormatPods(pods []PodInfo, namespace string) string {
	if len(pods) == 0 {
		return fmt.Sprintf("No pods found in namespace '%s'", namespace)
	}

	output := fmt.Sprintf("Found %d pod(s) in namespace '%s':\n\n", len(pods), namespace)
	for _, pod := range pods {
		output += fmt.Sprintf("**%s** (%s)\n", pod.Name, pod.Status.Phase)
		output += fmt.Sprintf("  - Node: %s\n", pod.Spec.NodeName)
		output += fmt.Sprintf("  - Ready: %s\n", pod.Ready)
		output += fmt.Sprintf("  - Restarts: %d\n", pod.Restarts)
		output += fmt.Sprintf("  - IP: %s\n", pod.Status.PodIP)
		output += fmt.Sprintf("  - Created: %s\n", pod.CreationTimestamp.Format("2006-01-02 15:04:05"))

		if len(pod.Containers) > 0 {
			output += "  - Containers:\n"
			for _, container := range pod.Containers {
				output += fmt.Sprintf("    - %s: %s (restarts: %d)\n",
					container.Name, container.State, container.RestartCount)
				if container.Reason != "" {
					output += fmt.Sprintf("      Reason: %s\n", container.Reason)
				}
			}
		}
		output += "\n"
	}

	return output
}

func FormatServices(services []ServiceInfo, namespace string) string {
	if len(services) == 0 {
		return fmt.Sprintf("No services found in namespace '%s'", namespace)
	}

	output := fmt.Sprintf("Found %d service(s) in namespace '%s':\n\n", len(services), namespace)
	for _, service := range services {
		svc := service.Spec
		output += fmt.Sprintf("**%s** (%s)\n", &service.Name, svc.Type)
		output += fmt.Sprintf("  - Cluster IP: %s\n", svc.ClusterIP)

		if len(svc.Ports) > 0 {
			output += "  - Ports:\n"
			for _, port := range svc.Ports {
				output += fmt.Sprintf("    - %d/%s", port.Port, port.Protocol)
				if port.TargetPort.IntVal != 0 {
					output += fmt.Sprintf(" -> %d", port.TargetPort.IntVal)
				} else if port.TargetPort.StrVal != "" {
					output += fmt.Sprintf(" -> %s", port.TargetPort.StrVal)
				}
				if port.Name != "" {
					output += fmt.Sprintf(" (%s)", port.Name)
				}
				output += "\n"
			}
		}

		if len(svc.Selector) > 0 {
			output += "  - Selector:\n"
			for k, v := range svc.Selector {
				output += fmt.Sprintf("    - %s: %s\n", k, v)
			}
		}

		output += fmt.Sprintf("  - Created: %s\n", service.CreationTimestamp.Format("2006-01-02 15:04:05"))
		output += "\n"
	}

	return output
}

func FormatDeployments(deployments []DeploymentInfo, namespace string) string {
	if len(deployments) == 0 {
		return fmt.Sprintf("No deployments found in namespace '%s'", namespace)
	}

	output := fmt.Sprintf("Found %d deployment(s) in namespace '%s':\n\n", len(deployments), namespace)
	for _, deploy := range deployments {
		output += fmt.Sprintf("**%s**\n", deploy.Name)
		output += fmt.Sprintf("  - Replicas: %d/%d ready\n", deploy.Status.ReadyReplicas, deploy.Status.Replicas)
		output += fmt.Sprintf("  - Updated: %d\n", deploy.Status.UpdatedReplicas)
		output += fmt.Sprintf("  - Available: %d\n", deploy.Status.AvailableReplicas)
		output += fmt.Sprintf("  - Strategy: %s\n", deploy.Spec.Strategy.Type)
		output += fmt.Sprintf("  - Created: %s\n", deploy.CreationTimestamp.Format("2006-01-02 15:04:05"))

		if len(deploy.Spec.Selector.MatchLabels) > 0 {
			output += "  - Selector:\n"
			for k, v := range deploy.Spec.Selector.MatchLabels {
				output += fmt.Sprintf("    - %s: %s\n", k, v)
			}
		}

		// Show conditions if there are any issues
		for _, condition := range deploy.Status.Conditions {
			if condition.Status != "True" {
				output += fmt.Sprintf("  - Issue: %s - %s\n", condition.Reason, condition.Message)
			}
		}

		output += "\n"
	}

	return output
}

// FormatAsJSON formats any data structure as pretty JSON
func FormatAsJSON(data interface{}) string {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting as JSON: %v", err)
	}
	return string(jsonData)
}
