package collector

import(
	"time"

	"github.com/nirmata/kube-netsee/tracker"
	"github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
)


var(
	activeConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_connections",
		Help: "Number of connections to the node",
	})
)

func StartCollector(tracker *tracker.Tracker, t time.Duration){
	go func(tick time.Duration){
		tc := time.NewTicker(t).C
		for{
			select{
				case <-tc:
				activeConnections.Set(float64(tracker.GetNumConnections()))
			}
		}
	}(t)
}
