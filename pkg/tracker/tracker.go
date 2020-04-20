package tracker

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	// Currently using the forked version
	"github.com/drewrip/datadog-agent/pkg/ebpf"
)

const (
	MaxConnBuffer = 256
)

func check(err error) {
	if err != nil {
		log.Fatalf("[%v] error: %s", time.Now(), err)
	}
}

type Tracker struct {
	Tick time.Duration
	// time idle before considering connection inactive
	Timeout        time.Duration
	Config         *ebpf.Config
	numConnections uint16

	tracer *ebpf.Tracer
	
	// These are the totals
	bytesSent uint64
	bytesRecv uint64
	totalSent uint64
	totalRecv uint64

	bytesSentPerSecond uint64
	bytesRecvPerSecond uint64

	// For use in calculation of bytesPerSecond
	trackLastUpdated time.Time

	// string key will be in the form ip:port
	dataHistory map[ConnectionID]*trackData

	NodeUpdateChan chan NodeUpdate
	ConnUpdateChan chan ConnUpdate
	stopChan chan struct{}
}

// Stats that are tracked for each connection
type trackData struct {
	// is this connection still active
	active bool
	// these are totals for each connection
	bytesSent uint64
	bytesRecv uint64
	// Used in calculation of bytesPerSecond
	lastBytesSent uint64
	lastBytesRecv uint64

	bytesSentPerSecond uint64
	bytesRecvPerSecond uint64

	// Used to tell when connection is inactive
	lastUpdated time.Time
}

type ExportData struct {
	BytesSent          uint64
	BytesRecv          uint64
	BytesSentPerSecond uint64
	BytesRecvPerSecond uint64
	LastUpdated        time.Time
}

type NodeUpdate struct {
	BytesSent          uint64
	BytesRecv          uint64
	BytesSentPerSecond uint64
	BytesRecvPerSecond uint64
	LastUpdated        time.Time

	NumConnections uint16
}

// Type to be piped through chan to collector for updates
type ConnUpdate struct {
	Connection ConnectionID
	Data       ExportData
}

type ConnectionID struct {
	DAddr string
	DPort uint16
	SAddr string
	SPort uint16
}

var (
	DefaultTracker = Tracker{
		Tick: 1 * time.Second,
		Config: &ebpf.Config{
			CollectTCPConns:              true,
			CollectUDPConns:              true,
			CollectIPv6Conns:             true,
			CollectLocalDNS:              false,
			DNSInspection:                false,
			UDPConnTimeout:               30 * time.Second,
			TCPConnTimeout:               2 * time.Minute,
			MaxTrackedConnections:        65536,
			ConntrackMaxStateSize:        65536,
			ConntrackShortTermBufferSize: 100,
			ProcRoot:                     "/proc",
			BPFDebug:                     false,
			EnableConntrack:              true,
			MaxClosedConnectionsBuffered: 50000,
			MaxConnectionsStateBuffered:  75000,
			ClientStateExpiry:            2 * time.Minute,
			ClosedChannelSize:            500,
		},
		numConnections:     0,
		bytesSent:          0,
		bytesRecv:          0,
		bytesSentPerSecond: 0.0,
		bytesRecvPerSecond: 0.0,
		dataHistory:        make(map[ConnectionID]*trackData),
		NodeUpdateChan:     make(chan NodeUpdate, MaxConnBuffer),
		ConnUpdateChan:     make(chan ConnUpdate, MaxConnBuffer),
		stopChan:           make(chan struct{}, 1),
	}
)

func NewTracker() *Tracker {
	return &DefaultTracker
}

func (t *Tracker) StartTracker() {
	err := checkSupport()
	check(err)
	err = t.run()
	check(err)
}

func checkSupport() error {
	_, err := ebpf.CurrentKernelVersion()
	if err != nil {
		return err
	}

	if supported, errtip := ebpf.IsTracerSupportedByOS(nil); !supported {
		return errors.New(errtip)
	}

	return nil
}

