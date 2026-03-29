// Package webhook provides Flowglad webhook signature verification and event parsing.
//
// Flowglad signs webhook payloads using HMAC-SHA256. The signature is delivered
// in the X-Flowglad-Signature HTTP header, encoded as a hex string, optionally
// prefixed with "sha256=".
//
// Usage:
//
//	event, err := webhook.Verify(body, r.Header.Get("X-Flowglad-Signature"), secret)
//	if err != nil {
//	    // signature invalid or payload unreadable
//	}
//	switch event.Type {
//	case webhook.EventTypeSubscriptionCanceled:
//	    // handle cancellation
//	}
package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
)

// EventType identifies the Flowglad webhook event type.
type EventType string

const (
	// EventTypeCustomerCreated fires when a customer is created.
	EventTypeCustomerCreated EventType = "customer.created"
	// EventTypeCustomerUpdated fires when a customer is updated.
	EventTypeCustomerUpdated EventType = "customer.updated"
	// EventTypePurchaseCompleted fires when a purchase is completed.
	EventTypePurchaseCompleted EventType = "purchase.completed"
	// EventTypePaymentFailed fires when a payment fails.
	EventTypePaymentFailed EventType = "payment.failed"
	// EventTypePaymentSucceeded fires when a payment succeeds.
	EventTypePaymentSucceeded EventType = "payment.succeeded"
	// EventTypeSubscriptionCreated fires when a subscription is created.
	EventTypeSubscriptionCreated EventType = "subscription.created"
	// EventTypeSubscriptionUpdated fires when a subscription is updated.
	EventTypeSubscriptionUpdated EventType = "subscription.updated"
	// EventTypeSubscriptionCanceled fires when a subscription is canceled.
	EventTypeSubscriptionCanceled EventType = "subscription.canceled"
	// EventTypeSyncEventsAvailable fires when sync events are available.
	EventTypeSyncEventsAvailable EventType = "sync.events_available"
)

// VerificationError is returned when webhook signature validation fails.
type VerificationError struct {
	// Reason explains why verification failed.
	Reason string
}

// Error implements the error interface.
func (e *VerificationError) Error() string {
	return "webhook: verification failed: " + e.Reason
}

// EventData holds the typed payload of a webhook event.
// Only the field matching the event type will be populated.
type EventData struct {
	// Customer is populated for customer.* events.
	Customer *EventCustomer `json:"customer,omitempty"`
	// Subscription is populated for subscription.* events.
	Subscription *EventSubscription `json:"subscription,omitempty"`
	// Payment is populated for payment.* events.
	Payment *EventPayment `json:"payment,omitempty"`
}

// EventCustomer is the customer payload embedded in customer.* webhook events.
type EventCustomer struct {
	// ID is the Flowglad customer identifier.
	ID string `json:"id"`
	// Name is the customer's display name.
	Name string `json:"name"`
	// Email is the customer's email address.
	Email string `json:"email"`
	// ExternalID is the caller-supplied external identifier.
	ExternalID string `json:"externalId"`
}

// EventSubscription is the subscription payload embedded in subscription.* webhook events.
type EventSubscription struct {
	// ID is the Flowglad subscription identifier.
	ID string `json:"id"`
	// CustomerID is the subscribed customer.
	CustomerID string `json:"customerId"`
	// Status is the subscription lifecycle status.
	Status string `json:"status"`
}

// EventPayment is the payment payload embedded in payment.* webhook events.
type EventPayment struct {
	// ID is the Flowglad payment identifier.
	ID string `json:"id"`
	// CustomerID is the paying customer.
	CustomerID string `json:"customerId"`
	// Amount is the payment amount in the smallest currency unit.
	Amount int64 `json:"amount"`
	// Currency is the ISO 4217 currency code.
	Currency string `json:"currency"`
	// Status is the payment status.
	Status string `json:"status"`
}

// Event is a parsed, strongly-typed Flowglad webhook event.
type Event struct {
	// ID is the unique event identifier.
	ID string `json:"id"`
	// Type is the event type.
	Type EventType `json:"type"`
	// CreatedAt is the Unix timestamp when the event was generated.
	CreatedAt int64 `json:"createdAt"`
	// Data holds the typed event payload.
	Data EventData `json:"data"`
}

// rawEvent is used for initial JSON parsing before signature verification.
type rawEvent struct {
	ID        string          `json:"id"`
	Type      EventType       `json:"type"`
	CreatedAt int64           `json:"createdAt"`
	Data      json.RawMessage `json:"data"`
}

// Verify validates the HMAC-SHA256 signature of a webhook request body and
// returns a typed Event if the signature is valid.
//
// body is the raw request body bytes.
// signature is the value of the X-Flowglad-Signature header. It may be a bare
// hex string or prefixed with "sha256=".
// secret is the webhook signing secret returned when the webhook was created.
//
// Returns a *VerificationError if the signature is invalid or missing.
func Verify(body []byte, signature, secret string) (*Event, error) {
	if signature == "" {
		return nil, &VerificationError{Reason: "missing signature header"}
	}
	if secret == "" {
		return nil, &VerificationError{Reason: "missing webhook secret"}
	}

	// Strip optional "sha256=" prefix.
	sig := signature
	if after, found := strings.CutPrefix(signature, "sha256="); found {
		sig = after
	}

	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		return nil, &VerificationError{Reason: "malformed signature: " + err.Error()}
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := mac.Sum(nil)

	if !hmac.Equal(sigBytes, expected) {
		return nil, &VerificationError{Reason: "signature mismatch"}
	}

	var raw rawEvent
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, &VerificationError{Reason: "invalid JSON payload: " + err.Error()}
	}

	var data EventData
	if raw.Data != nil {
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return nil, &VerificationError{Reason: "invalid event data: " + err.Error()}
		}
	}

	return &Event{
		ID:        raw.ID,
		Type:      raw.Type,
		CreatedAt: raw.CreatedAt,
		Data:      data,
	}, nil
}
