package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/insolar/vanilla/throw"
	"github.com/soverenio/instrumentation/server"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/bigbes/go-echo/metrics"
)

var (
	configurationPrefix = "SVRN_ECHO_"

	configurationAvailableOptions = map[string]struct{}{
		"INSTRUMENTATION_SERVER": {},
		"SERVER":                 {},
		"SERVER_H2C":             {},
	}
)

func main() {
	ctx := context.Background()
	variablesCheck()

	s := initializeMetrics(ctx, variableGet("INSTRUMENTATION_SERVER", ":9090"))
	defer func() { _ = s.Stop(ctx) }()

	errChan := make(chan error)
	httpSrv := serveHTTP(errChan)
	http2Srv := serveH2C(errChan)
	servers := []Shutdowner{httpSrv, http2Srv}

	// SIGINT/SIGTERM handling
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	waitForStopOrError(errChan, osSignals)

	log.Println("Shutting down server...")
	err := shutdown(servers...)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server stopped")
}

func serveHTTP(resChan chan error) Shutdowner {
	addr := variableGet("SERVER", ":9000")
	log.Println("Starting HTTP server on", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(throw.W(err, "failed to listen "+addr))
	}

	srv := &fasthttp.Server{
		Handler: httpHandler,
	}

	go func() {
		resChan <- srv.Serve(lis)
	}()

	return srv
}

func serveH2C(resChan chan error) Shutdowner {
	addr := variableGet("SERVER_H2C", ":9001")
	log.Println("Starting H2C server on", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(throw.W(err, "failed to listen: "+addr))
	}

	http2srv := &http2.Server{}
	srv := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(http.HandlerFunc(h2cHandler), http2srv),
	}

	err = http2.ConfigureServer(srv, http2srv)
	if err != nil {
		panic(throw.W(err, "failed to configure: "+addr))
	}

	go func() {
		resChan <- srv.Serve(lis)
	}()

	return &HttpServerWithShutdown{srv: srv}
}

func httpHandler(ctx *fasthttp.RequestCtx) {
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

func h2cHandler(w http.ResponseWriter, r *http.Request) {
	metrics.CallCount.Inc()

	if len(r.Referer()) > 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	// process body
	switch {
	case r.ContentLength > 0:
		if ct := r.Header.Get("Content-Type"); len(ct) > 0 {
			w.Header().Add("Content-Type", ct)
		}
		w.WriteHeader(http.StatusOK)
		defer r.Body.Close()

		bodyBytes, _ := io.ReadAll(r.Body)
		_, _ = io.WriteString(w, string(bodyBytes))

	case r.ContentLength == -1:
		_, _ = w.Write([]byte("unsupported"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
	}
}

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

func waitForStopOrError(errChan chan error, osSignals chan os.Signal) {
	for {
		select {
		case srvErr := <-errChan:
			if srvErr != nil {
				log.Printf("server error: %s", srvErr)
			}
			return

		case <-osSignals:
			log.Println("")
			log.Println("Shutdown signal received")
			return
		}
	}
}

func shutdown(servers ...Shutdowner) error {
	var errs []error
	for _, srv := range servers {
		err := srv.Shutdown()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return throw.W(errors.Join(errs...), "graceful shutdown error")
}

type Shutdowner interface {
	Shutdown() error
}

type HttpServerWithShutdown struct {
	srv *http.Server
}

// Shutdown gracefully stops the server.
// Implements Shutdowner interface.
func (s *HttpServerWithShutdown) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.srv.Shutdown(ctx)
}
