package main

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/insolar/vanilla/throw"
	"github.com/soverenio/instrumentation/server"

	"github.com/bigbes/go-echo/metrics"
)

func echoHandle(rw http.ResponseWriter, r *http.Request) {
	metrics.CallCount.Inc()

	if ct := r.Header.Get("Content-Type"); ct != "" {
		rw.Header().Set("Content-Type", ct)
	}

	// process body
	if r.ContentLength > 0 {
		_, _ = io.CopyN(rw, r.Body, r.ContentLength)
	} else if r.ContentLength == -1 || r.Body != nil {
		_, _ = io.Copy(rw, r.Body)
	}

	if r.Body != nil {
		_ = r.Body.Close()
	}
}

var (
	configurationPrefix = "GO_ECHO_"
)

func getVariable(suffix string, defaultValue string) string {
	if value, ok := os.LookupEnv(configurationPrefix + suffix); ok {
		return value
	}
	return defaultValue
}

func initializeMetrics(ctx context.Context, port string) server.DefaultServer {
	cfg := server.NewConfiguration()
	cfg.ListenAddress = port

	s := server.NewDefaultServer(cfg)
	s.AddMetrics(ctx, metrics.PrepareRegistry())

	if err := s.Run(ctx); err != nil {
		panic(throw.W(err, "failed to start metrics server"))
	}
	return s
}

func main() {
	ctx := context.Background()

	s := initializeMetrics(ctx, getVariable("INSTRUMENTATION_PORT", ":9090"))
	defer func() { _ = s.Stop(ctx) }()

	http.HandleFunc("/", echoHandle)
	err := http.ListenAndServe(getVariable("ECHO_PORT", ":9000"), nil)
	if err != http.ErrServerClosed {
		panic(throw.W(err, "failed to start server"))
	}
}