func (t *Tracker) run() error {
	tracer, err := ebpf.NewTracer(t.Config)
	if err != nil {
		return err
	}

	// Initial set
	t.trackLastUpdated = time.Now()

	ticker := time.NewTicker(t.Tick).C

ControlLoop:
	for {
		select{

			case <-ticker:
		for k, v := range t.dataHistory {
			if time.Since(v.lastUpdated) >= 20*time.Second {
				t.dataHistory[k].active = false
			}
		}

		cs, err := tracer.GetActiveConnections(fmt.Sprintf("%d", os.Getpid()))
		if err != nil {
			return err
		}

		conns := cs.Conns
		for _, c := range conns {
			id := ConnectionID{
				SAddr: c.Source.String(),
				SPort: c.SPort,
				DAddr: c.Dest.String(),
				DPort: c.DPort,
			}
			// Creating a new entry for this connection if it doesn't exist
			if _, ok := t.dataHistory[id]; !ok {
				t.dataHistory[id] = &trackData{
					bytesSent:     c.MonotonicSentBytes,
					bytesRecv:     c.MonotonicRecvBytes,
					lastBytesSent: c.MonotonicSentBytes,
					lastBytesRecv: c.MonotonicRecvBytes,
					active:        true,
					lastUpdated:   time.Now(),
				}
			}
			// Updating the entry if it does exist
			lastbSent := t.dataHistory[id].bytesSent
			lastbRecv := t.dataHistory[id].bytesRecv
			currbSent := c.MonotonicSentBytes
			currbRecv := c.MonotonicRecvBytes

			t.dataHistory[id].lastBytesSent = lastbSent
			t.dataHistory[id].lastBytesRecv = lastbRecv
			t.dataHistory[id].bytesSent = currbSent
			t.dataHistory[id].bytesRecv = currbRecv

			tConn := time.Since(t.dataHistory[id].lastUpdated)

			bRPS := uint64(float64(currbRecv-lastbRecv) / (float64(tConn) / float64(time.Second)))
			bSPS := uint64(float64(currbSent-lastbSent) / (float64(tConn) / float64(time.Second)))
			t.dataHistory[id].bytesRecvPerSecond = bRPS
			t.dataHistory[id].bytesSentPerSecond = bSPS

			lastUpdated := time.Now()

			t.dataHistory[id].lastUpdated = lastUpdated

			// Sending the updated stats through the pipe for collector to receive
			update := ConnUpdate{
				Connection: id,
				Data: ExportData{
					BytesSent:          currbSent,
					BytesRecv:          currbRecv,
					BytesSentPerSecond: bSPS,
					BytesRecvPerSecond: bRPS,
					LastUpdated:        lastUpdated,
				},
			}

			t.ConnUpdateChan <- update
		}

		// Adding the new bytes to the stats
		var (
			newSentBytes       uint64 = 0
			newRecvBytes       uint64 = 0
			newChangeSentBytes uint64 = 0
			newChangeRecvBytes uint64 = 0
			totalSent          uint64 = 0
			totalRecv          uint64 = 0
			numConnections     uint16 = 0
		)

		for _, v := range t.dataHistory {
			if v.active {
				numConnections++
				newSentBytes += v.bytesSent
				newRecvBytes += v.bytesRecv
				newChangeSentBytes += (v.bytesSent - v.lastBytesSent)
				newChangeRecvBytes += (v.bytesRecv - v.lastBytesRecv)
			}
			totalSent += v.bytesSent
			totalRecv += v.bytesRecv
		}

		t.numConnections = numConnections

		// t.totalSent/Recv is historical, counts bytes for connections that are no longer active
		t.totalSent = totalSent
		t.totalRecv = totalRecv
		t.bytesSent = newSentBytes
		t.bytesRecv = newRecvBytes
		tStop := time.Since(t.trackLastUpdated)
		// bytes per second for the sent bytes
		tBSPS := uint64(float64(newChangeSentBytes) / (float64(tStop) / float64(time.Second)))
		// bytes per second for the receive bytes
		tBRPS := uint64(float64(newChangeRecvBytes) / (float64(tStop) / float64(time.Second)))

		t.bytesSentPerSecond = tBSPS
		t.bytesRecvPerSecond = tBRPS

		tLastUpdated := time.Now()

		t.trackLastUpdated = tLastUpdated

		trackUpdate := NodeUpdate{
			BytesSent:          totalSent,
			BytesRecv:          totalRecv,
			BytesSentPerSecond: tBSPS,
			BytesRecvPerSecond: tBRPS,
			LastUpdated:        tLastUpdated,
			NumConnections:     numConnections,
		}

		t.NodeUpdateChan <- trackUpdate
			case <-t.stopChan:
			break ControlLoop
		}
	}
	return nil
}

func (t *Tracker) Stop() {
	t.stopChan<-struct{}{}
}

// Clears the current internal tracking data.
func (t *Tracker) ResetStats() error {
	t.dataHistory = make(map[ConnectionID]*trackData)
	return nil
}
