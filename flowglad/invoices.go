//nolint:dupl // resource clients intentionally share structure
package flowglad

import (
	"context"
	"iter"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// InvoicesClient provides methods for the Flowglad Invoices API.
type InvoicesClient struct {
	cfg *internal.HTTPConfig
}

// invoiceResponse is the JSON envelope for a single-invoice response.
type invoiceResponse struct {
	Data struct {
		Invoice Invoice `json:"invoice"`
	} `json:"data"`
}

// Get retrieves an invoice by ID.
func (c *InvoicesClient) Get(ctx context.Context, id string) (*Invoice, error) {
	resp, err := internal.Do[invoiceResponse](ctx, c.cfg, "GET", "/invoices/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Invoice, nil
}

// List returns an iterator over all invoices.
func (c *InvoicesClient) List(ctx context.Context, p *ListInvoicesParams) iter.Seq2[*Invoice, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[Invoice], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[Invoice]](ctx, c.cfg, "GET", "/invoices"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// InvoiceLineItemsClient provides methods for the Flowglad Invoice Line Items API.
type InvoiceLineItemsClient struct {
	cfg *internal.HTTPConfig
}

// invoiceLineItemResponse is the JSON envelope for a single line item response.
type invoiceLineItemResponse struct {
	Data struct {
		InvoiceLineItem InvoiceLineItem `json:"invoiceLineItem"`
	} `json:"data"`
}

// Get retrieves an invoice line item by ID.
func (c *InvoiceLineItemsClient) Get(ctx context.Context, id string) (*InvoiceLineItem, error) {
	path := "/invoice-line-items/" + id
	resp, err := internal.Do[invoiceLineItemResponse](ctx, c.cfg, "GET", path, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.InvoiceLineItem, nil
}

// List returns an iterator over all invoice line items.
func (c *InvoiceLineItemsClient) List(
	ctx context.Context, p *ListInvoiceLineItemsParams,
) iter.Seq2[*InvoiceLineItem, error] {
	var limit *int
	if p != nil {
		limit = p.Limit
	}
	return listIter(func(cursor string) (*page[InvoiceLineItem], error) {
		qs := internal.BuildQueryString(cursor, limit)
		resp, err := internal.Do[page[InvoiceLineItem]](ctx, c.cfg, "GET", "/invoice-line-items"+qs, nil, true)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}
