package main

import(
	"net/http"
	"fmt"
	
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/nirmata/kube-netc/pkg/collector"
	"github.com/nirmata/kube-netc/pkg/tracker"
)

func main(){
	t := tracker.NewTracker()
	go t.StartTracker()
	go collector.StartCollector(t)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("[SERVER STARTED ON :2112]")
	http.ListenAndServe(":2112", nil)
}
