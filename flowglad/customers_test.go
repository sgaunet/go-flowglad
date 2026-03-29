package flowglad_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/sgaunet/go-flowglad/flowglad"
	"github.com/sgaunet/go-flowglad/flowgladtest"
)

func mustNewClient(t *testing.T, baseURL string) *flowglad.Client {
	t.Helper()
	c, err := flowglad.NewClient("sk_test_fake", flowglad.WithBaseURL(baseURL), flowglad.WithNoRetry())
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func TestCustomers_Create(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /customers", flowgladtest.RespondWith(http.StatusOK, flowgladtest.CustomerFixture(
		"cus_test_123", "Test Co", "test@example.com", "ext_1",
	)))

	client := mustNewClient(t, srv.URL())
	cust, err := client.Customers.Create(context.Background(), &flowglad.CreateCustomerParams{
		Name:       "Test Co",
		Email:      "test@example.com",
		ExternalID: "ext_1",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if cust.ID != "cus_test_123" {
		t.Errorf("ID: got %q, want %q", cust.ID, "cus_test_123")
	}
	if cust.Name != "Test Co" {
		t.Errorf("Name: got %q, want %q", cust.Name, "Test Co")
	}
	if cust.Email != "test@example.com" {
		t.Errorf("Email: got %q, want %q", cust.Email, "test@example.com")
	}
}

func TestCustomers_Get(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /customers/{id}", flowgladtest.RespondWith(http.StatusOK, flowgladtest.CustomerFixture(
		"cus_abc", "Alice", "alice@example.com", "ext_alice",
	)))

	client := mustNewClient(t, srv.URL())
	cust, err := client.Customers.Get(context.Background(), "cus_abc")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if cust.ID != "cus_abc" {
		t.Errorf("ID: got %q, want %q", cust.ID, "cus_abc")
	}
}

func TestCustomers_List_SinglePage(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /customers", flowgladtest.RespondWith(http.StatusOK, flowgladtest.CustomerListFixture(
		[]map[string]any{
			{"id": "cus_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "name": "A", "email": "a@example.com", "externalId": "e1", "archived": false},
			{"id": "cus_2", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "name": "B", "email": "b@example.com", "externalId": "e2", "archived": false},
		},
		false, "",
	)))

	client := mustNewClient(t, srv.URL())
	var ids []string
	for cust, err := range client.Customers.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, cust.ID)
	}
	if len(ids) != 2 {
		t.Errorf("expected 2 customers, got %d", len(ids))
	}
	if ids[0] != "cus_1" || ids[1] != "cus_2" {
		t.Errorf("unexpected IDs: %v", ids)
	}
}

func TestCustomers_List_MultiPage(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	callCount := 0
	srv.On("GET /customers", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		cursor := r.URL.Query().Get("cursor")
		var fixture map[string]any
		if cursor == "" {
			// First page
			fixture = flowgladtest.CustomerListFixture(
				[]map[string]any{
					{"id": "cus_1", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "name": "A", "email": "a@example.com", "externalId": "e1", "archived": false},
				},
				true, "cursor_page2",
			)
		} else {
			// Second page
			fixture = flowgladtest.CustomerListFixture(
				[]map[string]any{
					{"id": "cus_2", "createdAt": 0, "updatedAt": 0, "livemode": false, "organizationId": "org", "name": "B", "email": "b@example.com", "externalId": "e2", "archived": false},
				},
				false, "",
			)
		}
		flowgladtest.RespondWith(http.StatusOK, fixture)(w, r)
	})

	client := mustNewClient(t, srv.URL())
	var ids []string
	for cust, err := range client.Customers.List(context.Background(), nil) {
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		ids = append(ids, cust.ID)
	}
	if len(ids) != 2 {
		t.Errorf("expected 2 customers across 2 pages, got %d: %v", len(ids), ids)
	}
	if callCount != 2 {
		t.Errorf("expected 2 HTTP calls for 2 pages, got %d", callCount)
	}
}

func TestCustomers_Update(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("PUT /customers/{id}", flowgladtest.RespondWith(http.StatusOK, flowgladtest.CustomerFixture(
		"cus_abc", "Updated Name", "alice@example.com", "ext_alice",
	)))

	client := mustNewClient(t, srv.URL())
	cust, err := client.Customers.Update(context.Background(), "cus_abc", &flowglad.UpdateCustomerParams{
		Name: flowglad.Ptr("Updated Name"),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if cust.Name != "Updated Name" {
		t.Errorf("Name: got %q, want %q", cust.Name, "Updated Name")
	}
}

func TestCustomers_Archive(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /customers/{id}/archive", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": map[string]any{
			"customer": map[string]any{
				"id": "cus_abc", "createdAt": 0, "updatedAt": 0, "livemode": false,
				"organizationId": "org", "name": "A", "email": "a@example.com",
				"externalId": "e1", "archived": true,
			},
		},
	}))

	client := mustNewClient(t, srv.URL())
	cust, err := client.Customers.Archive(context.Background(), "cus_abc")
	if err != nil {
		t.Fatalf("Archive: %v", err)
	}
	if !cust.Archived {
		t.Error("expected Archived to be true after archive")
	}
}

