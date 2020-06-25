package collector

// func init() lives in this file

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	labs = []string{
		"name",
		"component",
		"instance",
		"version",
		"part_of",
		"managed_by",
		"source_address",
		"destination_address",
		"source_name",
		"destination_name",
		"source_kind",
		"destination_kind",
		"source_namespace",
		"destination_namespace",
		"source_node",
		"destination_node",
	}
)

var (
	BytesSent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bytes_sent",
			Help: "Total bytes sent to a given connection",
		},
		labs,
	)

	BytesRecv = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bytes_recv",
			Help: "Total bytes received from a given connection",
		},
		labs,
	)

	BytesSentPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bytes_sent_per_second",
			Help: "Bytes per second being sent to a given connection",
		},
		labs,
	)

	BytesRecvPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bytes_recv_per_second",
			Help: "Bytes per second being received from a given connection",
		},
		labs,
	)

	// Number of connections
	ActiveConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "active_connections",
		Help: "Number of connections to the node",
	})
)

func init() {
	prometheus.MustRegister(BytesSent)
	prometheus.MustRegister(BytesRecv)
	prometheus.MustRegister(BytesRecvPerSecond)
	prometheus.MustRegister(BytesSentPerSecond)
	prometheus.MustRegister(ActiveConnections)
}
