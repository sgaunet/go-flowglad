//nolint:dupl // resource clients intentionally share structure
package flowglad

import (
	"context"
	"iter"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// CheckoutSessionsClient provides methods for the Flowglad Checkout Sessions API.
type CheckoutSessionsClient struct {
	cfg *internal.HTTPConfig
}

// checkoutSessionResponse is the JSON envelope for a single checkout session response.
type checkoutSessionResponse struct {
	Data struct {
		CheckoutSession CheckoutSession `json:"checkoutSession"`
	} `json:"data"`
}

// Create creates a new checkout session and returns the hosted URL.
func (c *CheckoutSessionsClient) Create(ctx context.Context, p *CreateCheckoutSessionParams) (*CheckoutSession, error) {
	resp, err := internal.Do[checkoutSessionResponse](ctx, c.cfg, "POST", "/checkout-sessions", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.CheckoutSession, nil
}

// Get retrieves a checkout session by ID.
func (c *CheckoutSessionsClient) Get(ctx context.Context, id string) (*CheckoutSession, error) {
	path := "/checkout-sessions/" + id
	resp, err := internal.Do[checkoutSessionResponse](ctx, c.cfg, "GET", path, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.CheckoutSession, nil
}

// List returns an iterator over all checkout sessions.
func (c *CheckoutSessionsClient) List(
	ctx context.Context, p *ListCheckoutSessionsParams,
) iter.Seq2[*CheckoutSession, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[CheckoutSession], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[CheckoutSession]](ctx, c.cfg, "GET", "/checkout-sessions"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}
