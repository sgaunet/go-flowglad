package flowglad

import (
	"context"
	"iter"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// ResourcesClient provides methods for the Flowglad Resources (Entitlements) API.
type ResourcesClient struct {
	cfg *internal.HTTPConfig
}

// resourceResponse is the JSON envelope for a single resource response.
type resourceResponse struct {
	Data struct {
		Resource Resource `json:"resource"`
	} `json:"data"`
}

// resourceClaimResponse is the JSON envelope for a resource claim response.
type resourceClaimResponse struct {
	Data struct {
		ResourceClaim ResourceClaim `json:"resourceClaim"`
	} `json:"data"`
}

// Create creates a new resource.
func (c *ResourcesClient) Create(ctx context.Context, p *CreateResourceParams) (*Resource, error) {
	resp, err := internal.Do[resourceResponse](ctx, c.cfg, "POST", "/resources", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Resource, nil
}

// Get retrieves a resource by ID.
func (c *ResourcesClient) Get(ctx context.Context, id string) (*Resource, error) {
	resp, err := internal.Do[resourceResponse](ctx, c.cfg, "GET", "/resources/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Resource, nil
}

// List returns an iterator over all resources.
func (c *ResourcesClient) List(ctx context.Context, p *ListResourcesParams) iter.Seq2[*Resource, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[Resource], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[Resource]](ctx, c.cfg, "GET", "/resources"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Update updates an existing resource.
func (c *ResourcesClient) Update(ctx context.Context, id string, p *UpdateResourceParams) (*Resource, error) {
	resp, err := internal.Do[resourceResponse](ctx, c.cfg, "PUT", "/resources/"+id, p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Resource, nil
}

// Claim claims a resource for a customer.
func (c *ResourcesClient) Claim(ctx context.Context, id string, p *ClaimResourceParams) (*ResourceClaim, error) {
	resp, err := internal.Do[resourceClaimResponse](ctx, c.cfg, "POST", "/resources/"+id+"/claim", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.ResourceClaim, nil
}

// Release releases a customer's claim on a resource.
func (c *ResourcesClient) Release(ctx context.Context, id string, p *ReleaseResourceParams) error {
	_, err := internal.Do[struct{}](ctx, c.cfg, "POST", "/resources/"+id+"/release", p, false)
	return err
}

// ListClaims returns an iterator over all claims on a resource.
func (c *ResourcesClient) ListClaims(ctx context.Context, id string) iter.Seq2[*ResourceClaim, error] {
	return listIter(func(cursor string) (*page[ResourceClaim], error) {
		qs := internal.BuildQueryString(cursor, nil)
		path := "/resources/" + id + "/claims" + qs
		resp, err := internal.Do[page[ResourceClaim]](ctx, c.cfg, "GET", path, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// ListUsages returns an iterator over all usage records for a resource.
func (c *ResourcesClient) ListUsages(ctx context.Context, id string) iter.Seq2[*ResourceUsage, error] {
	return listIter(func(cursor string) (*page[ResourceUsage], error) {
		qs := internal.BuildQueryString(cursor, nil)
		path := "/resources/" + id + "/usages" + qs
		resp, err := internal.Do[page[ResourceUsage]](ctx, c.cfg, "GET", path, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}
