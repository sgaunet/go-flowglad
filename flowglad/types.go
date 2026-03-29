package flowglad

import "iter"

// ---------------------------------------------------------------------------
// Pagination helpers
// ---------------------------------------------------------------------------

// ListParams are common cursor-based pagination parameters accepted by all
// List endpoints.
type ListParams struct {
	// Cursor is the opaque pagination cursor returned by the previous page.
	Cursor *string
	// Limit is the maximum number of items to return per page.
	Limit *int
}

// page is the internal JSON envelope for paginated list responses.
type page[T any] struct {
	Data          []T    `json:"data"`
	CurrentCursor string `json:"currentCursor"`
	NextCursor    string `json:"nextCursor"`
	HasMore       bool   `json:"hasMore"`
	Total         int    `json:"total"`
}

// listIter returns an iter.Seq2 that transparently fetches pages using do.
// do receives the current cursor (empty string on the first call) and returns
// the next page or an error.
func listIter[T any](do func(cursor string) (*page[T], error)) iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		cursor := ""
		for {
			p, err := do(cursor)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range p.Data {
				if !yield(&p.Data[i], nil) {
					return
				}
			}
			if !p.HasMore {
				return
			}
			cursor = p.NextCursor
		}
	}
}

// Ptr returns a pointer to v. Convenience helper for optional/nullable fields.
func Ptr[T any](v T) *T { return &v }

// ---------------------------------------------------------------------------
// Customers
// ---------------------------------------------------------------------------

// BillingAddress holds a postal address used for invoicing.
type BillingAddress struct {
	// Line1 is the first address line.
	Line1 string `json:"line1"`
	// Line2 is the optional second address line.
	Line2 *string `json:"line2,omitempty"`
	// City is the city or locality.
	City string `json:"city"`
	// State is the state, province, or region code.
	State *string `json:"state,omitempty"`
	// PostalCode is the ZIP or postal code.
	PostalCode string `json:"postalCode"`
	// Country is the ISO 3166-1 alpha-2 country code.
	Country string `json:"country"`
}

// Customer represents a Flowglad billing customer.
type Customer struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates whether the customer is in live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// Name is the customer's display name.
	Name string `json:"name"`
	// Email is the customer's email address.
	Email string `json:"email"`
	// ExternalID is the caller-supplied external identifier.
	ExternalID string `json:"externalId"`
	// Archived indicates whether this customer is archived.
	Archived bool `json:"archived"`
	// BillingAddress is the optional postal billing address.
	BillingAddress *BillingAddress `json:"billingAddress,omitempty"`
}

// CreateCustomerParams are the parameters for creating a new customer.
type CreateCustomerParams struct {
	// Name is the customer's display name.
	Name string `json:"name"`
	// Email is the customer's email address.
	Email string `json:"email"`
	// ExternalID is an optional caller-supplied external identifier.
	ExternalID string `json:"externalId,omitempty"`
	// BillingAddress is the optional postal billing address.
	BillingAddress *BillingAddress `json:"billingAddress,omitempty"`
}

// UpdateCustomerParams are the parameters for updating an existing customer.
// All fields are optional; omit a field to leave it unchanged.
type UpdateCustomerParams struct {
	// Name is the new display name.
	Name *string `json:"name,omitempty"`
	// Email is the new email address.
	Email *string `json:"email,omitempty"`
	// ExternalID is the new external identifier.
	ExternalID *string `json:"externalId,omitempty"`
	// BillingAddress is the new postal billing address.
	BillingAddress *BillingAddress `json:"billingAddress,omitempty"`
}

// ListCustomersParams are the query parameters for listing customers.
type ListCustomersParams struct {
	ListParams
}

// CustomerBillingDetails holds the billing details associated with a customer.
type CustomerBillingDetails struct {
	// CustomerID is the Flowglad customer identifier.
	CustomerID string `json:"customerId"`
	// BillingAddress is the customer's postal billing address.
	BillingAddress *BillingAddress `json:"billingAddress,omitempty"`
	// DefaultPaymentMethod is the customer's default payment method, if set.
	DefaultPaymentMethod *PaymentMethod `json:"defaultPaymentMethod,omitempty"`
}

// UsageBalance holds the current usage meter balance for a customer.
type UsageBalance struct {
	// UsageMeterID is the identifier of the usage meter.
	UsageMeterID string `json:"usageMeterId"`
	// UsageMeterName is the human-readable meter name.
	UsageMeterName string `json:"usageMeterName"`
	// CurrentBalance is the current consumption balance.
	CurrentBalance float64 `json:"currentBalance"`
}

