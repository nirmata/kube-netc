package main

import (
	"fmt"
	"net/http"
	"log"
	
	"github.com/nirmata/kube-netc/pkg/collector"
	"github.com/nirmata/kube-netc/pkg/tracker"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func check(err error){
	if err != nil {
		log.Fatalf("[ERR] %s", err)
	}
}

func main() {
	t := tracker.NewTracker()
	go t.StartTracker()
	go collector.StartCollector(t)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("[SERVER STARTED ON :2112]")
	err := http.ListenAndServe(":2112", nil)
	check(err)
}
