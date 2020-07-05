package cluster

import (
	"k8s.io/api/core/v1"
)

func (c *ClusterInfo) handleNewObject(obj interface{}) {

	var name string
	var kind string
	var namespace string
	var node string

	var ip string
	var labels map[string]string

	// Copies all of the attributes as the right type
	switch o := obj.(type) {
	case *v1.Pod:
		ip = o.Status.PodIP
		node = o.Spec.NodeName

		name = o.GetName()
		kind = getObjectType(o)
		namespace = o.GetNamespace()

		labels = o.GetLabels()
	case *v1.Service:
		ip = o.Spec.ClusterIP

		name = o.GetName()
		kind = getObjectType(o)
		namespace = o.GetNamespace()

		labels = o.GetLabels()
	case *v1.Node:
		internalIP, err := getNodeIP(o)
		check(err)
		ip = internalIP

		name = o.GetName()
		kind = getObjectType(o)
		namespace = o.GetNamespace()

		labels = o.GetLabels()
	}

	info := &ObjectInfo{
		Name:      name,
		Kind:      kind,
		Namespace: namespace,
		Node:      node,
	}

	info.LabelName = labels["name"]
	info.LabelComponent = labels["component"]
	info.LabelInstance = labels["instance"]
	info.LabelVersion = labels["version"]
	info.LabelPartOf = labels["part-of"]
	info.LabelManagedBy = labels["managed-by"]

	// Updating the map
	c.Set(ip, info)
}

func (c *ClusterInfo) handleUpdateObject(oldObj interface{}, obj interface{}) {
	_ = oldObj

	var name string
	var kind string
	var namespace string
	var node string

	var ip string
	var labels map[string]string

	// Copies all of the attributes as the right type
	switch o := obj.(type) {
	case *v1.Pod:
		ip = o.Status.PodIP
		node = o.Spec.NodeName

		name = o.GetName()
		kind = getObjectType(o)
		namespace = o.GetNamespace()

		labels = o.GetLabels()
	case *v1.Service:
		ip = o.Spec.ClusterIP

		name = o.GetName()
		kind = getObjectType(o)
		namespace = o.GetNamespace()

		labels = o.GetLabels()
	case *v1.Node:
		internalIP, err := getNodeIP(o)
		check(err)
		ip = internalIP

		name = o.GetName()
		kind = getObjectType(o)
		namespace = o.GetNamespace()

		labels = o.GetLabels()
	}

	info := &ObjectInfo{
		Name:      name,
		Kind:      kind,
		Namespace: namespace,
		Node:      node,
	}

	info.LabelName = labels["name"]
	info.LabelComponent = labels["component"]
	info.LabelInstance = labels["instance"]
	info.LabelVersion = labels["version"]
	info.LabelPartOf = labels["part-of"]
	info.LabelManagedBy = labels["managed-by"]

	// Updating the map
	c.Set(ip, info)
}

func (c *ClusterInfo) handleDeleteObject(obj interface{}) {

	var ip string

	switch o := obj.(type) {
	case *v1.Pod:
		ip = o.Status.PodIP
	case *v1.Service:
		ip = o.Spec.ClusterIP
	case *v1.Node:
		internalIP, err := getNodeIP(o)
		check(err)
		ip = internalIP
	}

	// Updating the map
	c.Set(ip, nil)
}