// CustomerUsageBalances holds all usage meter balances for a customer.
type CustomerUsageBalances struct {
	// CustomerID is the Flowglad customer identifier.
	CustomerID string `json:"customerId"`
	// Balances lists the current balance per usage meter.
	Balances []UsageBalance `json:"balances"`
}

// ---------------------------------------------------------------------------
// Subscriptions
// ---------------------------------------------------------------------------

// SubscriptionStatus represents the lifecycle state of a subscription.
type SubscriptionStatus string

const (
	// SubscriptionStatusTrialing indicates the subscription is in a free trial.
	SubscriptionStatusTrialing SubscriptionStatus = "trialing"
	// SubscriptionStatusActive indicates the subscription is active and billing normally.
	SubscriptionStatusActive SubscriptionStatus = "active"
	// SubscriptionStatusPastDue indicates a payment failure; the subscription is still active.
	SubscriptionStatusPastDue SubscriptionStatus = "past_due"
	// SubscriptionStatusUnpaid indicates persistent payment failure; access may be revoked.
	SubscriptionStatusUnpaid SubscriptionStatus = "unpaid"
	// SubscriptionStatusCancellationScheduled indicates cancellation is queued for period end.
	SubscriptionStatusCancellationScheduled SubscriptionStatus = "cancellation_scheduled"
	// SubscriptionStatusIncomplete indicates the subscription awaits initial payment confirmation.
	SubscriptionStatusIncomplete SubscriptionStatus = "incomplete"
	// SubscriptionStatusIncompleteExpired indicates the initial payment window has passed.
	SubscriptionStatusIncompleteExpired SubscriptionStatus = "incomplete_expired"
	// SubscriptionStatusCanceled indicates the subscription has been cancelled.
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
	// SubscriptionStatusPaused indicates the subscription is temporarily paused.
	SubscriptionStatusPaused SubscriptionStatus = "paused"
	// SubscriptionStatusCreditTrial indicates the subscription is in a credit-funded trial.
	SubscriptionStatusCreditTrial SubscriptionStatus = "credit_trial"
)

// CancellationTiming controls when a subscription cancellation takes effect.
type CancellationTiming string

const (
	// CancellationTimingImmediately cancels the subscription immediately.
	CancellationTimingImmediately CancellationTiming = "immediately"
	// CancellationTimingAtEndOfCurrentBillingPeriod cancels at the end of the current period.
	CancellationTimingAtEndOfCurrentBillingPeriod CancellationTiming = "at_end_of_current_billing_period"
)

// AdjustmentTiming controls when a subscription adjustment takes effect.
type AdjustmentTiming string

const (
	// AdjustmentTimingImmediately applies the adjustment immediately.
	AdjustmentTimingImmediately AdjustmentTiming = "immediately"
	// AdjustmentTimingAtEndOfCurrentBillingPeriod applies the adjustment at period end.
	AdjustmentTimingAtEndOfCurrentBillingPeriod AdjustmentTiming = "at_end_of_current_billing_period"
	// AdjustmentTimingAuto lets Flowglad choose the timing.
	AdjustmentTimingAuto AdjustmentTiming = "auto"
)

// CancellationDetails holds information about a scheduled or completed cancellation.
type CancellationDetails struct {
	// Reason is the cancellation reason provided by the caller.
	Reason *string `json:"reason,omitempty"`
	// EffectiveAt is the Unix timestamp when the cancellation will or did take effect.
	EffectiveAt *int64 `json:"effectiveAt,omitempty"`
}

// Subscription represents a recurring billing agreement between a customer and a price.
type Subscription struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates whether the subscription is in live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// CustomerID is the associated customer.
	CustomerID string `json:"customerId"`
	// PriceID is the price this subscription is on.
	PriceID string `json:"priceId"`
	// Status is the current subscription lifecycle status.
	Status SubscriptionStatus `json:"status"`
	// CurrentPeriodStart is the Unix timestamp of the current billing period start.
	CurrentPeriodStart *int64 `json:"currentPeriodStart,omitempty"`
	// CurrentPeriodEnd is the Unix timestamp of the current billing period end.
	CurrentPeriodEnd *int64 `json:"currentPeriodEnd,omitempty"`
	// CancellationDetails holds cancellation information if applicable.
	CancellationDetails *CancellationDetails `json:"cancellationDetails,omitempty"`
	// TrialEnd is the Unix timestamp when the trial ends, if applicable.
	TrialEnd *int64 `json:"trialEnd,omitempty"`
}

// SubscriptionItemInput specifies a subscription item when creating or adjusting.
type SubscriptionItemInput struct {
	// PriceID is the price for this subscription item.
	PriceID string `json:"priceId"`
	// Quantity is the number of units, for per-seat pricing.
	Quantity *int `json:"quantity,omitempty"`
}

