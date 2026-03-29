package flowglad

import (
	"context"
	"iter"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// CustomersClient provides methods for the Flowglad Customers API.
type CustomersClient struct {
	cfg *internal.HTTPConfig
}

// customerResponse is the JSON envelope for a single-customer response.
type customerResponse struct {
	Data struct {
		Customer Customer `json:"customer"`
	} `json:"data"`
}

// Create creates a new customer.
func (c *CustomersClient) Create(ctx context.Context, p *CreateCustomerParams) (*Customer, error) {
	resp, err := internal.Do[customerResponse](ctx, c.cfg, "POST", "/customers", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Customer, nil
}

// Get retrieves a customer by ID.
func (c *CustomersClient) Get(ctx context.Context, id string) (*Customer, error) {
	resp, err := internal.Do[customerResponse](ctx, c.cfg, "GET", "/customers/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Customer, nil
}

// List returns an iterator over all customers. It transparently fetches
// subsequent pages as the caller iterates.
func (c *CustomersClient) List(ctx context.Context, p *ListCustomersParams) iter.Seq2[*Customer, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[Customer], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[Customer]](ctx, c.cfg, "GET", "/customers"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Update updates an existing customer.
func (c *CustomersClient) Update(ctx context.Context, id string, p *UpdateCustomerParams) (*Customer, error) {
	resp, err := internal.Do[customerResponse](ctx, c.cfg, "PUT", "/customers/"+id, p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Customer, nil
}

// Archive archives a customer, preventing new charges.
func (c *CustomersClient) Archive(ctx context.Context, id string) (*Customer, error) {
	resp, err := internal.Do[customerResponse](ctx, c.cfg, "POST", "/customers/"+id+"/archive", nil, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Customer, nil
}

// billingDetailsResponse is the JSON envelope for the billing details endpoint.
type billingDetailsResponse struct {
	Data CustomerBillingDetails `json:"data"`
}

// GetBillingDetails retrieves billing details (address + default payment method) for a customer.
func (c *CustomersClient) GetBillingDetails(ctx context.Context, id string) (*CustomerBillingDetails, error) {
	path := "/customers/" + id + "/billing-details"
	resp, err := internal.Do[billingDetailsResponse](ctx, c.cfg, "GET", path, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// usageBalancesResponse is the JSON envelope for the usage balances endpoint.
type usageBalancesResponse struct {
	Data CustomerUsageBalances `json:"data"`
}

// GetUsageBalances retrieves all usage meter balances for a customer.
func (c *CustomersClient) GetUsageBalances(ctx context.Context, id string) (*CustomerUsageBalances, error) {
	path := "/customers/" + id + "/usage-balances"
	resp, err := internal.Do[usageBalancesResponse](ctx, c.cfg, "GET", path, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
