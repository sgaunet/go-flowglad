package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func pricingModelFixture(id, name string, isDefault bool) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"pricingModel": map[string]any{
				"id":             id,
				"createdAt":      0,
				"updatedAt":      0,
				"livemode":       false,
				"organizationId": "org",
				"name":           name,
				"isDefault":      isDefault,
			},
		},
	}
}

func TestPricingModels_Create(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /pricing-models", flowgladtest.RespondWith(http.StatusOK, pricingModelFixture("pm_1", "Standard", false)))

	client := mustNewClient(t, srv.URL())
	pm, err := client.PricingModels.Create(context.Background(), &flowglad.CreatePricingModelParams{
		Name: "Standard",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if pm.ID != "pm_1" {
		t.Errorf("ID: got %q, want pm_1", pm.ID)
	}
}

func TestPricingModels_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /pricing-models/{id}", flowgladtest.RespondWith(http.StatusOK, pricingModelFixture("pm_1", "Standard", false)))

	client := mustNewClient(t, srv.URL())
	pm, err := client.PricingModels.Get(context.Background(), "pm_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if pm.Name != "Standard" {
		t.Errorf("Name: got %q, want Standard", pm.Name)
	}
}

func TestPricingModels_List(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /pricing-models", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "pm_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "name": "Standard", "isDefault": false},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for pm, err := range client.PricingModels.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, pm.ID)
	}
	if len(ids) != 1 || ids[0] != "pm_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}

func TestPricingModels_Update(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("PUT /pricing-models/{id}", flowgladtest.RespondWith(http.StatusOK, pricingModelFixture("pm_1", "Enterprise", false)))

	client := mustNewClient(t, srv.URL())
	pm, err := client.PricingModels.Update(context.Background(), "pm_1", &flowglad.UpdatePricingModelParams{
		Name: flowglad.Ptr("Enterprise"),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if pm.Name != "Enterprise" {
		t.Errorf("Name: got %q, want Enterprise", pm.Name)
	}
}

func TestPricingModels_Clone(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /pricing-models/{id}/clone", flowgladtest.RespondWith(http.StatusOK, pricingModelFixture("pm_2", "Standard (Copy)", false)))

	client := mustNewClient(t, srv.URL())
	pm, err := client.PricingModels.Clone(context.Background(), "pm_1", &flowglad.ClonePricingModelParams{
		Name: flowglad.Ptr("Standard (Copy)"),
	})
	if err != nil {
		t.Fatalf("Clone: %v", err)
	}
	if pm.ID != "pm_2" {
		t.Errorf("ID: got %q, want pm_2", pm.ID)
	}
	if pm.Name != "Standard (Copy)" {
		t.Errorf("Name: got %q, want Standard (Copy)", pm.Name)
	}
}

func TestPricingModels_Export(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /pricing-models/{id}/export", flowgladtest.RespondWith(http.StatusOK, pricingModelFixture("pm_1", "Standard", false)))

	client := mustNewClient(t, srv.URL())
	pm, err := client.PricingModels.Export(context.Background(), "pm_1")
	if err != nil {
		t.Fatalf("Export: %v", err)
	}
	if pm.ID != "pm_1" {
		t.Errorf("ID: got %q, want pm_1", pm.ID)
	}
}

func TestPricingModels_Setup(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /pricing-models/setup", flowgladtest.RespondWith(http.StatusOK, pricingModelFixture("pm_3", "New Setup", true)))

	client := mustNewClient(t, srv.URL())
	pm, err := client.PricingModels.Setup(context.Background(), &flowglad.SetupPricingModelParams{
		Name: "New Setup",
	})
	if err != nil {
		t.Fatalf("Setup: %v", err)
	}
	if pm.ID != "pm_3" {
		t.Errorf("ID: got %q, want pm_3", pm.ID)
	}
	if !pm.IsDefault {
		t.Error("expected IsDefault to be true")
	}
}

func TestPricingModels_GetDefault(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /pricing-models/default", flowgladtest.RespondWith(http.StatusOK, pricingModelFixture("pm_1", "Standard", true)))

	client := mustNewClient(t, srv.URL())
	pm, err := client.PricingModels.GetDefault(context.Background())
	if err != nil {
		t.Fatalf("GetDefault: %v", err)
	}
	if pm.ID != "pm_1" {
		t.Errorf("ID: got %q, want pm_1", pm.ID)
	}
	if !pm.IsDefault {
		t.Error("expected IsDefault to be true")
	}
}
