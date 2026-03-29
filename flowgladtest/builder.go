package flowgladtest

import (
	"encoding/json"
	"net/http"
)

// RespondWith returns an http.HandlerFunc that writes a JSON response with
// the given status code and body. body will be JSON-marshalled before writing.
func RespondWith(statusCode int, body any) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}
}

// RespondWithError returns an http.HandlerFunc that writes a Flowglad-shaped
// error response body with the given HTTP status code, error code, and message.
func RespondWithError(statusCode int, code, message string) http.HandlerFunc {
	return RespondWith(statusCode, map[string]any{
		"code":    code,
		"message": message,
		"issues":  []any{},
	})
}
