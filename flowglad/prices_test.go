package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func priceFixture(id, name, productID string, amount int64, active bool) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"price": map[string]any{
				"id":             id,
				"createdAt":      0,
				"updatedAt":      0,
				"livemode":       false,
				"organizationId": "org",
				"productId":      productID,
				"name":           name,
				"currency":       "usd",
				"unitAmount":     amount,
				"active":         active,
			},
		},
	}
}

func TestPrices_Create(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /prices", flowgladtest.RespondWith(http.StatusOK, priceFixture("price_1", "Monthly", "prod_1", 1000, true)))

	client := mustNewClient(t, srv.URL())
	price, err := client.Prices.Create(context.Background(), &flowglad.CreatePriceParams{
		Name:       "Monthly",
		ProductID:  "prod_1",
		Currency:   "usd",
		UnitAmount: 1000,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if price.ID != "price_1" {
		t.Errorf("ID: got %q, want price_1", price.ID)
	}
	if price.UnitAmount != 1000 {
		t.Errorf("UnitAmount: got %d, want 1000", price.UnitAmount)
	}
}

func TestPrices_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /prices/{id}", flowgladtest.RespondWith(http.StatusOK, priceFixture("price_1", "Monthly", "prod_1", 1000, true)))

	client := mustNewClient(t, srv.URL())
	price, err := client.Prices.Get(context.Background(), "price_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if price.Currency != "usd" {
		t.Errorf("Currency: got %q, want usd", price.Currency)
	}
}

func TestPrices_List(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /prices", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "price_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "productId": "prod_1", "name": "Monthly", "currency": "usd", "unitAmount": 1000, "active": true},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for price, err := range client.Prices.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, price.ID)
	}
	if len(ids) != 1 || ids[0] != "price_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}

func TestPrices_Update(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("PUT /prices/{id}", flowgladtest.RespondWith(http.StatusOK, priceFixture("price_1", "Annual", "prod_1", 10000, true)))

	client := mustNewClient(t, srv.URL())
	price, err := client.Prices.Update(context.Background(), "price_1", &flowglad.UpdatePriceParams{
		Name: flowglad.Ptr("Annual"),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if price.Name != "Annual" {
		t.Errorf("Name: got %q, want Annual", price.Name)
	}
}
