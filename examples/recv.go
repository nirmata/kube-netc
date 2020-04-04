package main

import(
	"fmt"
	"time"
	"github.com/nirmata/kube-netsee/pkg/tracker"
	"github.com/dustin/go-humanize"
)



func main(){
	t := tracker.NewTracker()
	go t.StartTracker()
	ticker := time.NewTicker(3 * time.Second).C
	for{
		select{
		case <-ticker:
			fmt.Printf("%s bytes/s\n", humanize.Bytes(t.GetBytesRecvPerSecond()))
			fmt.Printf("%d connections\n", t.GetNumConnections())
			//fmt.Printf("%v\n", t.GetConnectionData())
		}
	}
}
