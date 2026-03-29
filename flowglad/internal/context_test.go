package internal

import (
	"context"
	"testing"
)

func TestWithNoRetry_SetsFlag(t *testing.T) {
	ctx := WithNoRetry(context.Background())
	if !NoRetryFromContext(ctx) {
		t.Error("expected NoRetryFromContext to return true")
	}
}

func TestNoRetryFromContext_DefaultFalse(t *testing.T) {
	if NoRetryFromContext(context.Background()) {
		t.Error("expected NoRetryFromContext to return false for plain context")
	}
}
