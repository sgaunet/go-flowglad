package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// HTTPConfig holds configuration for the HTTP layer.
type HTTPConfig struct {
	// BaseURL is the API base URL without a trailing slash.
	BaseURL string
	// APIKey is the Bearer token for authentication. Never logged.
	APIKey string //nolint:gosec // field name matches pattern but is not a raw secret
	// HTTPClient is the underlying transport.
	HTTPClient *http.Client
	// Logger receives structured request/response logs.
	Logger *slog.Logger
	// Retry configures automatic retry behaviour. May be nil (no retry).
	Retry *BackoffPolicy
	// NoRetry disables all automatic retries when true.
	NoRetry bool
}

// Do executes an authenticated HTTP request against the Flowglad API and
// JSON-decodes the response into a value of type T.
//
// method is the HTTP verb ("GET", "POST", "PUT", "PATCH", "DELETE").
// path is the URL path including any query string (e.g. "/customers?limit=20").
// body is the request payload marshalled to JSON, or nil.
// shouldRetry indicates whether the caller allows this particular request to be
// retried on transient failures (5xx, 429, network errors).
func Do[T any]( //nolint:ireturn // generic HTTP helper must return T
	ctx context.Context, cfg *HTTPConfig, method, path string, body any, shouldRetry bool,
) (T, error) {
	bodyBytes, err := marshalBody(body)
	if err != nil {
		var zero T
		return zero, err
	}

	var result T
	attempt := func(n int) (bool, error) {
		return executeOnce(ctx, cfg, method, path, bodyBytes, shouldRetry, n, &result)
	}

	if cfg.NoRetry || cfg.Retry == nil || NoRetryFromContext(ctx) {
		_, err = attempt(1)
		return result, err
	}
	err = cfg.Retry.DoWithContext(ctx, attempt)
	return result, err
}

// marshalBody JSON-encodes body, returning nil bytes when body is nil.
func marshalBody(body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("flowglad: marshal request: %w", err)
	}
	return b, nil
}

// executeOnce performs a single HTTP round-trip and decodes the result into *out.
func executeOnce[T any](
	ctx context.Context,
	cfg *HTTPConfig,
	method, path string,
	bodyBytes []byte,
	shouldRetry bool,
	attempt int,
	out *T,
) (bool, error) {
	req, err := buildRequest(ctx, cfg, method, path, bodyBytes)
	if err != nil {
		return false, err
	}

	start := time.Now()
	resp, err := cfg.HTTPClient.Do(req) //nolint:gosec // URL comes from trusted HTTPConfig.BaseURL
	if err != nil {
		return shouldRetry, fmt.Errorf("flowglad: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	cfg.Logger.DebugContext(ctx, "flowglad request",
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status", resp.StatusCode),
		slog.Duration("latency", time.Since(start)),
		slog.Int("attempt", attempt),
	)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return shouldRetry, fmt.Errorf("flowglad: read response: %w", err)
	}

	retry, apiErr := checkStatus(ctx, resp, respBody, shouldRetry)
	if apiErr != nil {
		return retry, apiErr
	}

	if err := json.Unmarshal(respBody, out); err != nil {
		return false, fmt.Errorf("flowglad: decode response: %w", err)
	}
	return false, nil
}

// buildRequest creates an authenticated HTTP request.
func buildRequest(ctx context.Context, cfg *HTTPConfig, method, path string, bodyBytes []byte) (*http.Request, error) {
	var reqBody io.Reader
	if bodyBytes != nil {
		reqBody = bytes.NewReader(bodyBytes)
	}
	req, err := http.NewRequestWithContext(ctx, method, cfg.BaseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("flowglad: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// checkStatus inspects the HTTP response status and handles rate-limit and error responses.
// It returns (shouldRetry, error). A nil error means the response is success (2xx).
func checkStatus(ctx context.Context, resp *http.Response, body []byte, shouldRetry bool) (bool, error) {
	switch {
	case resp.StatusCode == http.StatusTooManyRequests:
		if err := sleepRetryAfter(ctx, resp.Header.Get("Retry-After")); err != nil {
			return false, err
		}
		return shouldRetry, ParseError(resp.StatusCode, body)
	case resp.StatusCode >= http.StatusInternalServerError:
		return shouldRetry, ParseError(resp.StatusCode, body)
	case resp.StatusCode >= http.StatusBadRequest:
		return false, ParseError(resp.StatusCode, body)
	}
	return false, nil
}

// sleepRetryAfter honours the Retry-After header before the next retry attempt.
func sleepRetryAfter(ctx context.Context, header string) error {
	if header == "" {
		return nil
	}
	secs, err := strconv.Atoi(header)
	if err != nil {
		return nil //nolint:nilerr // invalid Retry-After values are intentionally ignored
	}
	select {
	case <-ctx.Done():
		return fmt.Errorf("flowglad: %w", ctx.Err())
	case <-time.After(time.Duration(secs) * time.Second):
		return nil
	}
}

// BuildQueryString builds a URL query string from cursor/limit pagination params.
func BuildQueryString(cursor string, limit *int) string {
	params := url.Values{}
	if cursor != "" {
		params.Set("cursor", cursor)
	}
	if limit != nil {
		params.Set("limit", strconv.Itoa(*limit))
	}
	if len(params) == 0 {
		return ""
	}
	return "?" + params.Encode()
}
