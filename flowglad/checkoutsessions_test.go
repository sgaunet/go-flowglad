package flowglad_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func TestCheckoutSessions_Create(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /checkout-sessions", flowgladtest.RespondWith(http.StatusOK, flowgladtest.CheckoutSessionFixture(
		"cs_1", "https://checkout.flowglad.com/cs_1", "open",
	)))

	client := mustNewClient(t, srv.URL())
	cs, err := client.CheckoutSessions.Create(context.Background(), &flowglad.CreateCheckoutSessionParams{
		PriceID: "price_1",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if cs.ID != "cs_1" {
		t.Errorf("ID: got %q, want cs_1", cs.ID)
	}
	if cs.URL == "" {
		t.Error("URL should not be empty")
	}
	if cs.Status != "open" {
		t.Errorf("Status: got %q, want open", cs.Status)
	}
}

func TestCheckoutSessions_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /checkout-sessions/{id}", flowgladtest.RespondWith(http.StatusOK, flowgladtest.CheckoutSessionFixture(
		"cs_1", "https://checkout.flowglad.com/cs_1", "complete",
	)))

	client := mustNewClient(t, srv.URL())
	cs, err := client.CheckoutSessions.Get(context.Background(), "cs_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if cs.Status != "complete" {
		t.Errorf("Status: got %q, want complete", cs.Status)
	}
}

func TestCheckoutSessions_Get_Expired(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /checkout-sessions/{id}", flowgladtest.RespondWith(http.StatusOK, flowgladtest.CheckoutSessionFixture(
		"cs_expired", "https://checkout.flowglad.com/cs_expired", "expired",
	)))

	client := mustNewClient(t, srv.URL())
	cs, err := client.CheckoutSessions.Get(context.Background(), "cs_expired")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if cs.Status != "expired" {
		t.Errorf("Status: got %q, want expired", cs.Status)
	}
}

func TestCheckoutSessions_List(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /checkout-sessions", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": []any{
			map[string]any{"id": "cs_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "url": "https://x.com/1", "status": "open"},
		},
		"hasMore": false, "nextCursor": "", "currentCursor": "", "total": 1,
	}))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for cs, err := range client.CheckoutSessions.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, cs.ID)
	}
	if len(ids) != 1 || ids[0] != "cs_1" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}
