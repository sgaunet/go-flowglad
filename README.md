# go-flowglad

[![GitHub release](https://img.shields.io/github/release/sgaunet/go-flowglad.svg)](https://github.com/sgaunet/go-flowglad/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/go-flowglad)](https://goreportcard.com/report/github.com/sgaunet/go-flowglad)
[![Test](https://github.com/sgaunet/go-flowglad/actions/workflows/test.yml/badge.svg)](https://github.com/sgaunet/go-flowglad/actions/workflows/test.yml)
[![GoDoc](https://godoc.org/github.com/sgaunet/go-flowglad?status.svg)](https://godoc.org/github.com/sgaunet/go-flowglad)
[![License](https://img.shields.io/github/license/sgaunet/go-flowglad.svg)](LICENSE)

Unofficial Go SDK for the [Flowglad](https://flowglad.com) billing API.

It provides typed clients for every Flowglad REST endpoint, cursor-based pagination via `iter.Seq2`, configurable retry with exponential backoff, webhook signature verification, and a companion `flowgladtest` package for deterministic unit testing — all with **zero non-stdlib dependencies** in the core package.

## Requirements

- **Go 1.23 or later** (uses `iter.Seq2` range-over-func pagination)

## Installation

```sh
go get github.com/sgaunet/go-flowglad
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/sgaunet/go-flowglad/flowglad"
)

func main() {
    client, err := flowglad.NewClient("sk_live_YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }

    // Create a customer
    cust, err := client.Customers.Create(context.Background(), &flowglad.CreateCustomerParams{
        Name:  "Alice Corp",
        Email: "alice@example.com",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("created:", cust.ID)

    // List customers with range-over-func pagination
    for c, err := range client.Customers.List(context.Background(), nil) {
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(c.ID, c.Email)
    }
}
```

## Supported Resources

| Resource | Methods |
|----------|---------|
| `Customers` | `Create`, `Get`, `Update`, `List`, `Archive`, `GetBillingDetails`, `GetUsageBalances` |
| `Subscriptions` | `Create`, `Get`, `List`, `Adjust`, `PreviewAdjust`, `Cancel`, `Uncancel`, `AddFeature`, `CancelScheduledAdjustment` |
| `Products` | `Create`, `Get`, `Update`, `List` |
| `Prices` | `Create`, `Get`, `Update`, `List` |
| `PricingModels` | `Create`, `Get`, `Update`, `List`, `Clone`, `Export`, `Setup`, `GetDefault` |
| `CheckoutSessions` | `Create`, `Get`, `List` |
| `Invoices` | `Get`, `List` |
| `InvoiceLineItems` | `Get`, `List` |
| `Payments` | `Get`, `List`, `Refund` |
| `PaymentMethods` | `Get`, `List` |
| `UsageEvents` | `Create`, `BulkCreate`, `Get`, `List` |
| `UsageMeters` | `Create`, `Get` |
| `Discounts` | `Create`, `Get`, `Update`, `List` |
| `Features` | `Create`, `Get`, `Update`, `List`, `AddProductFeature`, `ListSubscriptionFeatures` |
| `Resources` | `Create`, `Get`, `Update`, `List`, `Claim`, `Release`, `ListClaims`, `ListUsages` |
| `Webhooks` | `Create`, `Get`, `Update` |
| `APIKeys` | `Get` |

## Pagination

All `List` methods return an `iter.Seq2[*T, error]` iterator that transparently fetches pages. Use the standard Go range loop:

```go
for sub, err := range client.Subscriptions.List(ctx, nil) {
    if err != nil {
        return err
    }
    fmt.Println(sub.ID, sub.Status)
}
```

Pass list params to control page size:

```go
params := &flowglad.ListSubscriptionsParams{
    Limit: flowglad.Ptr(20),
}
for sub, err := range client.Subscriptions.List(ctx, params) {
    // ...
}
```

The `flowglad.Ptr` helper converts any value to a pointer, which is the convention for optional fields:

```go
flowglad.Ptr("a string")   // *string
flowglad.Ptr(42)            // *int
flowglad.Ptr(true)          // *bool
```

## Error Handling

API errors are returned as `*flowglad.Error`, which carries the HTTP status code, a machine-readable code, a human-readable message, and any validation issues.

Use `errors.As` to access structured fields:

```go
_, err := client.Customers.Get(ctx, "nonexistent")
if err != nil {
    var apiErr *flowglad.Error
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d | code: %s | message: %s\n",
            apiErr.StatusCode, apiErr.Code, apiErr.Message)
    }
}
```

## Client Options

`NewClient` accepts functional options to override defaults:

```go
client, err := flowglad.NewClient("sk_live_...",
    // Override the API base URL (useful for self-hosted or tests)
    flowglad.WithBaseURL("https://custom.flowglad.example.com/api/v1"),

    // Provide a custom *http.Client (e.g. with a timeout or tracing middleware)
    flowglad.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),

    // Set a structured slog logger for request/response debug tracing
    flowglad.WithLogger(slog.Default()),

    // Override the default retry policy
    flowglad.WithRetry(flowglad.RetryPolicy{
        MaxAttempts:    5,
        InitialBackoff: 500 * time.Millisecond,
        MaxBackoff:     30 * time.Second,
        Multiplier:     2.0,
        JitterFactor:   0.1,
    }),

    // Disable all automatic retries
    flowglad.WithNoRetry(),
)
```

You can also disable retries for a single call using `flowglad.NoRetry(ctx)`:

```go
ctx := flowglad.NoRetry(ctx)
_, err := client.Subscriptions.Cancel(ctx, id, params)
```

## Retry Policy

The SDK retries automatically on transient failures (network errors, HTTP 5xx, 429). The default policy:

| Setting | Default |
|---------|---------|
| `MaxAttempts` | 3 |
| `InitialBackoff` | 1 s |
| `MaxBackoff` | 30 s |
| `Multiplier` | 2.0 |
| `JitterFactor` | 0.1 (±10%) |

Use `WithNoRetry()` to disable retries entirely, or `WithRetry(p)` to customise.

## Webhook Verification

The `webhook` sub-package verifies Flowglad webhook signatures (HMAC-SHA256) and unmarshals events:

```go
import "github.com/sgaunet/go-flowglad/webhook"

http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    event, err := webhook.Verify(body, r.Header.Get("X-Flowglad-Signature"), webhookSecret)
    if err != nil {
        http.Error(w, "invalid signature", http.StatusUnauthorized)
        return
    }

    switch event.Type {
    case webhook.EventTypeSubscriptionCanceled:
        // handle cancellation
    case webhook.EventTypePaymentSucceeded:
        // handle payment succeeded
    case webhook.EventTypePurchaseCompleted:
        // handle purchase completed
    }
    w.WriteHeader(http.StatusNoContent)
})
```

### Supported Event Types

| Constant | Event |
|----------|-------|
| `EventTypeCustomerCreated` | `customer.created` |
| `EventTypeCustomerUpdated` | `customer.updated` |
| `EventTypePurchaseCompleted` | `purchase.completed` |
| `EventTypePaymentFailed` | `payment.failed` |
| `EventTypePaymentSucceeded` | `payment.succeeded` |
| `EventTypeSubscriptionCreated` | `subscription.created` |
| `EventTypeSubscriptionUpdated` | `subscription.updated` |
| `EventTypeSubscriptionCanceled` | `subscription.canceled` |
| `EventTypeSyncEventsAvailable` | `sync.events_available` |

## Testing with flowgladtest

The `flowgladtest` package provides a fake Flowglad HTTP server for deterministic unit tests — no network required.

```go
import (
    "testing"

    "github.com/sgaunet/go-flowglad/flowglad"
    "github.com/sgaunet/go-flowglad/flowgladtest"
)

func TestMyFeature(t *testing.T) {
    srv := flowgladtest.NewServer(t) // auto-cleaned up via t.Cleanup

    // Register route handlers with fixtures
    srv.On("GET /customers/{id}", flowgladtest.RespondWith(200,
        flowgladtest.CustomerFixture("cus_test_1", "Test Corp", "test@example.com", "ext_1"),
    ))

    client, err := flowglad.NewClient("sk_test_fake",
        flowglad.WithBaseURL(srv.URL()),
        flowglad.WithNoRetry(),
    )
    if err != nil {
        t.Fatal(err)
    }

    cust, err := client.Customers.Get(t.Context(), "cus_test_1")
    if err != nil {
        t.Fatal(err)
    }
    if cust.Name != "Test Corp" {
        t.Errorf("got name %q, want %q", cust.Name, "Test Corp")
    }

    // Inspect recorded requests
    if len(srv.Calls()) != 1 {
        t.Errorf("expected 1 call, got %d", len(srv.Calls()))
    }
}
```

Available helpers:

- `flowgladtest.NewServer(t)` — creates a fake server with automatic cleanup
- `flowgladtest.RespondWith(status, body)` — JSON response handler
- `flowgladtest.RespondWithError(status, code, message)` — error response handler
- `flowgladtest.CustomerFixture(...)` — pre-built customer response
- `flowgladtest.SubscriptionFixture(...)` — pre-built subscription response
- `flowgladtest.CheckoutSessionFixture(...)` — pre-built checkout session response
- `flowgladtest.CustomerListFixture(...)` — pre-built paginated customer list response
- `flowgladtest.ErrorFixture(code, message)` — pre-built error body

## Development

```sh
# Run all tests
go test -race ./...

# Run integration tests (requires FLOWGLAD_API_KEY env var)
go test -race -tags integration ./...

# Lint
golangci-lint run

# Vet
go vet ./...

# Tidy dependencies
go mod tidy
```

## License

[MIT](LICENSE)
