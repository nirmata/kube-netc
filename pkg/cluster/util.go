package cluster

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

func getObjectType(o interface{}) string {
	switch v := o.(type) {
	case *corev1.Pod:
		return "pod"
	case *corev1.Service:
		return "service"
	case *corev1.Node:
		return "node"
	default:
		_ = v
		return "unknown"
	}
}

func getNodeIP(n *corev1.Node) (string, error) {

	// Of type []NodeAddress
	addrs := n.Status.Addresses

	for _, a := range addrs {
		if string(a.Type) == "InternalIP" {
			return a.Address, nil
		}
	}

	return "", fmt.Errorf("Could not find internal IP of Node: %s", n.GetName())
}
