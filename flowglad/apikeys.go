package flowglad

import (
	"context"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// APIKeysClient provides methods for the Flowglad API Keys API.
type APIKeysClient struct {
	cfg *internal.HTTPConfig
}

// apiKeyResponse is the JSON envelope for a single API key response.
type apiKeyResponse struct {
	Data struct {
		APIKey APIKey `json:"apiKey"`
	} `json:"data"`
}

// Get retrieves an API key record by ID (metadata only; the secret is not returned).
func (c *APIKeysClient) Get(ctx context.Context, id string) (*APIKey, error) {
	resp, err := internal.Do[apiKeyResponse](ctx, c.cfg, "GET", "/api-keys/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.APIKey, nil
}
