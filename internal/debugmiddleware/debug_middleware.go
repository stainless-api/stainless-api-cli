package debugmiddleware

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

// For the time being these type definitions are duplicated here so that we can
// test this file in a non-generated context.
type (
	Middleware     = func(*http.Request, MiddlewareNext) (*http.Response, error)
	MiddlewareNext = func(*http.Request) (*http.Response, error)
)

const redactedPlaceholder = "<REDACTED>"

// Headers known to contain sensitive information like an API key.
var sensitiveHeaders = []string{}

// RequestLogger is a middleware that logs HTTP requests and responses.
type RequestLogger struct {
	logger           interface{ Printf(string, ...any) } // field for testability; usually log.Default()
	sensitiveHeaders []string                            // field for testability; usually sensitiveHeaders
}

// NewRequestLogger returns a new RequestLogger instance with default options.
func NewRequestLogger() *RequestLogger {
	return &RequestLogger{
		logger:           log.Default(),
		sensitiveHeaders: sensitiveHeaders,
	}
}

func (m *RequestLogger) Middleware() Middleware {
	return func(req *http.Request, mn MiddlewareNext) (*http.Response, error) {
		if reqBytes, err := httputil.DumpRequest(m.redactRequest(req), true); err == nil {
			m.logger.Printf("Request Content:\n%s\n", reqBytes)
		}

		resp, err := mn(req)
		if err != nil {
			return resp, err
		}

		if respBytes, err := httputil.DumpResponse(resp, true); err == nil {
			m.logger.Printf("Response Content:\n%s\n", respBytes)
		}

		return resp, err
	}
}

// redactRequest redacts sensitive information from the request for logging
// purposes. If redaction is necessary, the request is cloned before mutating
// the original and that clone is returned. As a small optimization, the
// original is request is returned unchanged if no redaction is necessary.
func (m *RequestLogger) redactRequest(req *http.Request) *http.Request {
	cloneReq := sync.OnceFunc(func() {
		req = req.Clone(req.Context())
	})

	// Notably, the clauses below are written so they can redact multiple
	// headers of the same name if necessary.
	if values := req.Header.Values("Authorization"); len(values) > 0 {
		cloneReq()
		req.Header.Del("Authorization")

		for _, value := range values {
			// In case we're using something like a bearer token (e.g. `Bearer
			// <my_token>`), keep the `Bearer` part for more debugging
			// information.
			if authKind, _, ok := strings.Cut(value, " "); ok {
				req.Header.Add("Authorization", authKind+" "+redactedPlaceholder)
			} else {
				req.Header.Add("Authorization", redactedPlaceholder)
			}
		}
	}

	for _, header := range m.sensitiveHeaders {
		values := req.Header.Values(header)
		if len(values) == 0 {
			continue
		}

		cloneReq()
		req.Header.Del(header)

		for range values {
			req.Header.Add(header, redactedPlaceholder)
		}
	}

	return req
}