// CreateSubscriptionParams are the parameters for creating a new subscription.
type CreateSubscriptionParams struct {
	// CustomerID is the customer to subscribe.
	CustomerID *string `json:"customerId,omitempty"`
	// PriceID is the price to subscribe to.
	PriceID *string `json:"priceId,omitempty"`
	// Items lists the subscription items (alternative to PriceID for multi-item subscriptions).
	Items []SubscriptionItemInput `json:"items,omitempty"`
	// TrialEnd is the optional Unix timestamp when a trial period ends.
	TrialEnd *int64 `json:"trialEnd,omitempty"`
}

// CancelSubscriptionParams are the parameters for cancelling a subscription.
type CancelSubscriptionParams struct {
	// CancellationTiming controls when the cancellation takes effect.
	CancellationTiming CancellationTiming `json:"cancellationTiming"`
	// Reason is an optional human-readable reason for the cancellation.
	Reason *string `json:"reason,omitempty"`
}

// AdjustmentRequest specifies a single adjustment operation on a subscription.
type AdjustmentRequest struct {
	// NewPriceID replaces the current price with this one, if set.
	NewPriceID *string `json:"newPriceId,omitempty"`
	// Quantity sets the new quantity for the subscription item.
	Quantity *int `json:"quantity,omitempty"`
	// Timing controls when the adjustment takes effect.
	Timing AdjustmentTiming `json:"timing,omitempty"`
}

// AdjustmentParams is the top-level request body for the adjust endpoint.
type AdjustmentParams struct {
	// Adjustment contains the adjustment details.
	Adjustment AdjustmentRequest `json:"adjustment"`
}

// PreviewAdjustmentParams are the parameters for previewing a subscription adjustment.
type PreviewAdjustmentParams struct {
	// Adjustment contains the proposed adjustment to preview.
	Adjustment AdjustmentRequest `json:"adjustment"`
}

// AdjustmentPreview holds the results of a proration preview calculation.
type AdjustmentPreview struct {
	// AmountDue is the amount in cents that will be charged immediately.
	AmountDue int64 `json:"amountDue"`
	// Currency is the ISO 4217 currency code.
	Currency string `json:"currency"`
	// EffectiveAt is the Unix timestamp when the adjustment will take effect.
	EffectiveAt *int64 `json:"effectiveAt,omitempty"`
	// ProrationDate is the Unix timestamp used for proration calculation.
	ProrationDate *int64 `json:"prorationDate,omitempty"`
}

// AddFeatureParams are the parameters for adding a feature to a subscription.
type AddFeatureParams struct {
	// FeatureID is the identifier of the feature to add.
	FeatureID string `json:"featureId"`
}

// ListSubscriptionsParams are the query parameters for listing subscriptions.
type ListSubscriptionsParams struct {
	ListParams

	// CustomerID filters subscriptions by customer.
	CustomerID *string
	// Status filters subscriptions by status.
	Status *SubscriptionStatus
}

// ---------------------------------------------------------------------------
// Products
// ---------------------------------------------------------------------------

// Product represents a billable offering in the Flowglad product catalog.
type Product struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates whether this product is in live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// Name is the product's display name.
	Name string `json:"name"`
	// Description is an optional human-readable description.
	Description *string `json:"description,omitempty"`
	// Active indicates whether the product is available for purchase.
	Active bool `json:"active"`
}

// CreateProductParams are the parameters for creating a new product.
type CreateProductParams struct {
	// Name is the product display name.
	Name string `json:"name"`
	// Description is an optional description.
	Description *string `json:"description,omitempty"`
	// Active sets the initial active state (defaults to true if omitted).
	Active *bool `json:"active,omitempty"`
}

// UpdateProductParams are the parameters for updating an existing product.
type UpdateProductParams struct {
	// Name is the new display name.
	Name *string `json:"name,omitempty"`
	// Description is the new description.
	Description *string `json:"description,omitempty"`
	// Active updates the active flag.
	Active *bool `json:"active,omitempty"`
}

// ListProductsParams are the query parameters for listing products.
type ListProductsParams struct {
	ListParams
}

// ---------------------------------------------------------------------------
// Prices
// ---------------------------------------------------------------------------

// Price represents a specific pricing configuration for a product.
type Price struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// ProductID is the associated product.
	ProductID string `json:"productId"`
	// Name is the price's display name.
	Name string `json:"name"`
	// Currency is the ISO 4217 currency code (e.g. "usd").
	Currency string `json:"currency"`
	// UnitAmount is the price in the currency's smallest unit (e.g. cents).
	UnitAmount int64 `json:"unitAmount"`
	// Interval is the billing interval ("month", "year", etc.) for recurring prices.
	Interval *string `json:"interval,omitempty"`
	// IntervalCount is the number of intervals between billings.
	IntervalCount *int `json:"intervalCount,omitempty"`
	// Active indicates whether this price is available for use.
	Active bool `json:"active"`
	// Slug is an optional unique human-readable identifier.
	Slug *string `json:"slug,omitempty"`
}

