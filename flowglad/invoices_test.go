package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func invoiceFixture(id, customerID, status string, amountDue int64) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"invoice": map[string]any{
				"id":             id,
				"createdAt":      0,
				"updatedAt":      0,
				"livemode":       false,
				"organizationId": "org",
				"customerId":     customerID,
				"status":         status,
				"amountDue":      amountDue,
				"amountPaid":     0,
				"currency":       "usd",
			},
		},
	}
}

func TestInvoices_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /invoices/{id}", flowgladtest.RespondWith(http.StatusOK, invoiceFixture("inv_1", "cus_1", "open", 5000)))

	client := mustNewClient(t, srv.URL())
	inv, err := client.Invoices.Get(context.Background(), "inv_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if inv.ID != "inv_1" {
		t.Errorf("ID: got %q, want inv_1", inv.ID)
	}
	if inv.Status != flowglad.InvoiceStatusOpen {
		t.Errorf("Status: got %q, want open", inv.Status)
	}
	if inv.AmountDue != 5000 {
		t.Errorf("AmountDue: got %d, want 5000", inv.AmountDue)
	}
}

func TestInvoices_List(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /invoices", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "inv_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "customerId": "cus_1", "status": "open", "amountDue": 5000, "amountPaid": 0, "currency": "usd"},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for inv, err := range client.Invoices.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, inv.ID)
	}
	if len(ids) != 1 || ids[0] != "inv_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}

func TestInvoiceLineItems_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /invoice-line-items/{id}", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": map[string]any{
			"invoiceLineItem": map[string]any{
				"id":          "li_1",
				"createdAt":   0,
				"updatedAt":   0,
				"invoiceId":   "inv_1",
				"description": "Pro Plan",
				"amount":      5000,
				"currency":    "usd",
			},
		},
	}))

	client := mustNewClient(t, srv.URL())
	li, err := client.InvoiceLineItems.Get(context.Background(), "li_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if li.ID != "li_1" {
		t.Errorf("ID: got %q, want li_1", li.ID)
	}
}

func TestInvoiceLineItems_List(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /invoice-line-items", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "li_1", "createdAt": 0, "updatedAt": 0, "invoiceId": "inv_1", "description": "Pro Plan", "amount": 5000, "currency": "usd"},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for li, err := range client.InvoiceLineItems.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, li.ID)
	}
	if len(ids) != 1 || ids[0] != "li_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}
