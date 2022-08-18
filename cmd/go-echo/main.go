package main

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/insolar/vanilla/throw"
	"github.com/soverenio/instrumentation/server"

	"github.com/bigbes/go-echo/metrics"
	http "github.com/valyala/fasthttp"
)

func echoHandle(ctx *http.RequestCtx) {
	metrics.CallCount.Inc()

	if len(ctx.Request.Header.Referer()) > 0 {
		ctx.SetStatusCode(http.StatusOK)
		return
	}

	// process body
	contentLength := ctx.Request.Header.ContentLength()
	if contentLength > 0 {
		if ct := ctx.Request.Header.ContentType(); len(ct) > 0 {
			ctx.SetContentTypeBytes(ct)
		}

		ctx.SetStatusCode(http.StatusOK)
		_ = ctx.Request.BodyWriteTo(ctx.Response.BodyWriter())
	} else if contentLength == -1 {
		ctx.Error("unsupported", http.StatusUnsupportedMediaType)
		return
	}
}

var (
	configurationPrefix = "SVRN_ECHO_"

	configurationAvailableOptions = map[string]struct{}{
		"INSTRUMENTATION_SERVER": struct{}{},
		"SERVER":                 struct{}{},
	}
)

func variableGet(suffix string, defaultValue string) string {
	if value, ok := os.LookupEnv(configurationPrefix + suffix); ok {
		return value
	}
	return defaultValue
}

func variablesCheck() {
	for _, envVar := range os.Environ() {
		key, _, ok := strings.Cut(envVar, "=")
		switch {
		case !ok:
			continue
		case !strings.HasPrefix(key, configurationPrefix):
			continue
		}

		if _, ok := configurationAvailableOptions[key[len(configurationPrefix):]]; !ok {
			panic("illegal environment variable " + key)
		}
	}
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
	variablesCheck()

	s := initializeMetrics(ctx, variableGet("INSTRUMENTATION_SERVER", ":9090"))
	defer func() { _ = s.Stop(ctx) }()

	err := http.ListenAndServe(variableGet("SERVER", ":9000"), echoHandle)
	if err != io.EOF {
		panic(throw.W(err, "failed to start server"))
	}
}
