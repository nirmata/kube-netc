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

			if foundName, ok := ci.PodIPMap[update.Connection.DAddr]; ok {
				fmt.Println("Found pod name from IP")
				labels = prometheus.Labels{
					"pod_name":    foundName.Name,
					"pod_address": update.Connection.DAddr,
				}
			} else {
				labels = prometheus.Labels{
					"pod_name":    "NOT_FOUND",
					"pod_address": update.Connection.DAddr,
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
