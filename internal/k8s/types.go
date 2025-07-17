package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type NodeInfo struct {
	corev1.Node `json:"node"`
	Status      string `json:"status"`
}

type PodInfo struct {
	corev1.Pod `json:"pod"`
	Ready      string          `json:"ready"`
	Restarts   int32           `json:"restarts"`
	Containers []ContainerInfo `json:"containers"`
}

type ContainerInfo struct {
	Name         string `json:"name"`
	Image        string `json:"image"`
	Ready        bool   `json:"ready"`
	RestartCount int32  `json:"restartCount"`
	State        string `json:"state"`
	Reason       string `json:"reason,omitempty"`
	Message      string `json:"message,omitempty"`
}

type ServiceInfo struct {
	corev1.Service `json:"service"`
}

type DeploymentInfo struct {
	appsv1.Deployment `json:"deployment"`
}

type RequestParams struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}
