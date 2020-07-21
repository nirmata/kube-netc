package tracker

import (
	"errors"
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/ebpf"
	"go.uber.org/zap"
	"os"
	"time"
)

const (
	MaxConnBuffer = 256
)

func (t *Tracker) check(err error) {
	if err != nil {
		t.Logger.Fatalw(err.Error(),
			"package", "tracker",
		)
	}
}

type Tracker struct {
	Tick time.Duration
	// time idle before considering connection inactive
	Timeout        time.Duration
	Config         *ebpf.Config
	numConnections uint16

	ConnUpdateChan chan ConnUpdate
	NodeUpdateChan chan NodeUpdate

	Logger *zap.SugaredLogger

	stopChan chan struct{}
}

type ConnData struct {
	BytesSent          uint64
	BytesRecv          uint64
	BytesSentPerSecond uint64
	BytesRecvPerSecond uint64
	Active             bool
	LastUpdated        time.Time
}

// Simple struct to pipe for any extra metrics not related to connections
// but rather the node as a whole
type NodeUpdate struct {
	NumConnections uint16
}

// Type to be piped through chan to collector for updates
type ConnUpdate struct {
	Connection ConnectionID
	Data       ConnData
}

type ConnectionID struct {
	DAddr string
	DPort uint16
	SAddr string
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
			ProcRoot:                     "/proc",
			BPFDebug:                     false,
			EnableConntrack:              true,
			MaxClosedConnectionsBuffered: 50000,
			MaxConnectionsStateBuffered:  75000,
			ClientStateExpiry:            2 * time.Minute,
			ClosedChannelSize:            500,
		},
		numConnections: 0,
		ConnUpdateChan: make(chan ConnUpdate, MaxConnBuffer),
		NodeUpdateChan: make(chan NodeUpdate, 16),
		stopChan:       make(chan struct{}, 1),
	}
)

func NewTracker(logger *zap.SugaredLogger) *Tracker {
	dt := &DefaultTracker
	dt.Logger = logger
	return &DefaultTracker
}

func (t *Tracker) StartTracker() {
	err := checkSupport()
	t.check(err)
	t.Logger.Debugw("finished checking for eBPF suport",
		"package", "tracker",
	)
	err = t.run()
	t.check(err)
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

	ticker := time.NewTicker(t.Tick).C

	t.Logger.Debugw("starting tracker control loop",
		"package", "tracker",
	)

ControlLoop:
	for {

		select {

		case <-ticker:

			cs, err := tracer.GetActiveConnections(fmt.Sprintf("%d", os.Getpid()))

			if err != nil {
				return err
			}

			conns := cs.Conns

			t.NodeUpdateChan <- NodeUpdate{
				NumConnections: uint16(len(conns)),
			}

			t.Logger.Debugw("num connections calculated",
				"package", "tracker",
				"connections", uint16(len(conns)),
			)

			for _, c := range conns {
				id := ConnectionID{
					SAddr: c.Source.String(),
					DAddr: c.Dest.String(),
					DPort: c.DPort,
				}

				// These values get used mored than once in calculations
				// and we want them to be uniform in this scope
				bytesSent := c.MonotonicSentBytes
				bytesRecv := c.MonotonicRecvBytes

				// Using runtime.nanotime(), see util.go
				now := Now()
				// In float64 seconds
				timeDiff := float64(now-c.LastUpdateEpoch) / 1000000000.0

				if timeDiff <= 0 {
					t.check(errors.New("no difference between LastUpdateEpoch and time.Now(), will create divide by zero error or negative"))
				}

				// Per Second Calculations
				bytesSentPerSecond := uint64(float64(bytesSent-c.LastSentBytes) / float64(timeDiff))
				bytesRecvPerSecond := uint64(float64(bytesRecv-c.LastRecvBytes) / float64(timeDiff))
				if bytesSentPerSecond > 50e9 {
					t.Logger.Warnw("extreme transfer rate",
						"package", "tracker",
						"source", IPPort(c.Source.String(), c.SPort),
						"direction", "sent",
						"bps", bytesSentPerSecond,
					)
				} else if bytesRecvPerSecond > 50e9 {
					t.Logger.Warnw("extreme transfer rate",
						"package", "tracker",
						"source", IPPort(c.Source.String(), c.SPort),
						"direction", "recv",
						"bps", bytesRecvPerSecond,
					)
				}
				// Sending the updated stats through the pipe for collector to receive
				update := ConnUpdate{
					Connection: id,
					Data: ConnData{
						BytesSent:          c.MonotonicSentBytes,
						BytesRecv:          c.MonotonicRecvBytes,
						BytesSentPerSecond: bytesSentPerSecond,
						BytesRecvPerSecond: bytesRecvPerSecond,
						Active:             true,
						LastUpdated:        time.Now(),
					},
				}

				t.ConnUpdateChan <- update
			}

		case <-t.stopChan:
			break ControlLoop
		}
	}

	return nil
}

func (t *Tracker) Stop() {
	t.stopChan <- struct{}{}
}
