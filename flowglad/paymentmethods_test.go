package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func paymentMethodFixture(id, customerID, pmType string) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"paymentMethod": map[string]any{
				"id":             id,
				"createdAt":      0,
				"updatedAt":      0,
				"livemode":       false,
				"organizationId": "org",
				"customerId":     customerID,
				"type":           pmType,
				"last4":          "4242",
				"brand":          "visa",
			},
		},
	}
}

func TestPaymentMethods_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /payment-methods/{id}", flowgladtest.RespondWith(http.StatusOK, paymentMethodFixture("pm_1", "cus_1", "card")))

	client := mustNewClient(t, srv.URL())
	pm, err := client.PaymentMethods.Get(context.Background(), "pm_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if pm.ID != "pm_1" {
		t.Errorf("ID: got %q, want pm_1", pm.ID)
	}
	if pm.Type != "card" {
		t.Errorf("Type: got %q, want card", pm.Type)
	}
}

func TestPaymentMethods_List(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /payment-methods", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "pm_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "customerId": "cus_1", "type": "card"},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for pm, err := range client.PaymentMethods.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, pm.ID)
	}
	if len(ids) != 1 || ids[0] != "pm_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}
