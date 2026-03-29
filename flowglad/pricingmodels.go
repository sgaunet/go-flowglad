//nolint:dupl // resource clients intentionally share structure
package flowglad

import (
	"context"
	"iter"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// PricingModelsClient provides methods for the Flowglad Pricing Models API.
type PricingModelsClient struct {
	cfg *internal.HTTPConfig
}

// pricingModelResponse is the JSON envelope for a single pricing model response.
type pricingModelResponse struct {
	Data struct {
		PricingModel PricingModel `json:"pricingModel"`
	} `json:"data"`
}

// Create creates a new pricing model.
func (c *PricingModelsClient) Create(ctx context.Context, p *CreatePricingModelParams) (*PricingModel, error) {
	resp, err := internal.Do[pricingModelResponse](ctx, c.cfg, "POST", "/pricing-models", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.PricingModel, nil
}

// Get retrieves a pricing model by ID.
func (c *PricingModelsClient) Get(ctx context.Context, id string) (*PricingModel, error) {
	resp, err := internal.Do[pricingModelResponse](ctx, c.cfg, "GET", "/pricing-models/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.PricingModel, nil
}

// List returns an iterator over all pricing models.
func (c *PricingModelsClient) List(ctx context.Context, p *ListPricingModelsParams) iter.Seq2[*PricingModel, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[PricingModel], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[PricingModel]](ctx, c.cfg, "GET", "/pricing-models"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Update updates an existing pricing model.
func (c *PricingModelsClient) Update(
	ctx context.Context, id string, p *UpdatePricingModelParams,
) (*PricingModel, error) {
	resp, err := internal.Do[pricingModelResponse](ctx, c.cfg, "PUT", "/pricing-models/"+id, p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.PricingModel, nil
}

// Clone duplicates an existing pricing model.
func (c *PricingModelsClient) Clone(ctx context.Context, id string, p *ClonePricingModelParams) (*PricingModel, error) {
	resp, err := internal.Do[pricingModelResponse](ctx, c.cfg, "POST", "/pricing-models/"+id+"/clone", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.PricingModel, nil
}

// Export exports a pricing model's configuration.
func (c *PricingModelsClient) Export(ctx context.Context, id string) (*PricingModel, error) {
	resp, err := internal.Do[pricingModelResponse](ctx, c.cfg, "POST", "/pricing-models/"+id+"/export", nil, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.PricingModel, nil
}

// Setup creates a pricing model using the setup workflow.
func (c *PricingModelsClient) Setup(ctx context.Context, p *SetupPricingModelParams) (*PricingModel, error) {
	resp, err := internal.Do[pricingModelResponse](ctx, c.cfg, "POST", "/pricing-models/setup", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.PricingModel, nil
}

// GetDefault retrieves the organization's default pricing model.
func (c *PricingModelsClient) GetDefault(ctx context.Context) (*PricingModel, error) {
	resp, err := internal.Do[pricingModelResponse](ctx, c.cfg, "GET", "/pricing-models/default", nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.PricingModel, nil
}