// CreatePriceParams are the parameters for creating a new price.
type CreatePriceParams struct {
	// Name is the price display name.
	Name string `json:"name"`
	// ProductID is the product this price belongs to.
	ProductID string `json:"productId"`
	// Currency is the ISO 4217 currency code.
	Currency string `json:"currency"`
	// UnitAmount is the price in the smallest currency unit.
	UnitAmount int64 `json:"unitAmount"`
	// Interval is the billing interval for recurring prices.
	Interval *string `json:"interval,omitempty"`
	// IntervalCount is the number of intervals between billings.
	IntervalCount *int `json:"intervalCount,omitempty"`
	// Active sets the initial active state.
	Active *bool `json:"active,omitempty"`
	// Slug is an optional human-readable identifier.
	Slug *string `json:"slug,omitempty"`
}

// UpdatePriceParams are the parameters for updating an existing price.
type UpdatePriceParams struct {
	// Name is the new display name.
	Name *string `json:"name,omitempty"`
	// Active updates the active flag.
	Active *bool `json:"active,omitempty"`
}

// ListPricesParams are the query parameters for listing prices.
type ListPricesParams struct {
	ListParams

	// ProductID filters prices by product.
	ProductID *string
}

// ---------------------------------------------------------------------------
// Checkout Sessions
// ---------------------------------------------------------------------------

// CheckoutSession represents a transient session for an in-progress purchase flow.
type CheckoutSession struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// URL is the hosted checkout URL to redirect the customer to.
	URL string `json:"url"`
	// Status is the current session state (e.g. "open", "complete", "expired").
	Status string `json:"status"`
	// PriceID is the price being purchased.
	PriceID *string `json:"priceId,omitempty"`
	// CustomerID is the pre-filled customer, if any.
	CustomerID *string `json:"customerId,omitempty"`
	// ExpiresAt is the Unix timestamp when the session expires.
	ExpiresAt *int64 `json:"expiresAt,omitempty"`
	// SuccessURL is the URL to redirect to after successful payment.
	SuccessURL *string `json:"successUrl,omitempty"`
	// CancelURL is the URL to redirect to if the customer cancels.
	CancelURL *string `json:"cancelUrl,omitempty"`
}

// CreateCheckoutSessionParams are the parameters for creating a checkout session.
type CreateCheckoutSessionParams struct {
	// PriceID is the price to checkout.
	PriceID string `json:"priceId"`
	// CustomerID pre-fills the customer on the checkout page.
	CustomerID *string `json:"customerId,omitempty"`
	// SuccessURL is the redirect target after successful payment.
	SuccessURL *string `json:"successUrl,omitempty"`
	// CancelURL is the redirect target if the customer cancels.
	CancelURL *string `json:"cancelUrl,omitempty"`
}

// ListCheckoutSessionsParams are the query parameters for listing checkout sessions.
type ListCheckoutSessionsParams struct {
	ListParams
}

// ---------------------------------------------------------------------------
// Invoices
// ---------------------------------------------------------------------------

// InvoiceStatus represents the billing state of an invoice.
type InvoiceStatus string

const (
	// InvoiceStatusDraft indicates the invoice is still being prepared.
	InvoiceStatusDraft InvoiceStatus = "draft"
	// InvoiceStatusOpen indicates the invoice is finalized and awaiting payment.
	InvoiceStatusOpen InvoiceStatus = "open"
	// InvoiceStatusPaid indicates the invoice has been fully paid.
	InvoiceStatusPaid InvoiceStatus = "paid"
	// InvoiceStatusVoid indicates the invoice has been voided.
	InvoiceStatusVoid InvoiceStatus = "void"
	// InvoiceStatusUncollectible indicates the invoice is deemed uncollectible.
	InvoiceStatusUncollectible InvoiceStatus = "uncollectible"
)

// Invoice represents a billing document for a subscription period or one-time charge.
type Invoice struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// CustomerID is the billed customer.
	CustomerID string `json:"customerId"`
	// SubscriptionID is the associated subscription, if any.
	SubscriptionID *string `json:"subscriptionId,omitempty"`
	// Status is the current invoice state.
	Status InvoiceStatus `json:"status"`
	// AmountDue is the total amount due in the smallest currency unit.
	AmountDue int64 `json:"amountDue"`
	// AmountPaid is the amount already paid.
	AmountPaid int64 `json:"amountPaid"`
	// Currency is the ISO 4217 currency code.
	Currency string `json:"currency"`
	// DueDate is the Unix timestamp when payment is due.
	DueDate *int64 `json:"dueDate,omitempty"`
}

