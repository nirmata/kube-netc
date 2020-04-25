package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/nirmata/kube-netc/pkg/tracker"
	"time"
)

func main() {
	t := tracker.NewTracker()
	go t.StartTracker()
	ticker := time.NewTicker(3 * time.Second).C
	for {
		select {
		case <-ticker:
			fmt.Printf("%s bytes/s\n", humanize.Bytes(t.GetBytesRecvPerSecond()))
			fmt.Printf("%d connections\n", t.GetNumConnections())
			//fmt.Printf("%v\n", t.GetConnectionData())
		}
	}
}
