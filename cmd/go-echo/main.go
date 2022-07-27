package main

import (
	"context"
	"io"
	"os"

	"github.com/insolar/vanilla/throw"
	"github.com/soverenio/instrumentation/server"

	"github.com/bigbes/go-echo/metrics"
	http "github.com/valyala/fasthttp"
)

func echoHandle(ctx *http.RequestCtx) {
	metrics.CallCount.Inc()

	// process body
	contentLength := ctx.Request.Header.ContentLength()
	if contentLength > 0 {
		if ct := ctx.Request.Header.ContentType(); len(ct) > 0 {
			ctx.SetContentTypeBytes(ct)
		}

		_ = ctx.Request.BodyWriteTo(ctx.Response.BodyWriter())
	} else if contentLength == -1 {
		ctx.Error("unsupported", http.StatusUnsupportedMediaType)
		return
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

	err := http.ListenAndServe(getVariable("ECHO_PORT", ":9000"), echoHandle)
	if err != io.EOF {
		panic(throw.W(err, "failed to start server"))
	}
}
