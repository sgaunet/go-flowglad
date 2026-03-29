package internal

import "context"

type noRetryCtxKey struct{}

// WithNoRetry returns a copy of ctx that disables automatic retries for a
// single API call.
func WithNoRetry(ctx context.Context) context.Context {
	return context.WithValue(ctx, noRetryCtxKey{}, true)
}

// NoRetryFromContext reports whether ctx carries the per-call no-retry flag.
func NoRetryFromContext(ctx context.Context) bool {
	v, _ := ctx.Value(noRetryCtxKey{}).(bool)
	return v
}