func TestCustomers_GetBillingDetails(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /customers/{id}/billing-details", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": map[string]any{
			"customerId": "cus_abc",
		},
	}))

	client := mustNewClient(t, srv.URL())
	bd, err := client.Customers.GetBillingDetails(context.Background(), "cus_abc")
	if err != nil {
		t.Fatalf("GetBillingDetails: %v", err)
	}
	if bd.CustomerID != "cus_abc" {
		t.Errorf("CustomerID: got %q, want %q", bd.CustomerID, "cus_abc")
	}
}

func TestCustomers_GetUsageBalances(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /customers/{id}/usage-balances", flowgladtest.RespondWith(http.StatusOK, map[string]any{
		"data": map[string]any{
			"customerId": "cus_abc",
			"balances": []any{
				map[string]any{
					"usageMeterId":   "um_1",
					"usageMeterName": "API Calls",
					"currentBalance": 42.0,
				},
			},
		},
	}))

	client := mustNewClient(t, srv.URL())
	ub, err := client.Customers.GetUsageBalances(context.Background(), "cus_abc")
	if err != nil {
		t.Fatalf("GetUsageBalances: %v", err)
	}
	if ub.CustomerID != "cus_abc" {
		t.Errorf("CustomerID: got %q, want %q", ub.CustomerID, "cus_abc")
	}
	if len(ub.Balances) != 1 {
		t.Fatalf("expected 1 balance, got %d", len(ub.Balances))
	}
	if ub.Balances[0].CurrentBalance != 42.0 {
		t.Errorf("CurrentBalance: got %v, want 42.0", ub.Balances[0].CurrentBalance)
	}
}

func TestCustomers_Create_Error(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /customers", flowgladtest.RespondWithError(http.StatusBadRequest, "validation", "email required"))

	client := mustNewClient(t, srv.URL())
	_, err := client.Customers.Create(context.Background(), &flowglad.CreateCustomerParams{
		Name: "Test", Email: "", ExternalID: "ext",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *flowglad.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *flowglad.Error, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode: got %d, want 400", apiErr.StatusCode)
	}
}

func TestCustomers_Update_Error(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("PUT /customers/{id}", flowgladtest.RespondWithError(http.StatusNotFound, "not_found", "not found"))

	client := mustNewClient(t, srv.URL())
	_, err := client.Customers.Update(context.Background(), "cus_gone", &flowglad.UpdateCustomerParams{
		Name: flowglad.Ptr("X"),
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCustomers_Archive_Error(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("POST /customers/{id}/archive", flowgladtest.RespondWithError(http.StatusNotFound, "not_found", "not found"))

	client := mustNewClient(t, srv.URL())
	_, err := client.Customers.Archive(context.Background(), "cus_gone")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCustomers_List_Error(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /customers", flowgladtest.RespondWithError(http.StatusForbidden, "forbidden", "access denied"))

	client := mustNewClient(t, srv.URL())
	for _, err := range client.Customers.List(context.Background(), nil) {
		if err == nil {
			t.Fatal("expected error from List iterator")
		}
		return // first iteration should yield an error
	}
}

func TestCustomers_InvalidKey_Returns401(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /customers/{id}", flowgladtest.RespondWithError(http.StatusUnauthorized, "unauthorized", "invalid API key"))

	client := mustNewClient(t, srv.URL())
	_, err := client.Customers.Get(context.Background(), "cus_any")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *flowglad.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *flowglad.Error, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode: got %d, want 401", apiErr.StatusCode)
	}
}

func TestCustomers_NotFound_Returns404(t *testing.T) {
	srv := flowgladtest.NewServer(t)
	srv.On("GET /customers/{id}", flowgladtest.RespondWithError(http.StatusNotFound, "not_found", "customer not found"))

	client := mustNewClient(t, srv.URL())
	_, err := client.Customers.Get(context.Background(), "cus_missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *flowglad.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *flowglad.Error, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode: got %d, want 404", apiErr.StatusCode)
	}
	if apiErr.Code != "not_found" {
		t.Errorf("Code: got %q, want not_found", apiErr.Code)
	}
}
