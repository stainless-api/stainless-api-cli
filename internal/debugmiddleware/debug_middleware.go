package debugmiddleware

import (
	"net/http"
	"net/http/httputil"
	"strings"
)

// For the time being these type definitions are duplicated here so that we can
// test this file in a non-generated context.
type (
	Middleware     = func(*http.Request, MiddlewareNext) (*http.Response, error)
	MiddlewareNext = func(*http.Request) (*http.Response, error)
)

const redactedPlaceholder = "<REDACTED>"

// DebugMiddleware returns a middleware that logs HTTP requests and responses.
//
// logWriter is log.Default() under most circumstances, but made low level so we
// can more easily inject a buffer to check in tests.
func DebugMiddleware(logger interface{ Printf(string, ...any) }) Middleware {
	return func(req *http.Request, mn MiddlewareNext) (*http.Response, error) {
		if reqBytes, err := httputil.DumpRequest(redactRequest(req), true); err == nil {
			logger.Printf("Request Content:\n%s\n", reqBytes)
		}

		resp, err := mn(req)
		if err != nil {
			return resp, err
		}

		if respBytes, err := httputil.DumpResponse(resp, true); err == nil {
			logger.Printf("Response Content:\n%s\n", respBytes)
		}

		return resp, err
	}
}

// redactRequest redacts sensitive information from the request for logging
// purposes. If redaction is necessary, the request is cloned before mutating
// the original and that clone is returned. As a small optimization, the
// original is request is returned unchanged if no redaction is necessary.
func redactRequest(req *http.Request) *http.Request {
	if auth := req.Header.Get("Authorization"); auth != "" {
		req = req.Clone(req.Context())

		// In case we're using something like a bearer token (e.g. `Bearer
		// <my_token>`), keep the `Bearer` part for more debugging
		// information.
		if authKind, _, ok := strings.Cut(auth, " "); ok {
			req.Header.Set("Authorization", authKind+" "+redactedPlaceholder)
		} else {
			req.Header.Set("Authorization", redactedPlaceholder)
		}
	}

	return req
}
