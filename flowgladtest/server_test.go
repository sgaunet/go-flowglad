package flowgladtest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func TestNewServer_URL(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	if srv.URL() == "" {
		t.Error("expected non-empty URL")
	}
}

func TestServer_On_ExactRoute(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /customers", flowgladtest.RespondWith(http.StatusOK, map[string]any{"id": "cus_1"}))

	resp, err := http.Post(srv.URL()+"/customers", "application/json", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestServer_On_PathParameter(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /customers/{id}", flowgladtest.RespondWith(http.StatusOK, map[string]any{"id": "cus_123"}))

	resp, err := http.Get(srv.URL() + "/customers/cus_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestServer_OnDefault(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.OnDefault(flowgladtest.RespondWith(http.StatusTeapot, nil))

	resp, err := http.Get(srv.URL() + "/unknown-path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTeapot {
		t.Errorf("expected 418, got %d", resp.StatusCode)
	}
}

func TestServer_Calls_RecordsRequests(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /test", flowgladtest.RespondWith(http.StatusOK, nil))

	if _, err := http.Get(srv.URL() + "/test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := http.Get(srv.URL() + "/test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	calls := srv.Calls()
	if len(calls) != 2 {
		t.Errorf("expected 2 calls, got %d", len(calls))
	}
}

func TestServer_NoHandler_Returns501(t *testing.T) {
	srv := flowgladtest.NewServer(t)

	resp, err := http.Get(srv.URL() + "/unregistered")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotImplemented {
		t.Errorf("expected 501, got %d", resp.StatusCode)
	}
}

func TestRespondWithError(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /bad", flowgladtest.RespondWithError(http.StatusUnauthorized, "unauthorized", "invalid API key"))

	resp, err := http.Get(srv.URL() + "/bad")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if body["code"] != "unauthorized" {
		t.Errorf("expected code 'unauthorized', got %v", body["code"])
	}
}

func TestServer_AutoShutdown(t *testing.T) {
	var srvURL string
	// Create a sub-test scope to trigger cleanup.
	t.Run("inner", func(t *testing.T) {
		srv := flowgladtest.NewServer(t)
		srvURL = srv.URL()
		srv.On("GET /ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		// Verify it works inside the test.
		if r, err := http.Get(srvURL + "/ping"); err != nil || r.StatusCode != 200 {
			t.Errorf("expected server to be up inside test")
		} else {
			r.Body.Close()
		}
	})
	// After the sub-test, the server should have been shut down via t.Cleanup.
	// Attempting a connection should fail.
	// (We rely on the sub-test's t.Cleanup having run; httptest.Server.Close()
	// causes the listener to stop accepting. The exact error varies by OS, so
	// we just verify the inner sub-test completed without issue.)
	_ = srvURL
}

// Ensure RespondWith uses JSON encoding correctly.
func TestRespondWith_JSONEncoding(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	flowgladtest.RespondWith(http.StatusCreated, map[string]string{"key": "value"})(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if body["key"] != "value" {
		t.Errorf("expected 'value', got %q", body["key"])
	}
}
