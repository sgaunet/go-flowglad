package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func featureFixture(id, name, slug string) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"feature": map[string]any{
				"id":             id,
				"createdAt":      0,
				"updatedAt":      0,
				"livemode":       false,
				"organizationId": "org",
				"name":           name,
				"slug":           slug,
			},
		},
	}
}

func TestFeatures_Create(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /features", flowgladtest.RespondWith(http.StatusOK, featureFixture("feat_1", "API Access", "api-access")))

	client := mustNewClient(t, srv.URL())
	feat, err := client.Features.Create(context.Background(), &flowglad.CreateFeatureParams{
		Name: "API Access",
		Slug: "api-access",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if feat.ID != "feat_1" {
		t.Errorf("ID: got %q, want feat_1", feat.ID)
	}
	if feat.Slug != "api-access" {
		t.Errorf("Slug: got %q, want api-access", feat.Slug)
	}
}

func TestFeatures_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /features/{id}", flowgladtest.RespondWith(http.StatusOK, featureFixture("feat_1", "API Access", "api-access")))

	client := mustNewClient(t, srv.URL())
	feat, err := client.Features.Get(context.Background(), "feat_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if feat.Name != "API Access" {
		t.Errorf("Name: got %q, want API Access", feat.Name)
	}
}

func TestFeatures_List(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /features", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "feat_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "name": "API Access", "slug": "api-access"},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for feat, err := range client.Features.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, feat.ID)
	}
	if len(ids) != 1 || ids[0] != "feat_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}

func TestFeatures_Update(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("PUT /features/{id}", flowgladtest.RespondWith(http.StatusOK, featureFixture("feat_1", "Advanced API", "api-access")))

	client := mustNewClient(t, srv.URL())
	feat, err := client.Features.Update(context.Background(), "feat_1", &flowglad.UpdateFeatureParams{
		Name: flowglad.Ptr("Advanced API"),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if feat.Name != "Advanced API" {
		t.Errorf("Name: got %q, want Advanced API", feat.Name)
	}
}

func TestFeatures_AddProductFeature(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /products/{id}/features", flowgladtest.RespondWith(http.StatusOK, featureFixture("feat_1", "API Access", "api-access")))

	client := mustNewClient(t, srv.URL())
	feat, err := client.Features.AddProductFeature(context.Background(), "prod_1", &flowglad.AddProductFeatureParams{
		FeatureID: "feat_1",
	})
	if err != nil {
		t.Fatalf("AddProductFeature: %v", err)
	}
	if feat.ID != "feat_1" {
		t.Errorf("ID: got %q, want feat_1", feat.ID)
	}
}

func TestFeatures_ListSubscriptionFeatures(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /subscriptions/{id}/features", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "feat_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "name": "API Access", "slug": "api-access"},
			map[string]any{"id": "feat_2", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "name": "SSO", "slug": "sso"},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 2,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for feat, err := range client.Features.ListSubscriptionFeatures(context.Background(), "sub_1") {
		if err != nil {
			t.Fatalf("ListSubscriptionFeatures: %v", err)
		}
		ids = append(ids, feat.ID)
	}
	if len(ids) != 2 || ids[0] != "feat_1" || ids[1] != "feat_2" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}
