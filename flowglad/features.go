//nolint:dupl // resource clients intentionally share structure
package flowglad

import (
	"context"
	"iter"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// FeaturesClient provides methods for the Flowglad Features API.
type FeaturesClient struct {
	cfg *internal.HTTPConfig
}

// featureResponse is the JSON envelope for a single feature response.
type featureResponse struct {
	Data struct {
		Feature Feature `json:"feature"`
	} `json:"data"`
}

// Create creates a new feature.
func (c *FeaturesClient) Create(ctx context.Context, p *CreateFeatureParams) (*Feature, error) {
	resp, err := internal.Do[featureResponse](ctx, c.cfg, "POST", "/features", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Feature, nil
}

// Get retrieves a feature by ID.
func (c *FeaturesClient) Get(ctx context.Context, id string) (*Feature, error) {
	resp, err := internal.Do[featureResponse](ctx, c.cfg, "GET", "/features/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Feature, nil
}

// List returns an iterator over all features.
func (c *FeaturesClient) List(ctx context.Context, p *ListFeaturesParams) iter.Seq2[*Feature, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[Feature], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[Feature]](ctx, c.cfg, "GET", "/features"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Update updates an existing feature.
func (c *FeaturesClient) Update(ctx context.Context, id string, p *UpdateFeatureParams) (*Feature, error) {
	resp, err := internal.Do[featureResponse](ctx, c.cfg, "PUT", "/features/"+id, p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Feature, nil
}

// AddProductFeature adds a feature to a product.
func (c *FeaturesClient) AddProductFeature(
	ctx context.Context, productID string, p *AddProductFeatureParams,
) (*Feature, error) {
	resp, err := internal.Do[featureResponse](ctx, c.cfg, "POST", "/products/"+productID+"/features", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Feature, nil
}

// ListSubscriptionFeatures returns an iterator over features for a subscription.
func (c *FeaturesClient) ListSubscriptionFeatures(
	ctx context.Context, subscriptionID string,
) iter.Seq2[*Feature, error] {
	return listIter(func(cursor string) (*page[Feature], error) {
		qs := internal.BuildQueryString(cursor, nil)
		path := "/subscriptions/" + subscriptionID + "/features" + qs
		resp, err := internal.Do[page[Feature]](ctx, c.cfg, "GET", path, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}
