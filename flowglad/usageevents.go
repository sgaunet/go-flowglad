package flowglad

import (
	"context"
	"iter"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// UsageEventsClient provides methods for the Flowglad Usage Events API.
type UsageEventsClient struct {
	cfg *internal.HTTPConfig
}

// usageEventResponse is the JSON envelope for a single usage event response.
type usageEventResponse struct {
	Data struct {
		UsageEvent UsageEvent `json:"usageEvent"`
	} `json:"data"`
}

// bulkUsageEventsResponse is the JSON envelope for a bulk usage event response.
type bulkUsageEventsResponse struct {
	Data struct {
		UsageEvents []UsageEvent `json:"usageEvents"`
	} `json:"data"`
}

// Create records a single usage event against a meter.
func (c *UsageEventsClient) Create(ctx context.Context, p *CreateUsageEventParams) (*UsageEvent, error) {
	resp, err := internal.Do[usageEventResponse](ctx, c.cfg, "POST", "/usage-events", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.UsageEvent, nil
}

// BulkCreate records multiple usage events atomically.
func (c *UsageEventsClient) BulkCreate(ctx context.Context, p *BulkCreateUsageEventsParams) ([]*UsageEvent, error) {
	resp, err := internal.Do[bulkUsageEventsResponse](ctx, c.cfg, "POST", "/usage-events/bulk", p, false)
	if err != nil {
		return nil, err
	}
	result := make([]*UsageEvent, len(resp.Data.UsageEvents))
	for i := range resp.Data.UsageEvents {
		result[i] = &resp.Data.UsageEvents[i]
	}
	return result, nil
}

// Get retrieves a usage event by ID.
func (c *UsageEventsClient) Get(ctx context.Context, id string) (*UsageEvent, error) {
	resp, err := internal.Do[usageEventResponse](ctx, c.cfg, "GET", "/usage-events/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.UsageEvent, nil
}

// List returns an iterator over all usage events.
func (c *UsageEventsClient) List(ctx context.Context, p *ListUsageEventsParams) iter.Seq2[*UsageEvent, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[UsageEvent], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[UsageEvent]](ctx, c.cfg, "GET", "/usage-events"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}
