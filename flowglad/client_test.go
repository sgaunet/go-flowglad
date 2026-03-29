package flowglad_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func TestNewClient_EmptyAPIKey(t *testing.T) {
	_, err := flowglad.NewClient("")
	if err == nil {
		t.Fatal("expected error for empty API key, got nil")
	}
}

func TestNewClient_ValidKey_AllSubClientsNonNil(t *testing.T) {
	c, err := flowglad.NewClient("sk_test_fake")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Customers == nil {
		t.Error("Customers is nil")
	}
	if c.Subscriptions == nil {
		t.Error("Subscriptions is nil")
	}
	if c.Products == nil {
		t.Error("Products is nil")
	}
	if c.Prices == nil {
		t.Error("Prices is nil")
	}
	if c.CheckoutSessions == nil {
		t.Error("CheckoutSessions is nil")
	}
	if c.Invoices == nil {
		t.Error("Invoices is nil")
	}
	if c.InvoiceLineItems == nil {
		t.Error("InvoiceLineItems is nil")
	}
	if c.Payments == nil {
		t.Error("Payments is nil")
	}
	if c.PaymentMethods == nil {
		t.Error("PaymentMethods is nil")
	}
	if c.UsageEvents == nil {
		t.Error("UsageEvents is nil")
	}
	if c.UsageMeters == nil {
		t.Error("UsageMeters is nil")
	}
	if c.Discounts == nil {
		t.Error("Discounts is nil")
	}
	if c.Webhooks == nil {
		t.Error("Webhooks is nil")
	}
	if c.PricingModels == nil {
		t.Error("PricingModels is nil")
	}
	if c.Features == nil {
		t.Error("Features is nil")
	}
	if c.Resources == nil {
		t.Error("Resources is nil")
	}
	if c.APIKeys == nil {
		t.Error("APIKeys is nil")
	}
}

func TestNewClient_WithBaseURL(t *testing.T) {
	// Should not error; we just verify construction succeeds.
	_, err := flowglad.NewClient("sk_test_fake", flowglad.WithBaseURL("http://localhost:9999"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewClient_WithNilLogger_FallsBackToDefault(t *testing.T) {
	// Passing nil logger should not panic.
	c, err := flowglad.NewClient("sk_test_fake", flowglad.WithLogger(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_WithHTTPClient(t *testing.T) {
	hc := &http.Client{Timeout: 5 * time.Second}
	c, err := flowglad.NewClient("sk_test_fake", flowglad.WithHTTPClient(hc))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_WithRetry(t *testing.T) {
	policy := flowglad.RetryPolicy{
		MaxAttempts:    5,
		InitialBackoff: 2 * time.Second,
		MaxBackoff:     60 * time.Second,
		Multiplier:     3.0,
		JitterFactor:   0.2,
	}
	c, err := flowglad.NewClient("sk_test_fake", flowglad.WithRetry(policy))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_WithNoRetry(t *testing.T) {
	c, err := flowglad.NewClient("sk_test_fake", flowglad.WithNoRetry())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNoRetry_PerCallDisablesRetry(t *testing.T) {
	var calls atomic.Int32
	srv := flowgladtest.NewServer(t)
	srv.OnDefault(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":"internal","message":"boom"}`))
	})

	// Client has retries enabled (default 3 attempts).
	client, err := flowglad.NewClient("sk_test_fake", flowglad.WithBaseURL(srv.URL()))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	// Per-call NoRetry should prevent any retry attempts.
	ctx := flowglad.NoRetry(context.Background())
	_, err = client.Customers.Get(ctx, "cus_1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := calls.Load(); got != 1 {
		t.Errorf("expected exactly 1 call (no retries), got %d", got)
	}
}
