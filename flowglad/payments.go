//nolint:dupl // resource clients intentionally share structure
package flowglad

import (
	"context"
	"iter"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// PaymentsClient provides methods for the Flowglad Payments API.
type PaymentsClient struct {
	cfg *internal.HTTPConfig
}

// paymentResponse is the JSON envelope for a single-payment response.
type paymentResponse struct {
	Data struct {
		Payment Payment `json:"payment"`
	} `json:"data"`
}

// Get retrieves a payment by ID.
func (c *PaymentsClient) Get(ctx context.Context, id string) (*Payment, error) {
	resp, err := internal.Do[paymentResponse](ctx, c.cfg, "GET", "/payments/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Payment, nil
}

// List returns an iterator over all payments.
func (c *PaymentsClient) List(ctx context.Context, p *ListPaymentsParams) iter.Seq2[*Payment, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[Payment], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[Payment]](ctx, c.cfg, "GET", "/payments"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Refund initiates a refund for a payment.
func (c *PaymentsClient) Refund(ctx context.Context, id string, p *RefundPaymentParams) (*Payment, error) {
	resp, err := internal.Do[paymentResponse](ctx, c.cfg, "POST", "/payments/"+id+"/refund", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Payment, nil
}