// InvoiceLineItem represents a single line on an invoice.
type InvoiceLineItem struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// InvoiceID is the parent invoice.
	InvoiceID string `json:"invoiceId"`
	// Description is the line item description.
	Description string `json:"description"`
	// Amount is the line item amount in the smallest currency unit.
	Amount int64 `json:"amount"`
	// Currency is the ISO 4217 currency code.
	Currency string `json:"currency"`
	// Quantity is the number of units.
	Quantity *int `json:"quantity,omitempty"`
	// PriceID is the associated price, if applicable.
	PriceID *string `json:"priceId,omitempty"`
}

// ListInvoicesParams are the query parameters for listing invoices.
type ListInvoicesParams struct {
	ListParams

	// CustomerID filters invoices by customer.
	CustomerID *string
}

// ListInvoiceLineItemsParams are the query parameters for listing invoice line items.
type ListInvoiceLineItemsParams struct {
	ListParams

	// InvoiceID filters line items by invoice.
	InvoiceID *string
}

// ---------------------------------------------------------------------------
// Payments
// ---------------------------------------------------------------------------

// Payment represents a payment transaction.
type Payment struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// CustomerID is the paying customer.
	CustomerID string `json:"customerId"`
	// InvoiceID is the associated invoice, if any.
	InvoiceID *string `json:"invoiceId,omitempty"`
	// Amount is the payment amount in the smallest currency unit.
	Amount int64 `json:"amount"`
	// Currency is the ISO 4217 currency code.
	Currency string `json:"currency"`
	// Status is the payment status (e.g. "succeeded", "failed", "refunded").
	Status string `json:"status"`
}

// RefundPaymentParams are the parameters for refunding a payment.
type RefundPaymentParams struct {
	// Amount is the amount to refund in the smallest currency unit.
	// Omit for a full refund.
	Amount *int64 `json:"amount,omitempty"`
	// Reason is an optional human-readable reason for the refund.
	Reason *string `json:"reason,omitempty"`
}

// ListPaymentsParams are the query parameters for listing payments.
type ListPaymentsParams struct {
	ListParams

	// CustomerID filters payments by customer.
	CustomerID *string
}

// ---------------------------------------------------------------------------
// Usage Events
// ---------------------------------------------------------------------------

// UsageEvent represents a single consumption record submitted against a usage meter.
type UsageEvent struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// Livemode indicates live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// UsageMeterID is the meter this event is recorded against.
	UsageMeterID string `json:"usageMeterId"`
	// CustomerID is the customer whose meter is incremented.
	CustomerID string `json:"customerId"`
	// Quantity is the consumption amount.
	Quantity float64 `json:"quantity"`
	// IdempotencyKey is the caller-supplied deduplication key.
	IdempotencyKey *string `json:"idempotencyKey,omitempty"`
	// EventAt is the optional Unix timestamp of when the event occurred.
	EventAt *int64 `json:"eventAt,omitempty"`
}

// CreateUsageEventParams are the parameters for recording a single usage event.
type CreateUsageEventParams struct {
	// UsageMeterID is the meter to record against.
	UsageMeterID string `json:"usageMeterId"`
	// CustomerID is the customer whose usage is incremented.
	CustomerID string `json:"customerId"`
	// Quantity is the consumption amount to record.
	Quantity float64 `json:"quantity"`
	// IdempotencyKey prevents duplicate recording if the request is retried.
	IdempotencyKey *string `json:"idempotencyKey,omitempty"`
	// EventAt is the optional timestamp of when the event occurred (Unix ms).
	EventAt *int64 `json:"eventAt,omitempty"`
}

// BulkCreateUsageEventsParams are the parameters for recording multiple usage events atomically.
type BulkCreateUsageEventsParams struct {
	// Events is the list of usage events to record.
	Events []CreateUsageEventParams `json:"events"`
}

// ListUsageEventsParams are the query parameters for listing usage events.
type ListUsageEventsParams struct {
	ListParams

	// UsageMeterID filters events by meter.
	UsageMeterID *string
	// CustomerID filters events by customer.
	CustomerID *string
}

// ---------------------------------------------------------------------------
// Usage Meters
// ---------------------------------------------------------------------------

// UsageMeterAggregation defines how a usage meter aggregates events.
type UsageMeterAggregation string

