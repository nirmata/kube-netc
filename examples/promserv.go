package main

import(
	"net/http"
	"time"
	
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/nirmata/kube-netsee/pkg/collector"
	"github.com/nirmata/kube-netsee/pkg/tracker"
)

func main(){
	t := tracker.NewTracker()
	go t.StartTracker()
	collector.StartCollector(t, 2 * time.Second)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
