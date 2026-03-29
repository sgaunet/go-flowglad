package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func paymentFixture(id, customerID, status string, amount int64) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"payment": map[string]any{
				"id":             id,
				"createdAt":      0,
				"updatedAt":      0,
				"livemode":       false,
				"organizationId": "org",
				"customerId":     customerID,
				"amount":         amount,
				"currency":       "usd",
				"status":         status,
			},
		},
	}
}

func TestPayments_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /payments/{id}", flowgladtest.RespondWith(http.StatusOK, paymentFixture("pay_1", "cus_1", "succeeded", 5000)))

	client := mustNewClient(t, srv.URL())
	pay, err := client.Payments.Get(context.Background(), "pay_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if pay.ID != "pay_1" {
		t.Errorf("ID: got %q, want pay_1", pay.ID)
	}
	if pay.Status != "succeeded" {
		t.Errorf("Status: got %q, want succeeded", pay.Status)
	}
}

func TestPayments_List(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /payments", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "pay_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "customerId": "cus_1", "amount": 5000, "currency": "usd", "status": "succeeded"},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for pay, err := range client.Payments.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, pay.ID)
	}
	if len(ids) != 1 || ids[0] != "pay_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}

func TestPayments_Refund(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /payments/{id}/refund", flowgladtest.RespondWith(http.StatusOK, paymentFixture("pay_1", "cus_1", "refunded", 5000)))

	client := mustNewClient(t, srv.URL())
	pay, err := client.Payments.Refund(context.Background(), "pay_1", &flowglad.RefundPaymentParams{})
	if err != nil {
		t.Fatalf("Refund: %v", err)
	}
	if pay.Status != "refunded" {
		t.Errorf("Status: got %q, want refunded", pay.Status)
	}
}
