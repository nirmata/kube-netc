package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/nirmata/kube-netc/pkg/cluster"
	"github.com/nirmata/kube-netc/pkg/collector"
	"github.com/nirmata/kube-netc/pkg/tracker"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func check(err error) {
	if err != nil {
		log.Printf("[ERR] %s", err)
	}
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "debug":
		return zapcore.DebugLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		log.Printf("[WARN] Unknown logging level %s: using INFO as default", level)
		return zapcore.InfoLevel
	}
}

func newZapLogger(level string) *zap.SugaredLogger {
	atom := zap.NewAtomicLevel()
	atom.SetLevel(getZapLevel(level))
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	return zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), atom)).Sugar()
}

func main() {
	// command line flags
	var logLevel string
	flag.StringVar(&logLevel, "v", "info", "the log level")

	flag.Parse()

	logger := newZapLogger(logLevel)
	//nolint:errcheck
	defer logger.Sync()

	t := tracker.NewTracker(logger)
	go t.StartTracker()

	clusterInfo := cluster.NewClusterInfo(logger)
	go clusterInfo.Run()

	go collector.StartCollector(t, clusterInfo, logger)

	http.Handle("/metrics", promhttp.Handler())
	logger.Infow("server started",
		"port", 9655,
	)
	err := http.ListenAndServe(":9655", nil)
	check(err)
}
