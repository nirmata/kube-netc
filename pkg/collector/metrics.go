package collector

// func init() lives in this file

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	BytesSent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bytes_sent",
			Help: "Total bytes sent to a given connection",
		},
		[]string{"pod_name", "pod_address"},
	)

	BytesRecv = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bytes_recv",
			Help: "Total bytes received from a given connection",
		},
		[]string{"pod_name", "pod_address"},
	)

	BytesSentPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bytes_sent_per_second",
			Help: "Bytes per second being sent to a given connection",
		},
		[]string{"pod_name", "pod_address"},
	)

	BytesRecvPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bytes_recv_per_second",
			Help: "Bytes per second being received from a given connection",
		},
		[]string{"pod_name", "pod_address"},
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
