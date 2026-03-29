package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func TestAPIKeys_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /api-keys/{id}", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": map[string]any{
			"apiKey": map[string]any{
				"id":             "key_1",
				"createdAt":      0,
				"livemode":       false,
				"organizationId": "org",
				"name":           "Production Key",
				"lastFour":       "k3y1",
				"active":         true,
			},
		},
	}))

	client := mustNewClient(t, srv.URL())
	key, err := client.APIKeys.Get(context.Background(), "key_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if key.ID != "key_1" {
		t.Errorf("ID: got %q, want key_1", key.ID)
	}
	if key.Name != "Production Key" {
		t.Errorf("Name: got %q, want Production Key", key.Name)
	}
	if !key.Active {
		t.Error("expected Active to be true")
	}
}
