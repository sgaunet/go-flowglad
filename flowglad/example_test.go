package flowglad_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/sgaunet/go-flowglad/flowglad"
)

// ExampleNewClient demonstrates constructing a Flowglad SDK client.
func ExampleNewClient() {
	client, err := flowglad.NewClient("sk_test_YOUR_API_KEY")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	_ = client
	fmt.Println("client ready")
	// Output: client ready
}

// ExamplePtr demonstrates using the Ptr helper for optional fields.
func ExamplePtr() {
	name := flowglad.Ptr("Alice Corp")
	fmt.Println(*name)
	// Output: Alice Corp
}

// ExampleCustomersClient_Create demonstrates creating a customer.
func ExampleCustomersClient_Create() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":{"customer":{"id":"cus_example_1","createdAt":0,"updatedAt":0,"livemode":false,"organizationId":"org","name":"Alice Corp","email":"alice@example.com","externalId":"","archived":false}}}`)
	}))
	defer srv.Close()

	client, _ := flowglad.NewClient("sk_test_fake", flowglad.WithBaseURL(srv.URL))
	cust, err := client.Customers.Create(context.Background(), &flowglad.CreateCustomerParams{
		Name:  "Alice Corp",
		Email: "alice@example.com",
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("created:", cust.ID)
	// Output: created: cus_example_1
}

// ExampleCustomersClient_List demonstrates iterating customers with range-over-func.
func ExampleCustomersClient_List() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[{"id":"cus_a","createdAt":0,"updatedAt":0,"livemode":false,"organizationId":"org","name":"A","email":"a@example.com","externalId":"","archived":false},{"id":"cus_b","createdAt":0,"updatedAt":0,"livemode":false,"organizationId":"org","name":"B","email":"b@example.com","externalId":"","archived":false}],"hasMore":false,"nextCursor":"","currentCursor":"","total":2}`)
	}))
	defer srv.Close()

	client, _ := flowglad.NewClient("sk_test_fake", flowglad.WithBaseURL(srv.URL))
	for cust, err := range client.Customers.List(context.Background(), nil) {
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println(cust.ID)
	}
	// Output:
	// cus_a
	// cus_b
}

// ExampleError demonstrates typed error handling with errors.As.
func ExampleError() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"customer not found","code":"not_found","issues":[]}`)
	}))
	defer srv.Close()

	client, _ := flowglad.NewClient("sk_test_fake", flowglad.WithBaseURL(srv.URL), flowglad.WithNoRetry())
	_, err := client.Customers.Get(context.Background(), "nonexistent")
	if err != nil {
		var apiErr *flowglad.Error
		if errors.As(err, &apiErr) {
			fmt.Printf("HTTP %d | code: %s\n", apiErr.StatusCode, apiErr.Code)
		}
	}
	// Output: HTTP 404 | code: not_found
}
