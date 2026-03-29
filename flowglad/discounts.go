//nolint:dupl // resource clients intentionally share structure
package flowglad

import (
	"context"
	"iter"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// DiscountsClient provides methods for the Flowglad Discounts API.
type DiscountsClient struct {
	cfg *internal.HTTPConfig
}

// discountResponse is the JSON envelope for a single-discount response.
type discountResponse struct {
	Data struct {
		Discount Discount `json:"discount"`
	} `json:"data"`
}

// Create creates a new discount.
func (c *DiscountsClient) Create(ctx context.Context, p *CreateDiscountParams) (*Discount, error) {
	resp, err := internal.Do[discountResponse](ctx, c.cfg, "POST", "/discounts", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Discount, nil
}

// Get retrieves a discount by ID.
func (c *DiscountsClient) Get(ctx context.Context, id string) (*Discount, error) {
	resp, err := internal.Do[discountResponse](ctx, c.cfg, "GET", "/discounts/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Discount, nil
}

// List returns an iterator over all discounts.
func (c *DiscountsClient) List(ctx context.Context, p *ListDiscountsParams) iter.Seq2[*Discount, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[Discount], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[Discount]](ctx, c.cfg, "GET", "/discounts"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Update updates an existing discount.
func (c *DiscountsClient) Update(ctx context.Context, id string, p *UpdateDiscountParams) (*Discount, error) {
	resp, err := internal.Do[discountResponse](ctx, c.cfg, "PUT", "/discounts/"+id, p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Discount, nil
}
