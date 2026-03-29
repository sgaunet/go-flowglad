package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func resourceFixture(id, name, featureID string) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"resource": map[string]any{
				"id":             id,
				"createdAt":      0,
				"updatedAt":      0,
				"livemode":       false,
				"organizationId": "org",
				"name":           name,
				"featureId":      featureID,
			},
		},
	}
}

func TestResources_Create(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /resources", flowgladtest.RespondWith(http.StatusOK, resourceFixture("res_1", "My Resource", "feat_1")))

	client := mustNewClient(t, srv.URL())
	res, err := client.Resources.Create(context.Background(), &flowglad.CreateResourceParams{
		Name:      "My Resource",
		FeatureID: "feat_1",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if res.ID != "res_1" {
		t.Errorf("ID: got %q, want res_1", res.ID)
	}
}

func TestResources_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /resources/{id}", flowgladtest.RespondWith(http.StatusOK, resourceFixture("res_1", "My Resource", "feat_1")))

	client := mustNewClient(t, srv.URL())
	res, err := client.Resources.Get(context.Background(), "res_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if res.FeatureID != "feat_1" {
		t.Errorf("FeatureID: got %q, want feat_1", res.FeatureID)
	}
}

func TestResources_List(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /resources", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "res_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "name": "My Resource", "featureId": "feat_1"},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for res, err := range client.Resources.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, res.ID)
	}
	if len(ids) != 1 || ids[0] != "res_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}

func TestResources_Update(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("PUT /resources/{id}", flowgladtest.RespondWith(http.StatusOK, resourceFixture("res_1", "Updated Resource", "feat_1")))

	client := mustNewClient(t, srv.URL())
	res, err := client.Resources.Update(context.Background(), "res_1", &flowglad.UpdateResourceParams{
		Name: flowglad.Ptr("Updated Resource"),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if res.Name != "Updated Resource" {
		t.Errorf("Name: got %q, want Updated Resource", res.Name)
	}
}

func TestResources_Claim(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /resources/{id}/claim", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": map[string]any{
			"resourceClaim": map[string]any{
				"id":         "claim_1",
				"resourceId": "res_1",
				"customerId": "cus_1",
				"createdAt":  0,
			},
		},
	}))

	client := mustNewClient(t, srv.URL())
	claim, err := client.Resources.Claim(context.Background(), "res_1", &flowglad.ClaimResourceParams{
		CustomerID: "cus_1",
	})
	if err != nil {
		t.Fatalf("Claim: %v", err)
	}
	if claim.ID != "claim_1" {
		t.Errorf("ID: got %q, want claim_1", claim.ID)
	}
}

func TestResources_Release(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /resources/{id}/release", flowgladtest.RespondWith(http.StatusOK, map[string]any{}))

	client := mustNewClient(t, srv.URL())
	err := client.Resources.Release(context.Background(), "res_1", &flowglad.ReleaseResourceParams{
		CustomerID: "cus_1",
	})
	if err != nil {
		t.Fatalf("Release: %v", err)
	}
}

func TestResources_ListClaims(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /resources/{id}/claims", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "claim_1", "resourceId": "res_1", "customerId": "cus_1", "createdAt": 0},
			map[string]any{"id": "claim_2", "resourceId": "res_1", "customerId": "cus_2", "createdAt": 0},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 2,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for claim, err := range client.Resources.ListClaims(context.Background(), "res_1") {
		if err != nil {
			t.Fatalf("ListClaims: %v", err)
		}
		ids = append(ids, claim.ID)
	}
	if len(ids) != 2 || ids[0] != "claim_1" || ids[1] != "claim_2" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}

func TestResources_ListUsages(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /resources/{id}/usages", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "usage_1", "resourceId": "res_1", "quantity": 5.0, "createdAt": 0},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for usage, err := range client.Resources.ListUsages(context.Background(), "res_1") {
		if err != nil {
			t.Fatalf("ListUsages: %v", err)
		}
		ids = append(ids, usage.ID)
	}
	if len(ids) != 1 || ids[0] != "usage_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}
