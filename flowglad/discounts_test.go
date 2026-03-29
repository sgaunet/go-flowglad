package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func discountFixture(id, name string, active bool) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"discount": map[string]any{
				"id":             id,
				"createdAt":      0,
				"updatedAt":      0,
				"livemode":       false,
				"organizationId": "org",
				"name":           name,
				"duration":       "once",
				"active":         active,
			},
		},
	}
}

func TestDiscounts_Create(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /discounts", flowgladtest.RespondWith(http.StatusOK, discountFixture("disc_1", "50% Off", true)))

	client := mustNewClient(t, srv.URL())
	disc, err := client.Discounts.Create(context.Background(), &flowglad.CreateDiscountParams{
		Name:       "50% Off",
		PercentOff: flowglad.Ptr(50.0),
		Duration:   "once",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if disc.ID != "disc_1" {
		t.Errorf("ID: got %q, want disc_1", disc.ID)
	}
}

func TestDiscounts_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /discounts/{id}", flowgladtest.RespondWith(http.StatusOK, discountFixture("disc_1", "50% Off", true)))

	client := mustNewClient(t, srv.URL())
	disc, err := client.Discounts.Get(context.Background(), "disc_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if disc.Name != "50% Off" {
		t.Errorf("Name: got %q, want '50%% Off'", disc.Name)
	}
}

func TestDiscounts_List(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /discounts", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "disc_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "name": "50% Off", "duration": "once", "active": true},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for disc, err := range client.Discounts.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, disc.ID)
	}
	if len(ids) != 1 || ids[0] != "disc_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}

func TestDiscounts_Update(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("PUT /discounts/{id}", flowgladtest.RespondWith(http.StatusOK, discountFixture("disc_1", "25% Off", true)))

	client := mustNewClient(t, srv.URL())
	disc, err := client.Discounts.Update(context.Background(), "disc_1", &flowglad.UpdateDiscountParams{
		Name: flowglad.Ptr("25% Off"),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if disc.Name != "25% Off" {
		t.Errorf("Name: got %q, want '25%% Off'", disc.Name)
	}
}
