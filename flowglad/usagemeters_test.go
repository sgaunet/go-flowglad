package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func usageMeterFixture(id, name, aggregation string) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"usageMeter": map[string]any{
				"id":             id,
				"createdAt":      0,
				"updatedAt":      0,
				"livemode":       false,
				"organizationId": "org",
				"name":           name,
				"aggregation":    aggregation,
			},
		},
	}
}

func TestUsageMeters_Create(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /usage-meters", flowgladtest.RespondWith(http.StatusOK, usageMeterFixture("um_1", "API Calls", "sum")))

	client := mustNewClient(t, srv.URL())
	meter, err := client.UsageMeters.Create(context.Background(), &flowglad.CreateUsageMeterParams{
		Name:        "API Calls",
		Aggregation: flowglad.UsageMeterAggregationSum,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if meter.ID != "um_1" {
		t.Errorf("ID: got %q, want um_1", meter.ID)
	}
	if meter.Aggregation != flowglad.UsageMeterAggregationSum {
		t.Errorf("Aggregation: got %q, want sum", meter.Aggregation)
	}
}

func TestUsageMeters_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /usage-meters/{id}", flowgladtest.RespondWith(http.StatusOK, usageMeterFixture("um_1", "API Calls", "sum")))

	client := mustNewClient(t, srv.URL())
	meter, err := client.UsageMeters.Get(context.Background(), "um_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if meter.Name != "API Calls" {
		t.Errorf("Name: got %q, want API Calls", meter.Name)
	}
}
