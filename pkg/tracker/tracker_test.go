package tracker

import(
	"testing"
	"time"
)



func TestStarted(t *testing.T){
	tracker := NewTracker()
	go tracker.StartTracker()
	time.Sleep(500 * time.Millisecond)
	if tracker.trackLastUpdated.UnixNano() == 0 {
		t.Errorf("tracker didn't start")
	}
}
