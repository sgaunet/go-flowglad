// Package internal provides unexported HTTP transport, retry, and error helpers
// used exclusively by the flowglad package. It MUST NOT be imported by any code
// outside the go-flowglad module.
package internal

import (
	"encoding/json"
	"fmt"
)

// APIError represents a Flowglad API error response. It is aliased as
// flowglad.Error so callers can use errors.As(err, &flowglad.Error{}).
type APIError struct {
	// StatusCode is the HTTP status code returned by the Flowglad API.
	StatusCode int
	// Code is the machine-readable Flowglad error code.
	Code string
	// Message is the human-readable error description.
	Message string
	// Issues contains individual validation issues, if any.
	Issues []string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if len(e.Issues) > 0 {
		return fmt.Sprintf("flowglad: HTTP %d %s: %s (issues: %v)", e.StatusCode, e.Code, e.Message, e.Issues)
	}
	return fmt.Sprintf("flowglad: HTTP %d %s: %s", e.StatusCode, e.Code, e.Message)
}

// apiErrorBody is the internal JSON structure of a Flowglad error response.
type apiErrorBody struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Issues  []any  `json:"issues"`
}

// ParseError decodes a Flowglad error response body into an APIError.
func ParseError(statusCode int, body []byte) *APIError {
	var eb apiErrorBody
	_ = json.Unmarshal(body, &eb)
	issues := make([]string, 0, len(eb.Issues))
	for _, iss := range eb.Issues {
		switch v := iss.(type) {
		case string:
			issues = append(issues, v)
		case map[string]any:
			if msg, ok := v["message"].(string); ok {
				issues = append(issues, msg)
			}
		}
	}
	return &APIError{
		StatusCode: statusCode,
		Code:       eb.Code,
		Message:    eb.Message,
		Issues:     issues,
	}
}
