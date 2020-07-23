package tracker

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestStarted(t *testing.T) {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	tracker := NewTracker(sugar)
	go tracker.StartTracker()
	time.Sleep(1000 * time.Millisecond)
	if tracker.numConnections == 0 {
		t.Errorf("tracker didn't start")
	}
}
