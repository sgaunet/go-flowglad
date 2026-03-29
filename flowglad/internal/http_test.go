package internal_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// testConfig returns a minimal HTTPConfig pointing at srv.
func testConfig(srv *httptest.Server, retry *internal.BackoffPolicy) *internal.HTTPConfig {
	return &internal.HTTPConfig{
		BaseURL:    srv.URL,
		APIKey:     "sk_test_fake",
		HTTPClient: srv.Client(),
		Logger:     slog.Default(),
		Retry:      retry,
	}
}

func TestDo_AuthorizationHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"test"}`))
	}))
	defer srv.Close()

	type resp struct{ ID string `json:"id"` }
	_, err := internal.Do[resp](context.Background(), testConfig(srv, nil), "GET", "/test", nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth != "Bearer sk_test_fake" {
		t.Errorf("got Authorization %q, want %q", gotAuth, "Bearer sk_test_fake")
	}
}

func TestDo_4xxReturnsAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"not found","code":"not_found","issues":[]}`))
	}))
	defer srv.Close()

	type resp struct{}
	_, err := internal.Do[resp](context.Background(), testConfig(srv, nil), "GET", "/missing", nil, false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*internal.APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
	if apiErr.Code != "not_found" {
		t.Errorf("expected code 'not_found', got %q", apiErr.Code)
	}
}

func TestDo_5xxRetries(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"server error","code":"internal","issues":[]}`))
	}))
	defer srv.Close()

	bp := &internal.BackoffPolicy{
		MaxAttempts:    3,
		InitialBackoff: time.Millisecond,
		MaxBackoff:     10 * time.Millisecond,
		Multiplier:     2.0,
		JitterFactor:   0.0,
	}

	type resp struct{}
	_, err := internal.Do[resp](context.Background(), testConfig(srv, bp), "GET", "/fail", nil, true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestDo_429RespectsRetryAfter(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.Header().Set("Retry-After", "0")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"message":"rate limited","code":"rate_limit","issues":[]}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"id": "ok"})
	}))
	defer srv.Close()

	bp := &internal.BackoffPolicy{
		MaxAttempts:    3,
		InitialBackoff: time.Millisecond,
		MaxBackoff:     10 * time.Millisecond,
		Multiplier:     2.0,
		JitterFactor:   0.0,
	}

	type resp struct{ ID string `json:"id"` }
	result, err := internal.Do[resp](context.Background(), testConfig(srv, bp), "GET", "/rate", nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "ok" {
		t.Errorf("expected id 'ok', got %q", result.ID)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestDo_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	type resp struct{}
	_, err := internal.Do[resp](ctx, testConfig(srv, nil), "GET", "/slow", nil, false)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestDo_NoRetryOn4xx(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"bad","code":"bad_request","issues":[]}`))
	}))
	defer srv.Close()

	bp := &internal.BackoffPolicy{
		MaxAttempts:    3,
		InitialBackoff: time.Millisecond,
		MaxBackoff:     10 * time.Millisecond,
		Multiplier:     2.0,
		JitterFactor:   0.0,
	}

	type resp struct{}
	_, err := internal.Do[resp](context.Background(), testConfig(srv, bp), "GET", "/bad", nil, true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt (no retry on 4xx), got %d", attempts)
	}
}

func TestBuildQueryString(t *testing.T) {
	tests := []struct {
		cursor string
		limit  *int
		want   string
	}{
		{"", nil, ""},
		{"cur123", nil, "?cursor=cur123"},
		{"", intPtr(10), "?limit=10"},
		{"cur123", intPtr(20), "?cursor=cur123&limit=20"},
	}
	for _, tc := range tests {
		got := internal.BuildQueryString(tc.cursor, tc.limit)
		if got != tc.want {
			t.Errorf("BuildQueryString(%q, %v) = %q, want %q", tc.cursor, tc.limit, got, tc.want)
		}
	}
}

func intPtr(v int) *int { return &v }
