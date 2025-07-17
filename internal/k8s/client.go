package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	clientset *kubernetes.Clientset
}

func NewClient() (*Client, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &Client{
		clientset: clientset,
	}, nil
}

func getKubeConfig() (*rest.Config, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	if _, err := os.Stat(kubeconfig); err == nil {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	return rest.InClusterConfig()
}

func (c *Client) GetNodes(ctx context.Context) ([]NodeInfo, error) {
	nodes, err := c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	result := make([]NodeInfo, len(nodes.Items))
	for i, node := range nodes.Items {
		result[i] = NodeInfo{
			Node:   node,
			Status: getNodeStatus(node),
		}
	}

	return result, nil
}

func (c *Client) GetPods(ctx context.Context, namespace string) ([]PodInfo, error) {
	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	result := make([]PodInfo, len(pods.Items))
	for i, pod := range pods.Items {
		containers := make([]ContainerInfo, len(pod.Status.ContainerStatuses))
		for j, container := range pod.Status.ContainerStatuses {
			state, reason, message := getContainerState(container)
			containers[j] = ContainerInfo{
				Name:         container.Name,
				Image:        container.Image,
				Ready:        container.Ready,
				RestartCount: container.RestartCount,
				State:        state,
				Reason:       reason,
				Message:      message,
			}
		}

		result[i] = PodInfo{
			Pod:        pod,
			Ready:      getPodReadyStatus(pod),
			Restarts:   getPodRestartCount(pod),
			Containers: containers,
		}
	}

	return result, nil
}

func (c *Client) GetServices(ctx context.Context, namespace string) ([]ServiceInfo, error) {
	services, err := c.clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	result := make([]ServiceInfo, len(services.Items))
	for i, svc := range services.Items {
		result[i] = ServiceInfo{
			Service: svc,
		}
	}

	return result, nil
}

func (c *Client) GetDeployments(ctx context.Context, namespace string) ([]DeploymentInfo, error) {
	deployments, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %w", err)
	}

	result := make([]DeploymentInfo, len(deployments.Items))
	for i, deploy := range deployments.Items {
		result[i] = DeploymentInfo{
			Deployment: deploy,
		}
	}

	return result, nil
}

func getNodeStatus(node corev1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			if condition.Status == corev1.ConditionTrue {
				return "Ready"
			}
			return "NotReady"
		}
	}
	return "Unknown"
}

func getPodReadyStatus(pod corev1.Pod) string {
	ready := 0
	total := len(pod.Status.ContainerStatuses)
	for _, status := range pod.Status.ContainerStatuses {
		if status.Ready {
			ready++
		}
	}
	return fmt.Sprintf("%d/%d", ready, total)
}

func getPodRestartCount(pod corev1.Pod) int32 {
	var total int32
	for _, status := range pod.Status.ContainerStatuses {
		total += status.RestartCount
	}
	return total
}

func getContainerState(container corev1.ContainerStatus) (state, reason, message string) {
	if container.State.Running != nil {
		return "Running", "", ""
	}
	if container.State.Waiting != nil {
		return "Waiting", container.State.Waiting.Reason, container.State.Waiting.Message
	}
	if container.State.Terminated != nil {
		return "Terminated", container.State.Terminated.Reason, container.State.Terminated.Message
	}
	return "Unknown", "", ""
}
