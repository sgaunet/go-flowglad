//nolint:dupl // resource clients intentionally share structure
package flowglad

import (
	"context"
	"iter"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// ProductsClient provides methods for the Flowglad Products API.
type ProductsClient struct {
	cfg *internal.HTTPConfig
}

// productResponse is the JSON envelope for a single-product response.
type productResponse struct {
	Data struct {
		Product Product `json:"product"`
	} `json:"data"`
}

// Create creates a new product.
func (c *ProductsClient) Create(ctx context.Context, p *CreateProductParams) (*Product, error) {
	resp, err := internal.Do[productResponse](ctx, c.cfg, "POST", "/products", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Product, nil
}

// Get retrieves a product by ID.
func (c *ProductsClient) Get(ctx context.Context, id string) (*Product, error) {
	resp, err := internal.Do[productResponse](ctx, c.cfg, "GET", "/products/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Product, nil
}

// List returns an iterator over all products.
func (c *ProductsClient) List(ctx context.Context, p *ListProductsParams) iter.Seq2[*Product, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[Product], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[Product]](ctx, c.cfg, "GET", "/products"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Update updates an existing product.
func (c *ProductsClient) Update(ctx context.Context, id string, p *UpdateProductParams) (*Product, error) {
	resp, err := internal.Do[productResponse](ctx, c.cfg, "PUT", "/products/"+id, p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Product, nil
}
