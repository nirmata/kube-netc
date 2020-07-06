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
	fmt.Printf("\t\t\tIn/s\tOut/s\tIn\tOut\tLast\n")
	for {
		select {
		case u := <-t.ConnUpdateChan:
			fmt.Printf("%s\t\t%s\t%s\t%s\t%s\t%s\n", formatID(u.Connection), bf(u.Data.BytesRecvPerSecond), bf(u.Data.BytesSentPerSecond), bf(u.Data.BytesRecv), bf(u.Data.BytesSent), u.Data.LastUpdated)
		}
	}
}
