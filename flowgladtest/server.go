// Package flowgladtest provides a configurable in-process HTTP server for
// testing code that uses the Flowglad Go SDK without making real API calls.
//
// Usage:
//
//	func TestMyCode(t *testing.T) {
//	    srv := flowgladtest.NewServer(t)
//	    srv.On("POST /customers", flowgladtest.RespondWith(200, customerFixture))
//
//	    client, _ := flowglad.NewClient("sk_test_fake", flowglad.WithBaseURL(srv.URL()))
//	    cust, err := client.Customers.Create(ctx, params)
//	    // assert cust fields ...
//	    // assert srv.Calls() ...
//	}
package flowgladtest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// methodPathParts is the number of parts in a "METHOD /path" pattern string.
const methodPathParts = 2

// Server is a configurable fake Flowglad HTTP server for use in tests.
// It is automatically shut down when the test ends via t.Cleanup.
type Server struct {
	srv      *httptest.Server
	mu       sync.Mutex
	routes   map[string]http.HandlerFunc // "METHOD /path" → handler
	fallback http.HandlerFunc
	calls    []*http.Request
}

// NewServer creates a new fake Flowglad server and registers cleanup with tb.
// Call srv.On() to register route handlers before making SDK calls.
func NewServer(tb testing.TB) *Server {
	tb.Helper()
	s := &Server{
		routes: make(map[string]http.HandlerFunc),
	}
	s.srv = httptest.NewServer(http.HandlerFunc(s.dispatch))
	tb.Cleanup(func() {
		s.srv.Close()
	})
	return s
}

// On registers an http.HandlerFunc for the given method+path pattern.
// Pattern may use Go 1.22 style path parameters: "GET /customers/{id}".
//
// The last registered handler for a pattern wins.
func (s *Server) On(pattern string, h http.HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes[pattern] = h
}

// OnDefault registers a fallback handler called when no specific pattern matches.
func (s *Server) OnDefault(h http.HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fallback = h
}

// URL returns the base URL to pass to flowglad.WithBaseURL().
func (s *Server) URL() string {
	return s.srv.URL
}

// Calls returns all recorded HTTP requests in the order they were received.
func (s *Server) Calls() []*http.Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]*http.Request, len(s.calls))
	copy(result, s.calls)
	return result
}

// dispatch is the underlying http.Handler that routes requests.
func (s *Server) dispatch(w http.ResponseWriter, r *http.Request) {
	// Record the call.
	s.mu.Lock()
	s.calls = append(s.calls, r)
	s.mu.Unlock()

	// Try exact match first, then prefix match for path parameters.
	key := r.Method + " " + r.URL.Path
	if h := s.routeHandler(key, r); h != nil {
		h(w, r)
		return
	}

	if s.fallback != nil {
		s.fallback(w, r)
		return
	}

	http.Error(w, "flowgladtest: no handler registered for "+key, http.StatusNotImplemented)
}

// routeHandler finds the best handler for the given key.
// It supports patterns like "GET /customers/{id}" using simple prefix matching.
func (s *Server) routeHandler(key string, r *http.Request) http.HandlerFunc {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Exact match.
	if h, ok := s.routes[key]; ok {
		return h
	}

	// Pattern match: try all registered patterns and pick the best.
	method := r.Method
	path := r.URL.Path
	for pattern, h := range s.routes {
		parts := strings.SplitN(pattern, " ", methodPathParts)
		if len(parts) != methodPathParts {
			continue
		}
		if parts[0] != method {
			continue
		}
		if matchPath(parts[1], path) {
			return h
		}
	}
	return nil
}

// matchPath checks whether urlPath matches a pattern that may contain {param} segments.
func matchPath(pattern, path string) bool {
	pParts := strings.Split(strings.Trim(pattern, "/"), "/")
	uParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pParts) != len(uParts) {
		return false
	}
	for i, pp := range pParts {
		if strings.HasPrefix(pp, "{") && strings.HasSuffix(pp, "}") {
			continue // wildcard segment
		}
		if pp != uParts[i] {
			return false
		}
	}
	return true
}
