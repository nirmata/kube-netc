package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nirmata/kube-netc/pkg/cluster"
	"github.com/nirmata/kube-netc/pkg/collector"
	"github.com/nirmata/kube-netc/pkg/tracker"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func check(err error) {
	if err != nil {
		log.Fatalf("[ERR] %s", err)
	}
}

func main() {

	t := tracker.NewTracker()
	go t.StartTracker()

	clusterInfo := cluster.NewClusterInfo()
	go clusterInfo.Run()

	go collector.StartCollector(t, clusterInfo)

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("[SERVER STARTED ON :9655]")
	err := http.ListenAndServe(":9655", nil)
	check(err)
}
