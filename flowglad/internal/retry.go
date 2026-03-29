package internal

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

// BackoffPolicy defines exponential backoff retry behaviour.
type BackoffPolicy struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// InitialBackoff is the delay before the second attempt.
	InitialBackoff time.Duration
	// MaxBackoff caps the delay at this value.
	MaxBackoff time.Duration
	// Multiplier is applied to the previous delay on each retry.
	Multiplier float64
	// JitterFactor controls ±random jitter: delay ± delay*JitterFactor.
	JitterFactor float64
}

// NextDelay returns the delay before the given 1-indexed attempt number.
// Attempt 1 is the first retry (after the initial attempt).
func (b *BackoffPolicy) NextDelay(attempt int) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}
	base := float64(b.InitialBackoff) * math.Pow(b.Multiplier, float64(attempt-1))
	if base > float64(b.MaxBackoff) {
		base = float64(b.MaxBackoff)
	}
	// Apply ±JitterFactor random jitter. math/rand is intentional; jitter does
	// not require cryptographic randomness.
	jitter := base * b.JitterFactor * (rand.Float64()*2 - 1) //nolint:gosec
	d := max(time.Duration(base+jitter), 0)
	return min(d, b.MaxBackoff)
}

// DoWithContext executes fn repeatedly according to the policy.
//
// fn receives the current attempt number (1-indexed). It must return
// (shouldRetry bool, err error). When shouldRetry is false or the maximum
// number of attempts is reached, DoWithContext stops and returns the last error.
// A nil error from fn means success and stops the loop immediately.
func (b *BackoffPolicy) DoWithContext(ctx context.Context, fn func(attempt int) (bool, error)) error {
	var lastErr error
	for attempt := 1; attempt <= b.MaxAttempts; attempt++ {
		retry, err := fn(attempt)
		if err == nil {
			return nil
		}
		lastErr = err
		if !retry || attempt == b.MaxAttempts {
			break
		}
		delay := b.NextDelay(attempt)
		select {
		case <-ctx.Done():
			return fmt.Errorf("flowglad/internal: %w", ctx.Err())
		case <-time.After(delay):
		}
	}
	return lastErr
}
