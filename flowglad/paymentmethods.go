//nolint:dupl // resource clients intentionally share structure
package flowglad

import (
	"context"
	"iter"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// PaymentMethodsClient provides methods for the Flowglad Payment Methods API.
type PaymentMethodsClient struct {
	cfg *internal.HTTPConfig
}

// paymentMethodResponse is the JSON envelope for a single payment method response.
type paymentMethodResponse struct {
	Data struct {
		PaymentMethod PaymentMethod `json:"paymentMethod"`
	} `json:"data"`
}

// Get retrieves a payment method by ID.
func (c *PaymentMethodsClient) Get(ctx context.Context, id string) (*PaymentMethod, error) {
	resp, err := internal.Do[paymentMethodResponse](ctx, c.cfg, "GET", "/payment-methods/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.PaymentMethod, nil
}

// List returns an iterator over all payment methods.
func (c *PaymentMethodsClient) List(ctx context.Context, p *ListPaymentMethodsParams) iter.Seq2[*PaymentMethod, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[PaymentMethod], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[PaymentMethod]](ctx, c.cfg, "GET", "/payment-methods"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}
