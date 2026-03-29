package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func usageEventFixture(id, meterID, customerID string, quantity float64) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"usageEvent": map[string]any{
				"id":             id,
				"createdAt":      0,
				"livemode":       false,
				"organizationId": "org",
				"usageMeterId":   meterID,
				"customerId":     customerID,
				"quantity":       quantity,
			},
		},
	}
}

func TestUsageEvents_Create(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /usage-events", flowgladtest.RespondWith(http.StatusOK, usageEventFixture("ue_1", "um_1", "cus_1", 10.0)))

	client := mustNewClient(t, srv.URL())
	ev, err := client.UsageEvents.Create(context.Background(), &flowglad.CreateUsageEventParams{
		UsageMeterID: "um_1",
		CustomerID:   "cus_1",
		Quantity:     10.0,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if ev.ID != "ue_1" {
		t.Errorf("ID: got %q, want ue_1", ev.ID)
	}
	if ev.Quantity != 10.0 {
		t.Errorf("Quantity: got %v, want 10.0", ev.Quantity)
	}
}

func TestUsageEvents_BulkCreate(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /usage-events/bulk", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": map[string]any{
			"usageEvents": []any{
				map[string]any{"id": "ue_1", "createdAt": 0, "livemode": false, "organizationId": "org", "usageMeterId": "um_1", "customerId": "cus_1", "quantity": 10.0},
				map[string]any{"id": "ue_2", "createdAt": 0, "livemode": false, "organizationId": "org", "usageMeterId": "um_1", "customerId": "cus_1", "quantity": 20.0},
			},
		},
	}))

	client := mustNewClient(t, srv.URL())
	events, err := client.UsageEvents.BulkCreate(context.Background(), &flowglad.BulkCreateUsageEventsParams{
		Events: []flowglad.CreateUsageEventParams{
			{UsageMeterID: "um_1", CustomerID: "cus_1", Quantity: 10.0},
			{UsageMeterID: "um_1", CustomerID: "cus_1", Quantity: 20.0},
		},
	})
	if err != nil {
		t.Fatalf("BulkCreate: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestUsageEvents_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /usage-events/{id}", flowgladtest.RespondWith(http.StatusOK, usageEventFixture("ue_1", "um_1", "cus_1", 10.0)))

	client := mustNewClient(t, srv.URL())
	ev, err := client.UsageEvents.Get(context.Background(), "ue_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if ev.ID != "ue_1" {
		t.Errorf("ID: got %q, want ue_1", ev.ID)
	}
}

func TestUsageEvents_List(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /usage-events", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "ue_1", "createdAt": 0, "livemode": false, "organizationId": "org", "usageMeterId": "um_1", "customerId": "cus_1", "quantity": 10.0},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for ev, err := range client.UsageEvents.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, ev.ID)
	}
	if len(ids) != 1 || ids[0] != "ue_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}
