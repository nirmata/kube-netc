package collector

import(
	"time"
	
	"github.com/nirmata/kube-netsee/pkg/tracker"
	"github.com/prometheus/client_golang/prometheus"
)

func StartCollector(tr *tracker.Tracker, t time.Duration){
	for{
		select{
			case update:=<-tr.NodeUpdateChan:

			ActiveConnections.Set(float64(update.NumConnections))

			case update:=<-tr.ConnUpdateChan:

			id := prometheus.Labels{"id": tracker.FormatCID(update.Connection)}
			BytesSent.With(id).Set(float64(update.Data.BytesSent))
			BytesRecv.With(id).Set(float64(update.Data.BytesRecv))
			BytesSentPerSecond.With(id).Set(float64(update.Data.BytesSentPerSecond))
			BytesRecvPerSecond.With(id).Set(float64(update.Data.BytesRecvPerSecond))
			
		}
	}
}
