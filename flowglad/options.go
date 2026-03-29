package flowglad

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// RetryPolicy configures automatic retry behaviour for the SDK client.
type RetryPolicy struct {
	// MaxAttempts is the total number of attempts (1 = no retries).
	MaxAttempts int
	// InitialBackoff is the delay before the first retry.
	InitialBackoff time.Duration
	// MaxBackoff caps the delay between retries.
	MaxBackoff time.Duration
	// Multiplier is applied to the backoff duration on each retry.
	Multiplier float64
	// JitterFactor adds ±random jitter to each delay (0.1 = ±10%).
	JitterFactor float64
}

const (
	defaultMaxAttempts  = 3
	defaultMaxBackoff   = 30 * time.Second
	defaultMultiplier   = 2.0
	defaultJitterFactor = 0.1
)

func defaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts:    defaultMaxAttempts,
		InitialBackoff: time.Second,
		MaxBackoff:     defaultMaxBackoff,
		Multiplier:     defaultMultiplier,
		JitterFactor:   defaultJitterFactor,
	}
}

// clientConfig holds the resolved configuration for a Client.
type clientConfig struct {
	baseURL    string
	httpClient *http.Client
	logger     *slog.Logger
	retry      RetryPolicy
	noRetry    bool
}

// Option is a functional option for configuring a Client.
type Option func(*clientConfig)

// WithBaseURL overrides the default Flowglad API base URL.
// Use this to point the SDK at a test server or a different API version.
func WithBaseURL(u string) Option {
	return func(c *clientConfig) { c.baseURL = u }
}

// WithHTTPClient sets a custom HTTP client (e.g. with a custom timeout or
// RoundTripper for tracing/logging middleware).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *clientConfig) { c.httpClient = hc }
}

// WithLogger sets a structured logger for request/response debug tracing.
// Pass nil to fall back to slog.Default().
func WithLogger(l *slog.Logger) Option {
	return func(c *clientConfig) { c.logger = l }
}

// WithRetry overrides the default retry policy.
func WithRetry(p RetryPolicy) Option {
	return func(c *clientConfig) { c.retry = p }
}

// WithNoRetry disables all automatic retries for the client.
func WithNoRetry() Option {
	return func(c *clientConfig) { c.noRetry = true }
}

// NoRetry returns a copy of ctx that disables automatic retries for a single
// API call. Use this to override the client-level retry policy on a per-call
// basis:
//
//	ctx := flowglad.NoRetry(ctx)
//	_, err := client.Subscriptions.Cancel(ctx, id, params)
func NoRetry(ctx context.Context) context.Context {
	return internal.WithNoRetry(ctx)
}
