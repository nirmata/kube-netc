package collector

import (
	"github.com/nirmata/kube-netc/pkg/cluster"
	"github.com/nirmata/kube-netc/pkg/tracker"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	emptyInfo = &cluster.ObjectInfo{
		Name:      "",
		Kind:      "",
		Namespace: "",
		Node:      "",

		LabelName:      "",
		LabelComponent: "",
		LabelInstance:  "",
		LabelVersion:   "",
		LabelPartOf:    "",
		LabelManagedBy: "",
	}
)

func generateLabels(connup tracker.ConnUpdate, ci *cluster.ClusterInfo) prometheus.Labels {

	conn := connup.Connection

	srcInfo, sok := ci.Get(conn.SAddr)

	if !sok {
		srcInfo = emptyInfo
	}
	
	destInfo, dok := ci.Get(conn.DAddr)

	if !dok {
		destInfo = emptyInfo
	}

	return prometheus.Labels{
		// Kubernetes labels
		"name":       srcInfo.LabelName,
		"component":  srcInfo.LabelComponent,
		"instance":   srcInfo.LabelInstance,
		"version":    srcInfo.LabelVersion,
		"part_of":    srcInfo.LabelPartOf,
		"managed_by": srcInfo.LabelManagedBy,
		// Nirmata networking labels
		"source_address":        tracker.IPPort(conn.SAddr, conn.SPort),
		"destination_address":   tracker.IPPort(conn.DAddr, conn.DPort),
		"source_name":           srcInfo.Name,
		"destination_name":      destInfo.Name,
		"source_kind":           srcInfo.Kind,
		"destination_kind":      destInfo.Kind,
		"source_namespace":      srcInfo.Namespace,
		"destination_namespace": destInfo.Namespace,
		"source_node":           srcInfo.Node,
		"destination_node":      destInfo.Node,
	}
}

func StartCollector(tr *tracker.Tracker, ci *cluster.ClusterInfo) {
	for {
		select {
		case update := <-tr.NodeUpdateChan:
			ActiveConnections.Set(float64(update.NumConnections))

		case update := <-tr.ConnUpdateChan:

			labels := generateLabels(update, ci)
			BytesSent.With(labels).Set(float64(update.Data.BytesSent))
			BytesRecv.With(labels).Set(float64(update.Data.BytesRecv))
			BytesSentPerSecond.With(labels).Set(float64(update.Data.BytesSentPerSecond))
			BytesRecvPerSecond.With(labels).Set(float64(update.Data.BytesRecvPerSecond))

		}

	}
}
