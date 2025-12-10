package debugmiddleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestDebugMiddleware(t *testing.T) {
	t.Parallel()

	setup := func() (*RequestLogger, *bytes.Buffer) {
		var (
			logBuf     bytes.Buffer
			middleware = NewRequestLogger()
		)
		middleware.logger = log.New(&logBuf, "", 0)
		return middleware, &logBuf
	}

	t.Run("DoesNotRedactMostHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		const stainlessUserAgent = "Stainless"

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set("User-Agent", stainlessUserAgent)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			// The request sent down through middleware shouldn't be mutated.
			if req.Header.Get("User-Agent") != stainlessUserAgent {
				t.Error("expected original request to be unmodified")
			}

			return &http.Response{}, nil
		})

		if !nextMiddlewareRan {
			t.Error("expected next middleware to have been run")
		}

		if !strings.Contains(logBuf.String(), "User-Agent: "+stainlessUserAgent) {
			t.Error("expected logged request headers to include `User-Agent: Stainless`")
		}
	})

	const secretToken = "secret-token"

	t.Run("RedactsAuthorizationHeader", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set("Authorization", secretToken)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			// The request sent down through middleware shouldn't be mutated.
			if req.Header.Get("Authorization") != secretToken {
				t.Error("expected original request to be unmodified")
			}

			return &http.Response{}, nil
		})

		if !nextMiddlewareRan {
			t.Error("expected next middleware to have been run")
		}

		if !strings.Contains(logBuf.String(), "Authorization: "+redactedPlaceholder) {
			t.Error("expected authorization header to be redacted")
		}
	})

	t.Run("RedactsOnlySecretInAuthorizationHeader", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set("Authorization", "Bearer "+secretToken)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			return &http.Response{}, nil
		})

		if !nextMiddlewareRan {
			t.Error("expected next middleware to have been run")
		}

		if !strings.Contains(logBuf.String(), "Authorization: Bearer "+redactedPlaceholder) {
			t.Error("expected authorization header to be redacted")
		}
	})

	t.Run("RedactsMultipleAuthorizationHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Add("Authorization", secretToken+"1")
		req.Header.Add("Authorization", secretToken+"2")

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			// The request sent down through middleware shouldn't be mutated.
			if !reflect.DeepEqual(req.Header.Values("Authorization"), []string{secretToken + "1", secretToken + "2"}) {
				t.Errorf("expected original request to be unmodified")
			}

			return &http.Response{}, nil
		})

		if !nextMiddlewareRan {
			t.Error("expected next middleware to have been run")
		}

		if strings.Count(logBuf.String(), "Authorization: "+redactedPlaceholder) != 2 {
			t.Error("expected exactly two redacted placeholders in authorization headers")
		}
	})

	const customAPIKeyHeader = "X-My-Api-Key"

	t.Run("RedactsSensitiveHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		middleware.sensitiveHeaders = []string{customAPIKeyHeader}

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set(customAPIKeyHeader, secretToken)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			// The request sent down through middleware shouldn't be mutated.
			if req.Header.Get(customAPIKeyHeader) != secretToken {
				t.Error("expected original request to be unmodified")
			}

			return &http.Response{}, nil
		})

		if !nextMiddlewareRan {
			t.Error("expected next middleware to have been run")
		}

		if !strings.Contains(logBuf.String(), customAPIKeyHeader+": "+redactedPlaceholder) {
			t.Errorf("expected %s header to be redacted", customAPIKeyHeader)
		}
	})

	t.Run("RedactsMultipleSensitiveHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		middleware.sensitiveHeaders = []string{customAPIKeyHeader}

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Add(customAPIKeyHeader, secretToken+"1")
		req.Header.Add(customAPIKeyHeader, secretToken+"2")

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			// The request sent down through middleware shouldn't be mutated.
			if !reflect.DeepEqual(req.Header.Values(customAPIKeyHeader), []string{secretToken + "1", secretToken + "2"}) {
				t.Error("expected original request to be unmodified")
			}

			return &http.Response{}, nil
		})

		if !nextMiddlewareRan {
			t.Error("expected next middleware to have been run")
		}

		if strings.Count(logBuf.String(), customAPIKeyHeader+": "+redactedPlaceholder) != 2 {
			t.Errorf("expected %s header to be redacted", customAPIKeyHeader)
		}
	})
}
