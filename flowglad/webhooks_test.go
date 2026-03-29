package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func webhookFixture(id, name, url string, active bool) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"webhook": map[string]any{
				"id":             id,
				"createdAt":      0,
				"updatedAt":      0,
				"livemode":       false,
				"organizationId": "org",
				"name":           name,
				"url":            url,
				"filterTypes":    []any{},
				"active":         active,
			},
		},
	}
}

func TestWebhooks_Create(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /webhooks", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": map[string]any{
			"webhook": map[string]any{
				"id": "wh_1", "createdAt": 0, "updatedAt": 0, "livemode": false,
				"organizationId": "org", "name": "My Webhook",
				"url": "https://example.com/wh", "filterTypes": []any{}, "active": true,
				"secret": "whsec_abc123",
			},
		},
	}))

	client := mustNewClient(t, srv.URL())
	result, err := client.Webhooks.Create(context.Background(), &flowglad.CreateWebhookParams{
		Name: "My Webhook",
		URL:  "https://example.com/wh",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if result.ID != "wh_1" {
		t.Errorf("ID: got %q, want wh_1", result.ID)
	}
	if result.Secret != "whsec_abc123" {
		t.Errorf("Secret: got %q, want whsec_abc123", result.Secret)
	}
}

func TestWebhooks_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /webhooks/{id}", flowgladtest.RespondWith(http.StatusOK, webhookFixture("wh_1", "My Webhook", "https://example.com/wh", true)))

	client := mustNewClient(t, srv.URL())
	wh, err := client.Webhooks.Get(context.Background(), "wh_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if wh.ID != "wh_1" {
		t.Errorf("ID: got %q, want wh_1", wh.ID)
	}
}

func TestWebhooks_Update(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("PUT /webhooks/{id}", flowgladtest.RespondWith(http.StatusOK, webhookFixture("wh_1", "Updated", "https://example.com/wh2", true)))

	client := mustNewClient(t, srv.URL())
	wh, err := client.Webhooks.Update(context.Background(), "wh_1", &flowglad.UpdateWebhookParams{
		Name: flowglad.Ptr("Updated"),
		URL:  flowglad.Ptr("https://example.com/wh2"),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if wh.Name != "Updated" {
		t.Errorf("Name: got %q, want Updated", wh.Name)
	}
}
