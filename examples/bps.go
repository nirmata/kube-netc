package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/nirmata/kube-netc/pkg/tracker"
	"strconv"
	"time"
)

func formatID(id tracker.ConnectionID) string {
	return id.DAddr + ":" + strconv.Itoa(int(id.DPort))
}

// This just shortens the humanize
func bf(b uint64) string {
	return humanize.Bytes(b)
}

func main() {
	t := tracker.NewTracker()
	go t.StartTracker()
	time.Sleep(5 * time.Second)
	conns := t.GetConnectionData()
	fmt.Printf("\t\t\tIn/s\tOut/s\tIn\tOut\tLast\n")
	for k, v := range conns {
		fmt.Printf("%s\t\t%s\t%s\t%s\t%s\t%s\n", formatID(k), bf(v.BytesRecvPerSecond), bf(v.BytesSentPerSecond), bf(v.BytesRecv), bf(v.BytesSent), v.LastUpdated)
	}
}
