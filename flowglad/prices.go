//nolint:dupl // resource clients intentionally share structure
package flowglad

import (
	"context"
	"iter"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// PricesClient provides methods for the Flowglad Prices API.
type PricesClient struct {
	cfg *internal.HTTPConfig
}

// priceResponse is the JSON envelope for a single-price response.
type priceResponse struct {
	Data struct {
		Price Price `json:"price"`
	} `json:"data"`
}

// Create creates a new price.
func (c *PricesClient) Create(ctx context.Context, p *CreatePriceParams) (*Price, error) {
	resp, err := internal.Do[priceResponse](ctx, c.cfg, "POST", "/prices", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Price, nil
}

// Get retrieves a price by ID.
func (c *PricesClient) Get(ctx context.Context, id string) (*Price, error) {
	resp, err := internal.Do[priceResponse](ctx, c.cfg, "GET", "/prices/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Price, nil
}

// List returns an iterator over all prices.
func (c *PricesClient) List(ctx context.Context, p *ListPricesParams) iter.Seq2[*Price, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[Price], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[Price]](ctx, c.cfg, "GET", "/prices"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Update updates an existing price.
func (c *PricesClient) Update(ctx context.Context, id string, p *UpdatePriceParams) (*Price, error) {
	resp, err := internal.Do[priceResponse](ctx, c.cfg, "PUT", "/prices/"+id, p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Price, nil
}
