package internal_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

func TestBackoffPolicy_NextDelay(t *testing.T) {
	bp := &internal.BackoffPolicy{
		MaxAttempts:    3,
		InitialBackoff: time.Second,
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
		JitterFactor:   0.1,
	}

	for attempt := 1; attempt <= 3; attempt++ {
		d := bp.NextDelay(attempt)
		base := min(
			time.Duration(float64(time.Second)*math.Pow(2.0, float64(attempt-1))),
			30*time.Second,
		)
		lower := time.Duration(float64(base) * 0.9)
		upper := time.Duration(float64(base) * 1.1)
		if d < lower || d > upper {
			t.Errorf("attempt %d: delay %v not in [%v, %v]", attempt, d, lower, upper)
		}
	}
}

func TestBackoffPolicy_CapAt30s(t *testing.T) {
	bp := &internal.BackoffPolicy{
		MaxAttempts:    10,
		InitialBackoff: time.Second,
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
		JitterFactor:   0.0,
	}
	d := bp.NextDelay(10) // would be 512s uncapped
	if d > 30*time.Second {
		t.Errorf("expected delay capped at 30s, got %v", d)
	}
}

func TestBackoffPolicy_ContextCancellation(t *testing.T) {
	bp := &internal.BackoffPolicy{
		MaxAttempts:    5,
		InitialBackoff: 10 * time.Second, // long enough to be cancelled
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
		JitterFactor:   0.0,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	attempts := 0
	err := bp.DoWithContext(ctx, func(attempt int) (bool, error) {
		attempts++
		return true, fmt.Errorf("transient error")
	})
	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Errorf("expected context error, got %v", err)
	}
	if attempts > 2 {
		t.Errorf("expected at most 2 attempts before cancel, got %d", attempts)
	}
}

func TestBackoffPolicy_NoRetryOnFalse(t *testing.T) {
	bp := &internal.BackoffPolicy{
		MaxAttempts:    3,
		InitialBackoff: time.Millisecond,
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
		JitterFactor:   0.0,
	}

	calls := 0
	err := bp.DoWithContext(context.Background(), func(attempt int) (bool, error) {
		calls++
		return false, fmt.Errorf("non-retryable error")
	})
	if err == nil {
		t.Error("expected error, got nil")
	}
	if calls != 1 {
		t.Errorf("expected 1 call (no retry), got %d", calls)
	}
}

func TestBackoffPolicy_SuccessOnSecondAttempt(t *testing.T) {
	bp := &internal.BackoffPolicy{
		MaxAttempts:    3,
		InitialBackoff: time.Millisecond,
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
		JitterFactor:   0.0,
	}

	calls := 0
	err := bp.DoWithContext(context.Background(), func(attempt int) (bool, error) {
		calls++
		if attempt < 2 {
			return true, fmt.Errorf("transient")
		}
		return false, nil
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestBackoffPolicy_MaxAttemptsExhausted(t *testing.T) {
	bp := &internal.BackoffPolicy{
		MaxAttempts:    3,
		InitialBackoff: time.Millisecond,
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
		JitterFactor:   0.0,
	}

	calls := 0
	err := bp.DoWithContext(context.Background(), func(attempt int) (bool, error) {
		calls++
		return true, fmt.Errorf("always fails")
	})
	if err == nil {
		t.Error("expected error, got nil")
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestBackoffPolicy_ZeroDelay(t *testing.T) {
	bp := &internal.BackoffPolicy{
		MaxAttempts:    3,
		InitialBackoff: time.Millisecond,
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
		JitterFactor:   0.0,
	}
	d := bp.NextDelay(0) // should not panic
	if d < 0 {
		t.Errorf("expected non-negative delay, got %v", d)
	}
}
