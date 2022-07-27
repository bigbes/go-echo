package trace

import (
	"net/http"
)

const (
	XTraceID = "X-Trace-Id"
	XSpanID  = "X-Span-Id"
)

// Client is wrapper around http.Client that injects traceID to outgoing http requests
type Client struct {
	*http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{
		Client: client,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	trace := ID(req.Context())
	if trace != "" {
		req.Header.Set(XTraceID, trace)
	}
	return c.Client.Do(req)
}

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type WrappedHttpRequestDoer struct {
	client HttpRequestDoer
}

func NewHttpRequestDoer(client HttpRequestDoer) WrappedHttpRequestDoer {
	return WrappedHttpRequestDoer{
		client: client,
	}
}

func (c WrappedHttpRequestDoer) Do(req *http.Request) (*http.Response, error) {
	trace := ID(req.Context())
	if trace != "" {
		req.Header.Set(XTraceID, trace)
	}
	return c.client.Do(req)
}
