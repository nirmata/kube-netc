package collector

import(
	"time"
	
	"github.com/nirmata/kube-netsee/pkg/tracker"
	"github.com/prometheus/client_golang/prometheus"
)

func StartCollector(tr *tracker.Tracker, t time.Duration){
	go func(tick time.Duration){
		tc := time.NewTicker(t).C
		for{
			select{
				case <-tc:
				
				ActiveConnections.Set(float64(tr.GetNumConnections()))
				
				conns := tr.GetConnectionData()
				for k, v := range conns {
					dest := prometheus.Labels{"dest": tracker.FormatCID(k)}
					BytesSent.With(dest).Set(float64(v.BytesSent))
					BytesRecv.With(dest).Set(float64(v.BytesRecv))
					BytesSentPerSecond.With(dest).Set(float64(v.BytesSentPerSecond))
					BytesRecvPerSecond.With(dest).Set(float64(v.BytesRecvPerSecond))
				}
			}
		}
	}(t)
}
