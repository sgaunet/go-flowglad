// Package flowglad is the official Go SDK for the Flowglad billing API.
// It provides typed clients for every Flowglad REST endpoint, cursor-based
// pagination via iter.Seq2, configurable retry with exponential backoff,
// and a companion flowgladtest package for deterministic testing.
package flowglad

import "github.com/sgaunet/go-flowglad/flowglad/internal"

// Error is a typed API error returned by Flowglad. It carries the HTTP status
// code, a machine-readable error code, a human-readable message, and any
// individual validation issues.
//
// Use errors.As to access structured fields:
//
//	var apiErr *flowglad.Error
//	if errors.As(err, &apiErr) {
//	    log.Printf("HTTP %d: %s", apiErr.StatusCode, apiErr.Message)
//	}
type Error = internal.APIError
