package tracker


import(
	"time"
	"os"
	"fmt"
	"log"
	"errors"
	// Currently using the forked version
	"github.com/drewrip/datadog-agent/pkg/ebpf"
)

func check(err error){
	if err != nil{
		log.Fatalf("[%v] error: %s", time.Now(), err)
	}
}

type Tracker struct {
	Tick time.Duration
	// time idle before considering connection inactive
	Timeout time.Duration
	Config *ebpf.Config
	numConnections uint16
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
	dataHistory map[string]*trackData
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
	// Used to tell when connection is inactive
	lastUpdated time.Time
}

var(
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
		numConnections: 0,
		bytesSent: 0,
		bytesRecv: 0,
		bytesSentPerSecond: 0.0,
		bytesRecvPerSecond: 0.0,
		dataHistory: make(map[string]*trackData),
		
	}
)

func NewTracker() *Tracker {
	return &DefaultTracker
}

func (t *Tracker) StartTracker() {
	err := checkSupport()
	check(err)
	fmt.Println("Running tracker...")
	err = t.run()
	check(err)
}

func (t *Tracker) GetBytesRecvPerSecond() uint64 {
	return t.bytesRecvPerSecond
}

func (t *Tracker) GetBytesRecv() uint64 {
	return t.bytesRecv
}

func (t *Tracker) GetTotalBytesRecv() uint64 {
	return t.totalRecv
}

func checkSupport() error {
	_, err := ebpf.CurrentKernelVersion()
	if err != nil{
		return err
	}

	if supported, errtip := ebpf.IsTracerSupportedByOS(nil); !supported{
		return errors.New(errtip)
	}

	return nil
}

func (t *Tracker) run() error {
	fmt.Println("Tracker Running")
	tracer, err := ebpf.NewTracer(t.Config)
	if err != nil{
		return err
	}

	// Initial set
	t.trackLastUpdated = time.Now()
	
	ticker := time.NewTicker(t.Tick).C
	for{
		select{
		case <-ticker:

			for k, v := range t.dataHistory {
				if time.Since(v.lastUpdated) >= 5 * time.Second {
					t.dataHistory[k].active = false
				}
			}
			
			cs, err := tracer.GetActiveConnections(fmt.Sprintf("%d", os.Getpid()))
			if err != nil{
				return err
			}

			conns := cs.Conns
			for _, c := range conns{
				id := fmtAddress(c.String(), c.DPort)
				if _, ok := t.dataHistory[id]; !ok{
					t.dataHistory[id] = &trackData{
						bytesSent: c.MonotonicSentBytes,
						bytesRecv: c.MonotonicRecvBytes,
						active: true,
						lastUpdated: time.Now(),
					}	
				} else {
					t.dataHistory[id].lastBytesSent = t.dataHistory[id].bytesSent
					t.dataHistory[id].lastBytesRecv = t.dataHistory[id].bytesRecv
					t.dataHistory[id].bytesSent = c.MonotonicSentBytes
					t.dataHistory[id].bytesRecv = c.MonotonicRecvBytes
					t.dataHistory[id].active = true
					t.dataHistory[id].lastUpdated = time.Now()
					
				}

			}

			// Adding the new bytes to the stats
			var newSentBytes uint64
			var newRecvBytes uint64
			var newChangeSentBytes uint64
			var newChangeRecvBytes uint64
			var totalSent uint64
			var totalRecv uint64
			
			for _, v := range t.dataHistory {
				if v.active {
					newSentBytes += v.bytesSent
					newRecvBytes += v.bytesRecv
					newChangeSentBytes += (v.bytesSent - v.lastBytesSent)
					newChangeRecvBytes += (v.bytesRecv - v.lastBytesRecv)
				}
				totalSent += v.bytesSent
				totalRecv += v.bytesRecv
			}

			t.totalSent = totalSent
			t.totalRecv = totalRecv
			t.bytesSent = newSentBytes
			t.bytesRecv = newRecvBytes
			t.bytesSentPerSecond = uint64(float64(newChangeSentBytes)/float64(time.Since(t.trackLastUpdated)/time.Second))
			t.bytesRecvPerSecond = uint64(float64(newChangeRecvBytes)/float64(time.Since(t.trackLastUpdated)/time.Second))



			t.trackLastUpdated = time.Now()

		}
	}
}
func fmtAddress(addr string, port uint16) string {
	return addr + ":" + string(port)
}

// Clears the current internal tracking data.
func (t *Tracker) ResetStats() error {
	t.dataHistory = make(map[string]*trackData)
	return nil
}