const (
	// UsageMeterAggregationSum sums all event quantities.
	UsageMeterAggregationSum UsageMeterAggregation = "sum"
	// UsageMeterAggregationCount counts the number of events.
	UsageMeterAggregationCount UsageMeterAggregation = "count"
	// UsageMeterAggregationMax takes the maximum event quantity.
	UsageMeterAggregationMax UsageMeterAggregation = "max"
)

// UsageMeter tracks consumption of a metered resource for usage-based billing.
type UsageMeter struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// Name is the meter's display name.
	Name string `json:"name"`
	// Aggregation is how the meter accumulates event quantities.
	Aggregation UsageMeterAggregation `json:"aggregation"`
}

// CreateUsageMeterParams are the parameters for creating a new usage meter.
type CreateUsageMeterParams struct {
	// Name is the meter display name.
	Name string `json:"name"`
	// Aggregation specifies how events are aggregated.
	Aggregation UsageMeterAggregation `json:"aggregation"`
}

// ---------------------------------------------------------------------------
// Discounts
// ---------------------------------------------------------------------------

// Discount represents a promotional price reduction.
type Discount struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// Name is the discount's display name.
	Name string `json:"name"`
	// AmountOff is the fixed amount to deduct (in the smallest currency unit).
	AmountOff *int64 `json:"amountOff,omitempty"`
	// PercentOff is the percentage to deduct (0–100).
	PercentOff *float64 `json:"percentOff,omitempty"`
	// Duration is the number of months the discount applies ("forever", "once", "repeating").
	Duration string `json:"duration"`
	// DurationMonths is the number of months if Duration is "repeating".
	DurationMonths *int `json:"durationMonths,omitempty"`
	// Active indicates whether the discount is currently usable.
	Active bool `json:"active"`
}

// CreateDiscountParams are the parameters for creating a new discount.
type CreateDiscountParams struct {
	// Name is the discount display name.
	Name string `json:"name"`
	// AmountOff is the fixed amount to deduct.
	AmountOff *int64 `json:"amountOff,omitempty"`
	// PercentOff is the percentage to deduct.
	PercentOff *float64 `json:"percentOff,omitempty"`
	// Duration controls how long the discount applies.
	Duration string `json:"duration"`
	// DurationMonths specifies the number of months for "repeating" duration.
	DurationMonths *int `json:"durationMonths,omitempty"`
	// Active sets the initial active state.
	Active *bool `json:"active,omitempty"`
}

// UpdateDiscountParams are the parameters for updating an existing discount.
type UpdateDiscountParams struct {
	// Name is the new display name.
	Name *string `json:"name,omitempty"`
	// AmountOff is the new fixed amount.
	AmountOff *int64 `json:"amountOff,omitempty"`
	// PercentOff is the new percentage.
	PercentOff *float64 `json:"percentOff,omitempty"`
	// Duration updates the duration type.
	Duration *string `json:"duration,omitempty"`
	// DurationMonths updates the repeating duration months.
	DurationMonths *int `json:"durationMonths,omitempty"`
	// Active updates the active flag.
	Active *bool `json:"active,omitempty"`
}

// ListDiscountsParams are the query parameters for listing discounts.
type ListDiscountsParams struct {
	ListParams
}

// ---------------------------------------------------------------------------
// Webhooks (management)
// ---------------------------------------------------------------------------

// WebhookEventType is a Flowglad webhook event type constant.
type WebhookEventType string

const (
	// WebhookEventTypeCustomerCreated fires when a customer is created.
	WebhookEventTypeCustomerCreated WebhookEventType = "customer.created"
	// WebhookEventTypeCustomerUpdated fires when a customer is updated.
	WebhookEventTypeCustomerUpdated WebhookEventType = "customer.updated"
	// WebhookEventTypePurchaseCompleted fires when a purchase is completed.
	WebhookEventTypePurchaseCompleted WebhookEventType = "purchase.completed"
	// WebhookEventTypePaymentFailed fires when a payment fails.
	WebhookEventTypePaymentFailed WebhookEventType = "payment.failed"
	// WebhookEventTypePaymentSucceeded fires when a payment succeeds.
	WebhookEventTypePaymentSucceeded WebhookEventType = "payment.succeeded"
	// WebhookEventTypeSubscriptionCreated fires when a subscription is created.
	WebhookEventTypeSubscriptionCreated WebhookEventType = "subscription.created"
	// WebhookEventTypeSubscriptionUpdated fires when a subscription is updated.
	WebhookEventTypeSubscriptionUpdated WebhookEventType = "subscription.updated"
	// WebhookEventTypeSubscriptionCanceled fires when a subscription is canceled.
	WebhookEventTypeSubscriptionCanceled WebhookEventType = "subscription.canceled"
	// WebhookEventTypeSyncEventsAvailable fires when sync events are available.
	WebhookEventTypeSyncEventsAvailable WebhookEventType = "sync.events_available"
)

