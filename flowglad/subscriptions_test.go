package flowglad_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func subFixture(id, customerID, priceID, status string) map[string]any {
	return flowgladtest.SubscriptionFixture(id, customerID, priceID, status)
}

func TestSubscriptions_Create(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /subscriptions", flowgladtest.RespondWith(http.StatusOK, subFixture(
		"sub_1", "cus_1", "price_1", "active",
	)))

	client := mustNewClient(t, srv.URL())
	sub, err := client.Subscriptions.Create(context.Background(), &flowglad.CreateSubscriptionParams{
		CustomerID: flowglad.Ptr("cus_1"),
		PriceID:    flowglad.Ptr("price_1"),
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if sub.ID != "sub_1" {
		t.Errorf("ID: got %q, want sub_1", sub.ID)
	}
	if sub.Status != flowglad.SubscriptionStatusActive {
		t.Errorf("Status: got %q, want active", sub.Status)
	}
}

func TestSubscriptions_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /subscriptions/{id}", flowgladtest.RespondWith(http.StatusOK, subFixture(
		"sub_1", "cus_1", "price_1", "active",
	)))

	client := mustNewClient(t, srv.URL())
	sub, err := client.Subscriptions.Get(context.Background(), "sub_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if sub.ID != "sub_1" {
		t.Errorf("ID: got %q, want sub_1", sub.ID)
	}
}

func TestSubscriptions_List(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /subscriptions", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "sub_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "customerId": "cus_1", "priceId": "price_1", "status": "active"},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for sub, err := range client.Subscriptions.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, sub.ID)
	}
	if len(ids) != 1 || ids[0] != "sub_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}

func TestSubscriptions_PreviewAdjust(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /subscriptions/{id}/preview-adjust", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": map[string]any{
			"adjustmentPreview": map[string]any{
				"amountDue": 5000,
				"currency":  "usd",
			},
		},
	}))

	client := mustNewClient(t, srv.URL())
	preview, err := client.Subscriptions.PreviewAdjust(context.Background(), "sub_1", &flowglad.PreviewAdjustmentParams{
		Adjustment: flowglad.AdjustmentRequest{
			NewPriceID: flowglad.Ptr("price_new"),
		},
	})
	if err != nil {
		t.Fatalf("PreviewAdjust: %v", err)
	}
	if preview.AmountDue != 5000 {
		t.Errorf("AmountDue: got %d, want 5000", preview.AmountDue)
	}
}

func TestSubscriptions_Adjust(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /subscriptions/{id}/adjust", flowgladtest.RespondWith(http.StatusOK, subFixture(
		"sub_1", "cus_1", "price_new", "active",
	)))

	client := mustNewClient(t, srv.URL())
	sub, err := client.Subscriptions.Adjust(context.Background(), "sub_1", &flowglad.AdjustmentParams{
		Adjustment: flowglad.AdjustmentRequest{
			NewPriceID: flowglad.Ptr("price_new"),
		},
	})
	if err != nil {
		t.Fatalf("Adjust: %v", err)
	}
	if sub.PriceID != "price_new" {
		t.Errorf("PriceID: got %q, want price_new", sub.PriceID)
	}
}

func TestSubscriptions_Cancel_Immediately(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /subscriptions/{id}/cancel", flowgladtest.RespondWith(http.StatusOK, subFixture(
		"sub_1", "cus_1", "price_1", "canceled",
	)))

	client := mustNewClient(t, srv.URL())
	sub, err := client.Subscriptions.Cancel(context.Background(), "sub_1", &flowglad.CancelSubscriptionParams{
		CancellationTiming: flowglad.CancellationTimingImmediately,
	})
	if err != nil {
		t.Fatalf("Cancel: %v", err)
	}
	if sub.Status != flowglad.SubscriptionStatusCanceled {
		t.Errorf("Status: got %q, want canceled", sub.Status)
	}
}

