package main

import(
	"fmt"
	"time"
	"github.com/nirmata/kube-netsee/tracker"
	"github.com/dustin/go-humanize"
)



func main(){
	t := tracker.NewTracker()
	go t.StartTracker()
	ticker := time.NewTicker(2 * time.Second).C
	for{
		select{
		case <-ticker:
			fmt.Printf("%s bytes recv\n", humanize.Bytes(t.GetTotalBytesRecv()))
			fmt.Printf("%d connections\n", t.GetNumConnections())
		}
	}
}
