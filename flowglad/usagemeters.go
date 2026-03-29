package flowglad

import (
	"context"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// UsageMetersClient provides methods for the Flowglad Usage Meters API.
type UsageMetersClient struct {
	cfg *internal.HTTPConfig
}

// usageMeterResponse is the JSON envelope for a single usage meter response.
type usageMeterResponse struct {
	Data struct {
		UsageMeter UsageMeter `json:"usageMeter"`
	} `json:"data"`
}

// Create creates a new usage meter.
func (c *UsageMetersClient) Create(ctx context.Context, p *CreateUsageMeterParams) (*UsageMeter, error) {
	resp, err := internal.Do[usageMeterResponse](ctx, c.cfg, "POST", "/usage-meters", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.UsageMeter, nil
}

// Get retrieves a usage meter by ID.
func (c *UsageMetersClient) Get(ctx context.Context, id string) (*UsageMeter, error) {
	resp, err := internal.Do[usageMeterResponse](ctx, c.cfg, "GET", "/usage-meters/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.UsageMeter, nil
}