func TestSubscriptions_Cancel_AtPeriodEnd(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /subscriptions/{id}/cancel", flowgladtest.RespondWith(http.StatusOK, subFixture(
		"sub_1", "cus_1", "price_1", "cancellation_scheduled",
	)))

	client := mustNewClient(t, srv.URL())
	sub, err := client.Subscriptions.Cancel(context.Background(), "sub_1", &flowglad.CancelSubscriptionParams{
		CancellationTiming: flowglad.CancellationTimingAtEndOfCurrentBillingPeriod,
	})
	if err != nil {
		t.Fatalf("Cancel: %v", err)
	}
	if sub.Status != flowglad.SubscriptionStatusCancellationScheduled {
		t.Errorf("Status: got %q, want cancellation_scheduled", sub.Status)
	}
}

func TestSubscriptions_Uncancel(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /subscriptions/{id}/uncancel", flowgladtest.RespondWith(http.StatusOK, subFixture(
		"sub_1", "cus_1", "price_1", "active",
	)))

	client := mustNewClient(t, srv.URL())
	sub, err := client.Subscriptions.Uncancel(context.Background(), "sub_1")
	if err != nil {
		t.Fatalf("Uncancel: %v", err)
	}
	if sub.Status != flowglad.SubscriptionStatusActive {
		t.Errorf("Status: got %q, want active", sub.Status)
	}
}

func TestSubscriptions_AddFeature(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /subscriptions/{id}/add-feature", flowgladtest.RespondWith(http.StatusOK, subFixture(
		"sub_1", "cus_1", "price_1", "active",
	)))

	client := mustNewClient(t, srv.URL())
	sub, err := client.Subscriptions.AddFeature(context.Background(), "sub_1", &flowglad.AddFeatureParams{
		FeatureID: "feat_1",
	})
	if err != nil {
		t.Fatalf("AddFeature: %v", err)
	}
	if sub.ID != "sub_1" {
		t.Errorf("ID: got %q, want sub_1", sub.ID)
	}
}

func TestSubscriptions_CancelScheduledAdjustment(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /subscriptions/{id}/cancel-scheduled-adjustment", flowgladtest.RespondWith(http.StatusOK, subFixture(
		"sub_1", "cus_1", "price_1", "active",
	)))

	client := mustNewClient(t, srv.URL())
	sub, err := client.Subscriptions.CancelScheduledAdjustment(context.Background(), "sub_1")
	if err != nil {
		t.Fatalf("CancelScheduledAdjustment: %v", err)
	}
	if sub.ID != "sub_1" {
		t.Errorf("ID: got %q, want sub_1", sub.ID)
	}
}

func TestSubscriptions_NotFound(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /subscriptions/{id}", flowgladtest.RespondWithError(http.StatusNotFound, "not_found", "subscription not found"))

	client := mustNewClient(t, srv.URL())
	_, err := client.Subscriptions.Get(context.Background(), "sub_missing")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *flowglad.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *flowglad.Error, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("StatusCode: got %d, want 404", apiErr.StatusCode)
	}
}

func TestSubscriptions_Create_Error(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /subscriptions", flowgladtest.RespondWithError(http.StatusBadRequest, "validation", "bad"))

	client := mustNewClient(t, srv.URL())
	_, err := client.Subscriptions.Create(context.Background(), &flowglad.CreateSubscriptionParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSubscriptions_Cancel_Error(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /subscriptions/{id}/cancel", flowgladtest.RespondWithError(http.StatusBadRequest, "invalid", "bad"))

	client := mustNewClient(t, srv.URL())
	_, err := client.Subscriptions.Cancel(context.Background(), "sub_1", &flowglad.CancelSubscriptionParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSubscriptions_Adjust_Error(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /subscriptions/{id}/adjust", flowgladtest.RespondWithError(http.StatusBadRequest, "invalid", "bad"))

	client := mustNewClient(t, srv.URL())
	_, err := client.Subscriptions.Adjust(context.Background(), "sub_1", &flowglad.AdjustmentParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSubscriptions_Uncancel_Error(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /subscriptions/{id}/uncancel", flowgladtest.RespondWithError(http.StatusNotFound, "not_found", "not found"))

	client := mustNewClient(t, srv.URL())
	_, err := client.Subscriptions.Uncancel(context.Background(), "sub_gone")
	if err == nil {
		t.Fatal("expected error")
	}
}
