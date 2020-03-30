package main

import(
	"net/http"
	"time"
	
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/nirmata/kube-netsee/collector"
	"github.com/nirmata/kube-netsee/tracker"
)

func main(){
	t := tracker.NewTracker()
	go t.StartTracker()
	collector.StartCollector(t, 2 * time.Second)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
