package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func productFixture(id, name string, active bool) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"product": map[string]any{
				"id":             id,
				"createdAt":      0,
				"updatedAt":      0,
				"livemode":       false,
				"organizationId": "org",
				"name":           name,
				"active":         active,
			},
		},
	}
}

func TestProducts_Create(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /products", flowgladtest.RespondWith(http.StatusOK, productFixture("prod_1", "Pro Plan", true)))

	client := mustNewClient(t, srv.URL())
	prod, err := client.Products.Create(context.Background(), &flowglad.CreateProductParams{
		Name: "Pro Plan",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if prod.ID != "prod_1" {
		t.Errorf("ID: got %q, want prod_1", prod.ID)
	}
}

func TestProducts_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /products/{id}", flowgladtest.RespondWith(http.StatusOK, productFixture("prod_1", "Pro Plan", true)))

	client := mustNewClient(t, srv.URL())
	prod, err := client.Products.Get(context.Background(), "prod_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if prod.Name != "Pro Plan" {
		t.Errorf("Name: got %q, want Pro Plan", prod.Name)
	}
}

func TestProducts_List(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /products", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "prod_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "name": "Pro Plan", "active": true},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for prod, err := range client.Products.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, prod.ID)
	}
	if len(ids) != 1 || ids[0] != "prod_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}

func TestProducts_Update(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("PUT /products/{id}", flowgladtest.RespondWith(http.StatusOK, productFixture("prod_1", "Updated Plan", true)))

	client := mustNewClient(t, srv.URL())
	prod, err := client.Products.Update(context.Background(), "prod_1", &flowglad.UpdateProductParams{
		Name: flowglad.Ptr("Updated Plan"),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if prod.Name != "Updated Plan" {
		t.Errorf("Name: got %q, want Updated Plan", prod.Name)
	}
}

func TestProducts_Create_Error(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /products", flowgladtest.RespondWithError(http.StatusBadRequest, "validation", "name required"))

	client := mustNewClient(t, srv.URL())
	_, err := client.Products.Create(context.Background(), &flowglad.CreateProductParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestProducts_Get_Error(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /products/{id}", flowgladtest.RespondWithError(http.StatusNotFound, "not_found", "not found"))

	client := mustNewClient(t, srv.URL())
	_, err := client.Products.Get(context.Background(), "prod_gone")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestProducts_Update_Error(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("PUT /products/{id}", flowgladtest.RespondWithError(http.StatusNotFound, "not_found", "not found"))

	client := mustNewClient(t, srv.URL())
	_, err := client.Products.Update(context.Background(), "prod_gone", &flowglad.UpdateProductParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}
