package webhook_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"

	"github.com/sgaunet/go-flowglad/webhook"
)

const testSecret = "whsec_test_secret_abc123"

// signBody creates a valid HMAC-SHA256 signature for body using secret.
func signBody(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func makePayload(eventType webhook.EventType, data any) []byte {
	payload := map[string]any{
		"id":        "evt_test_123",
		"type":      string(eventType),
		"createdAt": int64(1700000000),
		"data":      data,
	}
	b, _ := json.Marshal(payload)
	return b
}

func TestVerify_ValidSignature(t *testing.T) {
	body := makePayload(webhook.EventTypeCustomerCreated, map[string]any{
		"customer": map[string]any{
			"id": "cus_1", "name": "Alice", "email": "alice@example.com", "externalId": "ext_1",
		},
	})
	sig := signBody(body, testSecret)

	event, err := webhook.Verify(body, sig, testSecret)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if event.Type != webhook.EventTypeCustomerCreated {
		t.Errorf("Type: got %q, want %q", event.Type, webhook.EventTypeCustomerCreated)
	}
	if event.ID != "evt_test_123" {
		t.Errorf("ID: got %q, want evt_test_123", event.ID)
	}
}

func TestVerify_ValidSignature_WithPrefix(t *testing.T) {
	body := makePayload(webhook.EventTypePaymentSucceeded, map[string]any{
		"payment": map[string]any{
			"id": "pay_1", "customerId": "cus_1", "amount": 5000, "currency": "usd", "status": "succeeded",
		},
	})
	sig := "sha256=" + signBody(body, testSecret)

	event, err := webhook.Verify(body, sig, testSecret)
	if err != nil {
		t.Fatalf("Verify with sha256= prefix: %v", err)
	}
	if event.Type != webhook.EventTypePaymentSucceeded {
		t.Errorf("Type: got %q, want %q", event.Type, webhook.EventTypePaymentSucceeded)
	}
}

func TestVerify_WrongSecret_ReturnsVerificationError(t *testing.T) {
	body := makePayload(webhook.EventTypeCustomerCreated, nil)
	sig := signBody(body, "wrong_secret")

	_, err := webhook.Verify(body, sig, testSecret)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var verErr *webhook.VerificationError
	if !errors.As(err, &verErr) {
		t.Fatalf("expected *VerificationError, got %T: %v", err, err)
	}
}

func TestVerify_EmptySignature_ReturnsVerificationError(t *testing.T) {
	body := makePayload(webhook.EventTypeCustomerCreated, nil)

	_, err := webhook.Verify(body, "", testSecret)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var verErr *webhook.VerificationError
	if !errors.As(err, &verErr) {
		t.Fatalf("expected *VerificationError, got %T", err)
	}
}

func TestVerify_EmptySecret_ReturnsVerificationError(t *testing.T) {
	body := makePayload(webhook.EventTypeCustomerCreated, nil)
	sig := signBody(body, testSecret)

	_, err := webhook.Verify(body, sig, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var verErr *webhook.VerificationError
	if !errors.As(err, &verErr) {
		t.Fatalf("expected *VerificationError, got %T", err)
	}
}

func TestVerify_TamperedBody_ReturnsVerificationError(t *testing.T) {
	body := makePayload(webhook.EventTypeCustomerCreated, nil)
	sig := signBody(body, testSecret)

	// Tamper with the body after signing.
	tamperedBody := append(body, []byte(" tampered")...)

	_, err := webhook.Verify(tamperedBody, sig, testSecret)
	if err == nil {
		t.Fatal("expected error for tampered body, got nil")
	}
	var verErr *webhook.VerificationError
	if !errors.As(err, &verErr) {
		t.Fatalf("expected *VerificationError, got %T", err)
	}
}

func TestVerify_MalformedSignature_ReturnsVerificationError(t *testing.T) {
	body := makePayload(webhook.EventTypeCustomerCreated, nil)

	_, err := webhook.Verify(body, "not-valid-hex!!", testSecret)
	if err == nil {
		t.Fatal("expected error for malformed signature, got nil")
	}
	var verErr *webhook.VerificationError
	if !errors.As(err, &verErr) {
		t.Fatalf("expected *VerificationError, got %T", err)
	}
}

func TestVerify_SubscriptionCanceled_PopulatesSubscriptionData(t *testing.T) {
	body := makePayload(webhook.EventTypeSubscriptionCanceled, map[string]any{
		"subscription": map[string]any{
			"id":         "sub_1",
			"customerId": "cus_1",
			"status":     "canceled",
		},
	})
	sig := signBody(body, testSecret)

	event, err := webhook.Verify(body, sig, testSecret)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if event.Data.Subscription == nil {
		t.Fatal("expected Subscription to be populated")
	}
	if event.Data.Subscription.ID != "sub_1" {
		t.Errorf("Subscription.ID: got %q, want sub_1", event.Data.Subscription.ID)
	}
	if event.Data.Subscription.Status != "canceled" {
		t.Errorf("Subscription.Status: got %q, want canceled", event.Data.Subscription.Status)
	}
}

func TestVerify_PaymentSucceeded_PopulatesPaymentData(t *testing.T) {
	body := makePayload(webhook.EventTypePaymentSucceeded, map[string]any{
		"payment": map[string]any{
			"id":         "pay_1",
			"customerId": "cus_1",
			"amount":     int64(5000),
			"currency":   "usd",
			"status":     "succeeded",
		},
	})
	sig := signBody(body, testSecret)

	event, err := webhook.Verify(body, sig, testSecret)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if event.Data.Payment == nil {
		t.Fatal("expected Payment to be populated")
	}
	if event.Data.Payment.Amount != 5000 {
		t.Errorf("Payment.Amount: got %d, want 5000", event.Data.Payment.Amount)
	}
}

func TestVerify_CustomerCreated_PopulatesCustomerData(t *testing.T) {
	body := makePayload(webhook.EventTypeCustomerCreated, map[string]any{
		"customer": map[string]any{
			"id":         "cus_1",
			"name":       "Alice",
			"email":      "alice@example.com",
			"externalId": "ext_1",
		},
	})
	sig := signBody(body, testSecret)

	event, err := webhook.Verify(body, sig, testSecret)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if event.Data.Customer == nil {
		t.Fatal("expected Customer to be populated")
	}
	if event.Data.Customer.Email != "alice@example.com" {
		t.Errorf("Customer.Email: got %q, want alice@example.com", event.Data.Customer.Email)
	}
}

func TestVerify_InvalidJSON_ReturnsVerificationError(t *testing.T) {
	body := []byte("not json at all")
	sig := signBody(body, testSecret)

	_, err := webhook.Verify(body, sig, testSecret)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
	var verErr *webhook.VerificationError
	if !errors.As(err, &verErr) {
		t.Fatalf("expected *VerificationError, got %T", err)
	}
}

func TestVerificationError_ErrorString(t *testing.T) {
	e := &webhook.VerificationError{Reason: "test reason"}
	if e.Error() == "" {
		t.Error("expected non-empty error string")
	}
}
