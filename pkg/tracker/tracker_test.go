package tracker

import (
	"testing"
	"time"
)

func TestStarted(t *testing.T) {
	tracker := NewTracker()
	go tracker.StartTracker()
	time.Sleep(1000 * time.Millisecond)
	if tracker.numConnections == 0 {
		t.Errorf("tracker didn't start")
	}
}
