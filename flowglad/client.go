package flowglad

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// ErrEmptyAPIKey is returned by NewClient when an empty API key is provided.
var ErrEmptyAPIKey = errors.New("flowglad: apiKey must not be empty")

const defaultBaseURL = "https://app.flowglad.com/api/v1"

// Client is the Flowglad SDK entry point. It is immutable after construction
// and safe for concurrent use from multiple goroutines.
//
// Access resources via the typed sub-client fields:
//
//	client.Customers.Create(ctx, params)
//	client.Subscriptions.List(ctx, nil)
type Client struct {
	// Customers provides methods for the Customers API.
	Customers *CustomersClient
	// Subscriptions provides methods for the Subscriptions API.
	Subscriptions *SubscriptionsClient
	// Products provides methods for the Products API.
	Products *ProductsClient
	// Prices provides methods for the Prices API.
	Prices *PricesClient
	// CheckoutSessions provides methods for the Checkout Sessions API.
	CheckoutSessions *CheckoutSessionsClient
	// Invoices provides methods for the Invoices API.
	Invoices *InvoicesClient
	// InvoiceLineItems provides methods for the Invoice Line Items API.
	InvoiceLineItems *InvoiceLineItemsClient
	// Payments provides methods for the Payments API.
	Payments *PaymentsClient
	// PaymentMethods provides methods for the Payment Methods API.
	PaymentMethods *PaymentMethodsClient
	// UsageEvents provides methods for the Usage Events API.
	UsageEvents *UsageEventsClient
	// UsageMeters provides methods for the Usage Meters API.
	UsageMeters *UsageMetersClient
	// Discounts provides methods for the Discounts API.
	Discounts *DiscountsClient
	// Webhooks provides methods for the Webhooks management API.
	Webhooks *WebhooksClient
	// PricingModels provides methods for the Pricing Models API.
	PricingModels *PricingModelsClient
	// Features provides methods for the Features API.
	Features *FeaturesClient
	// Resources provides methods for the Resources (Entitlements) API.
	Resources *ResourcesClient
	// APIKeys provides methods for the API Keys API.
	APIKeys *APIKeysClient
}

// NewClient constructs an immutable, goroutine-safe Flowglad SDK client.
//
// apiKey is required and must be non-empty; passing an empty string returns an
// error. All other settings have sensible defaults and can be overridden with
// functional options.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, ErrEmptyAPIKey
	}

	cfg := &clientConfig{
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
		logger:     slog.Default(),
		retry:      defaultRetryPolicy(),
	}
	for _, o := range opts {
		o(cfg)
	}
	if cfg.logger == nil {
		cfg.logger = slog.Default()
	}

	var retryPolicy *internal.BackoffPolicy
	if !cfg.noRetry {
		retryPolicy = &internal.BackoffPolicy{
			MaxAttempts:    cfg.retry.MaxAttempts,
			InitialBackoff: cfg.retry.InitialBackoff,
			MaxBackoff:     cfg.retry.MaxBackoff,
			Multiplier:     cfg.retry.Multiplier,
			JitterFactor:   cfg.retry.JitterFactor,
		}
	}

	httpCfg := &internal.HTTPConfig{
		BaseURL:    cfg.baseURL,
		APIKey:     apiKey,
		HTTPClient: cfg.httpClient,
		Logger:     cfg.logger,
		Retry:      retryPolicy,
		NoRetry:    cfg.noRetry,
	}

	return &Client{
		Customers:        &CustomersClient{cfg: httpCfg},
		Subscriptions:    &SubscriptionsClient{cfg: httpCfg},
		Products:         &ProductsClient{cfg: httpCfg},
		Prices:           &PricesClient{cfg: httpCfg},
		CheckoutSessions: &CheckoutSessionsClient{cfg: httpCfg},
		Invoices:         &InvoicesClient{cfg: httpCfg},
		InvoiceLineItems: &InvoiceLineItemsClient{cfg: httpCfg},
		Payments:         &PaymentsClient{cfg: httpCfg},
		PaymentMethods:   &PaymentMethodsClient{cfg: httpCfg},
		UsageEvents:      &UsageEventsClient{cfg: httpCfg},
		UsageMeters:      &UsageMetersClient{cfg: httpCfg},
		Discounts:        &DiscountsClient{cfg: httpCfg},
		Webhooks:         &WebhooksClient{cfg: httpCfg},
		PricingModels:    &PricingModelsClient{cfg: httpCfg},
		Features:         &FeaturesClient{cfg: httpCfg},
		Resources:        &ResourcesClient{cfg: httpCfg},
		APIKeys:          &APIKeysClient{cfg: httpCfg},
	}, nil
}