// Webhook represents a registered webhook endpoint.
type Webhook struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// Name is the webhook's display name.
	Name string `json:"name"`
	// URL is the HTTPS endpoint that receives events.
	URL string `json:"url"`
	// FilterTypes are the event types this webhook subscribes to.
	FilterTypes []WebhookEventType `json:"filterTypes"`
	// Active indicates whether this webhook is currently receiving events.
	Active bool `json:"active"`
}

// CreateWebhookResult is returned when creating a webhook. It extends Webhook
// with the signing secret (shown only once at creation time).
type CreateWebhookResult struct {
	Webhook

	// Secret is the HMAC-SHA256 signing secret. Store this securely; it will not
	// be returned again.
	Secret string `json:"secret"` //nolint:gosec // webhook signing secret returned once at creation
}

// CreateWebhookParams are the parameters for registering a new webhook.
type CreateWebhookParams struct {
	// Name is the webhook display name.
	Name string `json:"name"`
	// URL is the HTTPS endpoint to deliver events to.
	URL string `json:"url"`
	// FilterTypes limits delivery to these event types. Omit to receive all events.
	FilterTypes []WebhookEventType `json:"filterTypes,omitempty"`
	// Active sets the initial active state (defaults to true).
	Active *bool `json:"active,omitempty"`
}

// UpdateWebhookParams are the parameters for updating an existing webhook.
type UpdateWebhookParams struct {
	// Name is the new display name.
	Name *string `json:"name,omitempty"`
	// URL is the new endpoint URL.
	URL *string `json:"url,omitempty"`
	// FilterTypes replaces the event type subscription list.
	FilterTypes []WebhookEventType `json:"filterTypes,omitempty"`
	// Active updates the active flag.
	Active *bool `json:"active,omitempty"`
}

// ListWebhooksParams are the query parameters for listing webhooks.
type ListWebhooksParams struct {
	ListParams
}

// ---------------------------------------------------------------------------
// Payment Methods
// ---------------------------------------------------------------------------

// PaymentMethod represents a saved payment instrument for a customer.
type PaymentMethod struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// CustomerID is the owning customer.
	CustomerID string `json:"customerId"`
	// Type is the payment method type (e.g. "card", "bank_account").
	Type string `json:"type"`
	// Last4 is the last four digits of the card or account number.
	Last4 *string `json:"last4,omitempty"`
	// Brand is the card network brand (e.g. "visa", "mastercard").
	Brand *string `json:"brand,omitempty"`
	// ExpMonth is the card expiry month.
	ExpMonth *int `json:"expMonth,omitempty"`
	// ExpYear is the card expiry year.
	ExpYear *int `json:"expYear,omitempty"`
}

// ListPaymentMethodsParams are the query parameters for listing payment methods.
type ListPaymentMethodsParams struct {
	ListParams

	// CustomerID filters payment methods by customer.
	CustomerID *string
}

// ---------------------------------------------------------------------------
// Pricing Models
// ---------------------------------------------------------------------------

// PricingModel represents a named collection of prices for a product offering.
type PricingModel struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// Name is the pricing model's display name.
	Name string `json:"name"`
	// Description is an optional description.
	Description *string `json:"description,omitempty"`
	// IsDefault indicates whether this is the organization's default pricing model.
	IsDefault bool `json:"isDefault"`
}

// CreatePricingModelParams are the parameters for creating a new pricing model.
type CreatePricingModelParams struct {
	// Name is the pricing model display name.
	Name string `json:"name"`
	// Description is an optional description.
	Description *string `json:"description,omitempty"`
}

// UpdatePricingModelParams are the parameters for updating an existing pricing model.
type UpdatePricingModelParams struct {
	// Name is the new display name.
	Name *string `json:"name,omitempty"`
	// Description is the new description.
	Description *string `json:"description,omitempty"`
}

// ListPricingModelsParams are the query parameters for listing pricing models.
type ListPricingModelsParams struct {
	ListParams
}

// ClonePricingModelParams are the parameters for cloning a pricing model.
type ClonePricingModelParams struct {
	// Name is an optional override for the cloned model's name.
	Name *string `json:"name,omitempty"`
	// Description is an optional override for the cloned model's description.
	Description *string `json:"description,omitempty"`
}

// SetupPricingModelParams are the parameters for setting up a pricing model.
type SetupPricingModelParams struct {
	// Name is the pricing model display name.
	Name string `json:"name"`
	// Description is an optional description.
	Description *string `json:"description,omitempty"`
}

// ---------------------------------------------------------------------------
// Features
// ---------------------------------------------------------------------------

