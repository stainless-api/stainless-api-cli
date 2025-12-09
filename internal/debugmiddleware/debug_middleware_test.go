package debugmiddleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDebugMiddleware(t *testing.T) {
	t.Parallel()

	setup := func() (Middleware, *bytes.Buffer) {
		var logBuf bytes.Buffer
		return DebugMiddleware(log.New(&logBuf, "", 0)), &logBuf
	}

	t.Run("DoesNotRedactMostHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		const stainlessUserAgent = "Stainless"

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set("User-Agent", stainlessUserAgent)

		var nextMiddlewareRan bool
		middleware(req, func(req *http.Request) (*http.Response, error) {
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
		middleware(req, func(req *http.Request) (*http.Response, error) {
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
		middleware(req, func(req *http.Request) (*http.Response, error) {
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
}
