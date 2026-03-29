package flowglad

import (
	"context"
	"iter"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// SubscriptionsClient provides methods for the Flowglad Subscriptions API.
type SubscriptionsClient struct {
	cfg *internal.HTTPConfig
}

// subscriptionResponse is the JSON envelope for a single-subscription response.
type subscriptionResponse struct {
	Data struct {
		Subscription Subscription `json:"subscription"`
	} `json:"data"`
}

// adjustmentPreviewResponse is the JSON envelope for the preview-adjust endpoint.
type adjustmentPreviewResponse struct {
	Data struct {
		AdjustmentPreview AdjustmentPreview `json:"adjustmentPreview"`
	} `json:"data"`
}

// Create creates a new subscription.
func (c *SubscriptionsClient) Create(ctx context.Context, p *CreateSubscriptionParams) (*Subscription, error) {
	resp, err := internal.Do[subscriptionResponse](ctx, c.cfg, "POST", "/subscriptions", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Subscription, nil
}

// Get retrieves a subscription by ID.
func (c *SubscriptionsClient) Get(ctx context.Context, id string) (*Subscription, error) {
	resp, err := internal.Do[subscriptionResponse](ctx, c.cfg, "GET", "/subscriptions/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Subscription, nil
}

// List returns an iterator over all subscriptions.
func (c *SubscriptionsClient) List(ctx context.Context, p *ListSubscriptionsParams) iter.Seq2[*Subscription, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[Subscription], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[Subscription]](ctx, c.cfg, "GET", "/subscriptions"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Adjust applies an adjustment (e.g. plan change, quantity change) to a subscription.
func (c *SubscriptionsClient) Adjust(ctx context.Context, id string, p *AdjustmentParams) (*Subscription, error) {
	resp, err := internal.Do[subscriptionResponse](ctx, c.cfg, "POST", "/subscriptions/"+id+"/adjust", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Subscription, nil
}

// PreviewAdjust previews the proration impact of an adjustment without applying it.
func (c *SubscriptionsClient) PreviewAdjust(
	ctx context.Context, id string, p *PreviewAdjustmentParams,
) (*AdjustmentPreview, error) {
	path := "/subscriptions/" + id + "/preview-adjust"
	resp, err := internal.Do[adjustmentPreviewResponse](ctx, c.cfg, "POST", path, p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.AdjustmentPreview, nil
}

// Cancel cancels a subscription. The timing is controlled by p.CancellationTiming.
func (c *SubscriptionsClient) Cancel(
	ctx context.Context, id string, p *CancelSubscriptionParams,
) (*Subscription, error) {
	resp, err := internal.Do[subscriptionResponse](ctx, c.cfg, "POST", "/subscriptions/"+id+"/cancel", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Subscription, nil
}

// Uncancel reverses a scheduled cancellation and restores the subscription to active.
func (c *SubscriptionsClient) Uncancel(ctx context.Context, id string) (*Subscription, error) {
	resp, err := internal.Do[subscriptionResponse](ctx, c.cfg, "POST", "/subscriptions/"+id+"/uncancel", nil, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Subscription, nil
}

// AddFeature adds a feature entitlement to a subscription.
func (c *SubscriptionsClient) AddFeature(ctx context.Context, id string, p *AddFeatureParams) (*Subscription, error) {
	path := "/subscriptions/" + id + "/add-feature"
	resp, err := internal.Do[subscriptionResponse](ctx, c.cfg, "POST", path, p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Subscription, nil
}

// CancelScheduledAdjustment cancels a pending scheduled adjustment on a subscription.
func (c *SubscriptionsClient) CancelScheduledAdjustment(ctx context.Context, id string) (*Subscription, error) {
	path := "/subscriptions/" + id + "/cancel-scheduled-adjustment"
	resp, err := internal.Do[subscriptionResponse](ctx, c.cfg, "POST", path, nil, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Subscription, nil
}
