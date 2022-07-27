package logger

import (
	"net/http"
	"strings"
	"time"

	"github.com/insolar/vanilla/throw"
	"github.com/labstack/echo/v4"
	"github.com/soverenio/log"

	"github.com/soverenio/instrumentation/trace"
)

type (
	// RequestIDConfig defines the config for RequestID middleware.
	RequestIDConfig struct {
		// Generator defines a function to generate an ID.
		// Optional. Default value random.String(32).
		Generator func() string
	}

	requestLog struct {
		*log.Msg `txt:"Incoming request"`

		Host          string
		Path          string
		RawPath       string
		Method        string
		RemoteIP      string
		ContentLength string

		Fields map[string]string
	}
	responseLog struct {
		*log.Msg `txt:"Request result"`

		Latency       string
		Status        int
		ContentLength int64
	}
)

var (
	// DefaultRequestIDConfig is the default RequestID middleware config.
	DefaultRequestIDConfig = RequestIDConfig{
		Generator: generator,
	}
)

// TraceID returns a tracing logger middleware.
func TraceID() echo.MiddlewareFunc {
	return TraceIDWithConfig(DefaultRequestIDConfig)
}

// TraceIDWithConfig returns a tracing logger middleware with config.
func TraceIDWithConfig(config RequestIDConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Generator == nil {
		config.Generator = generator
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			log := FromContext(c.Request().Context())

			req := c.Request()
			res := c.Response()

			traceID := req.Header.Get(trace.XTraceID)
			if traceID == "" {
				traceID = config.Generator()
			}
			log = log.WithField("trace_id", traceID)

			res.Header().Set(trace.XTraceID, traceID)
			ctx, err := trace.SetID(req.Context(), traceID)
			if err != nil {
				log.Error(throw.W(err, "failed to set trace_id"))
			}
			ctx = SetLogger(ctx, log)

			c.SetRequest(req.WithContext(ctx))

			requestLogger := log.WithFields(map[string]interface{}{
				"host":           req.Host,
				"path":           c.Path(),
				"raw_path":       req.RequestURI,
				"method":         req.Method,
				"remote_ip":      c.RealIP(),
				"content_length": getContentLength(req),
				"fields":         getParams(c),
			})
			requestLogger.Info("Incoming request")

			defer func() {
				responseLogger := log.WithFields(map[string]interface{}{
					"latency":        time.Since(start).String(),
					"status":         res.Status,
					"content_length": res.Size,
				})
				responseLogger.Info("Request result")
			}()

			return next(c)
		}
	}
}

func getContentLength(req *http.Request) string {
	cl := req.Header.Get(echo.HeaderContentLength)
	if cl == "" {
		cl = "0"
	}
	return cl
}

func getParams(ctx echo.Context) map[string]string {
	requestParams := make(map[string]string)
	values := ctx.ParamValues()
	for i, name := range ctx.ParamNames() {
		if i < len(values) {
			requestParams[getOneStyleString(name)] = values[i]
		}
	}
	return requestParams
}

func getOneStyleString(s string) string {
	// KeyName, keyName, key-name => key_name
	result := make([]rune, 0, len(s))
	for _, v := range s {
		// upper value
		if v >= 'A' && v <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, v)
	}
	resultSt := string(result)
	resultSt = strings.ToLower(resultSt)
	resultSt = strings.ReplaceAll(resultSt, "-", "_")
	resultSt = strings.Trim(resultSt, "_")
	return resultSt
}

func generator() string {
	return trace.RandID()
}
