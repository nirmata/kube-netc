package collector

import (
	"fmt"

	"github.com/nirmata/kube-netc/pkg/cluster"
	"github.com/nirmata/kube-netc/pkg/tracker"
	"github.com/prometheus/client_golang/prometheus"
)

func StartCollector(tr *tracker.Tracker, ci *cluster.ClusterInfo) {
	for {
		select {
		case update := <-tr.NodeUpdateChan:
			ActiveConnections.Set(float64(update.NumConnections))

		case update := <-tr.ConnUpdateChan:

			var labels prometheus.Labels
			conn := update.Connection

			if foundName, ok := ci.PodIPMap[update.Connection.DAddr]; ok {
				//The pod name is filled if there is a corresponding pod in the cluster
				labels = prometheus.Labels{
					"pod_name":            foundName.Name,
					"source_address":      tracker.IPPort(conn.SAddr, conn.SPort),
					"destination_address": tracker.IPPort(conn.DAddr, conn.DPort),
				}
			} else {
				labels = prometheus.Labels{
					"source_address":      tracker.IPPort(conn.SAddr, conn.SPort),
					"destination_address": tracker.IPPort(conn.DAddr, conn.DPort),
				}
			}

			BytesSent.With(labels).Set(float64(update.Data.BytesSent))
			BytesRecv.With(labels).Set(float64(update.Data.BytesRecv))
			BytesSentPerSecond.With(labels).Set(float64(update.Data.BytesSentPerSecond))
			BytesRecvPerSecond.With(labels).Set(float64(update.Data.BytesRecvPerSecond))

		case update := <-ci.MapUpdateChan:

			if update.Info != nil {
				fmt.Printf("Pod Update: %s -> %s\n", update.IP, update.Info.Name)

			} else {
				fmt.Printf("Pod Deleted: %s\n", update.IP)
			}

		}
	}
}