// Feature represents a named product feature used for entitlement management.
type Feature struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// Name is the feature's display name.
	Name string `json:"name"`
	// Slug is the unique human-readable identifier used in entitlement checks.
	Slug string `json:"slug"`
	// Description is an optional description.
	Description *string `json:"description,omitempty"`
}

// CreateFeatureParams are the parameters for creating a new feature.
type CreateFeatureParams struct {
	// Name is the feature display name.
	Name string `json:"name"`
	// Slug is the unique human-readable identifier.
	Slug string `json:"slug"`
	// Description is an optional description.
	Description *string `json:"description,omitempty"`
}

// UpdateFeatureParams are the parameters for updating an existing feature.
type UpdateFeatureParams struct {
	// Name is the new display name.
	Name *string `json:"name,omitempty"`
	// Slug is the new unique identifier.
	Slug *string `json:"slug,omitempty"`
	// Description is the new description.
	Description *string `json:"description,omitempty"`
}

// ListFeaturesParams are the query parameters for listing features.
type ListFeaturesParams struct {
	ListParams
}

// AddProductFeatureParams are the parameters for adding a feature to a product.
type AddProductFeatureParams struct {
	// FeatureID is the feature to add.
	FeatureID string `json:"featureId"`
}

// ---------------------------------------------------------------------------
// Resources (Entitlements)
// ---------------------------------------------------------------------------

// Resource represents an instance of a feature assigned for entitlement management.
type Resource struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt is the last-update timestamp (Unix ms).
	UpdatedAt int64 `json:"updatedAt"`
	// Livemode indicates live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// Name is the resource's display name.
	Name string `json:"name"`
	// FeatureID is the feature this resource is an instance of.
	FeatureID string `json:"featureId"`
	// CustomerID is the customer this resource is assigned to, if any.
	CustomerID *string `json:"customerId,omitempty"`
	// Quantity is the resource quantity (e.g. seat count).
	Quantity *float64 `json:"quantity,omitempty"`
}

// ResourceClaim represents a customer's claim on a resource.
type ResourceClaim struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// ResourceID is the resource being claimed.
	ResourceID string `json:"resourceId"`
	// CustomerID is the claiming customer.
	CustomerID string `json:"customerId"`
	// CreatedAt is the Unix timestamp when the claim was created.
	CreatedAt int64 `json:"createdAt"`
}

// ResourceUsage represents a usage record for a resource.
type ResourceUsage struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// ResourceID is the resource this usage is recorded for.
	ResourceID string `json:"resourceId"`
	// Quantity is the usage amount.
	Quantity float64 `json:"quantity"`
	// CreatedAt is the Unix timestamp when the usage was recorded.
	CreatedAt int64 `json:"createdAt"`
}

// CreateResourceParams are the parameters for creating a new resource.
type CreateResourceParams struct {
	// Name is the resource display name.
	Name string `json:"name"`
	// FeatureID is the feature this resource is an instance of.
	FeatureID string `json:"featureId"`
	// CustomerID is the customer to assign this resource to.
	CustomerID *string `json:"customerId,omitempty"`
	// Quantity is the resource quantity.
	Quantity *float64 `json:"quantity,omitempty"`
}

// UpdateResourceParams are the parameters for updating an existing resource.
type UpdateResourceParams struct {
	// Name is the new display name.
	Name *string `json:"name,omitempty"`
	// Quantity is the new quantity.
	Quantity *float64 `json:"quantity,omitempty"`
}

// ListResourcesParams are the query parameters for listing resources.
type ListResourcesParams struct {
	ListParams
}

// ClaimResourceParams are the parameters for claiming a resource for a customer.
type ClaimResourceParams struct {
	// CustomerID is the customer claiming the resource.
	CustomerID string `json:"customerId"`
}

// ReleaseResourceParams are the parameters for releasing a resource claim.
type ReleaseResourceParams struct {
	// CustomerID is the customer releasing the resource.
	CustomerID string `json:"customerId"`
}

// ---------------------------------------------------------------------------
// API Keys
// ---------------------------------------------------------------------------

// APIKey represents a Flowglad API key (read-only; create via the dashboard).
type APIKey struct {
	// ID is the unique Flowglad identifier.
	ID string `json:"id"`
	// CreatedAt is the creation timestamp (Unix ms).
	CreatedAt int64 `json:"createdAt"`
	// Livemode indicates live or test mode.
	Livemode bool `json:"livemode"`
	// OrganizationID is the owning Flowglad organization.
	OrganizationID string `json:"organizationId"`
	// Name is the key's display name.
	Name string `json:"name"`
	// LastFour is the last four characters of the key (for display purposes).
	LastFour string `json:"lastFour"`
	// Active indicates whether this key is currently active.
	Active bool `json:"active"`
}
