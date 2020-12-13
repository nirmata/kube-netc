package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/nirmata/kube-netc/pkg/tracker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
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

func newZapLogger() *zap.SugaredLogger {
	atom := zap.NewAtomicLevel()
	atom.SetLevel(zapcore.InfoLevel)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	return zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), atom)).Sugar()
}

func main() {
	logger := newZapLogger()
	t := tracker.NewTracker(logger)
	go t.StartTracker()
	time.Sleep(5 * time.Second)
	fmt.Printf("\t\t\tIn/s\tOut/s\tIn\tOut\tLast\n")
	for {
		select {
		case u := <-t.ConnUpdateChan:
			fmt.Printf("%s\t\t%s\t%s\t%s\t%s\t%s\n", formatID(u.Connection), bf(u.Data.BytesRecvPerSecond), bf(u.Data.BytesSentPerSecond), bf(u.Data.BytesRecv), bf(u.Data.BytesSent), u.Data.LastUpdated.Format("2006-01-02 15:04:05"))
		}
	}
}
