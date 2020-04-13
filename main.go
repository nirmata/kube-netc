package main

import (
	"fmt"
	"net/http"

	"github.com/nirmata/kube-netc/pkg/collector"
	"github.com/nirmata/kube-netc/pkg/tracker"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	t := tracker.NewTracker()
	go t.StartTracker()
	go collector.StartCollector(t)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("[SERVER STARTED ON :2112]")
	http.ListenAndServe(":2112", nil)
}
