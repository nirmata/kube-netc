package main

import (
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/ebpf"
	"os"
	"time"
	_ "unsafe"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var config *ebpf.Config = &ebpf.Config{
	CollectTCPConns:              true,
	CollectUDPConns:              true,
	CollectIPv6Conns:             true,
	CollectLocalDNS:              false,
	DNSInspection:                false,
	UDPConnTimeout:               30 * time.Second,
	TCPConnTimeout:               2 * time.Minute,
	MaxTrackedConnections:        65536,
	ConntrackMaxStateSize:        65536,
	ProcRoot:                     "/proc",
	BPFDebug:                     false,
	EnableConntrack:              true,
	MaxClosedConnectionsBuffered: 50000,
	MaxConnectionsStateBuffered:  75000,
	ClientStateExpiry:            2 * time.Minute,
	ClosedChannelSize:            500,
}

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64

// time.Now() is 45ns, runtime.nanotime is 20ns
// I can not create an exported symbol with //go:linkname
// I need a wrapper
// Go does not inline functions? https://lemire.me/blog/2017/09/05/go-does-not-inline-functions-when-it-should/
// The wrapper costs 5ns per call
func Nanotime() int64 {
	return nanotime()
}

func checkSupport() {
	_, err := ebpf.CurrentKernelVersion()
	check(err)

	if supported, errtip := ebpf.IsTracerSupportedByOS(nil); !supported {
		panic(errtip)
	}
}

func main() {
	start := Nanotime()
	fmt.Println(start)
	tracer, err := ebpf.NewTracer(config)
	check(err)
	fmt.Printf("Add\t\tCurr\t\tLast\t\tLastUpdate\n")
	tick := time.NewTicker(1 * time.Second).C
	for {
		cs, err := tracer.GetActiveConnections(fmt.Sprintf("%d", os.Getpid()))
		check(err)
		conns := cs.Conns
		select {
		case <-tick:
			for _, c := range conns {
				fmt.Printf("%v\t\t%d\t\t%d\t\t%d\n", c.Dest, c.MonotonicSentBytes, c.LastSentBytes, c.LastUpdateEpoch)
			}
		}
	}
}
