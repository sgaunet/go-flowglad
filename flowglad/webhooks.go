package flowglad

import (
	"context"

	"github.com/sgaunet/go-flowglad/flowglad/internal"
)

// WebhooksClient provides methods for the Flowglad Webhooks management API.
type WebhooksClient struct {
	cfg *internal.HTTPConfig
}

// webhookResponse is the JSON envelope for a single webhook response.
type webhookResponse struct {
	Data struct {
		Webhook Webhook `json:"webhook"`
	} `json:"data"`
}

// createWebhookResponse is the JSON envelope for the webhook create response
// (includes the one-time secret).
type createWebhookResponse struct {
	Data struct {
		Webhook CreateWebhookResult `json:"webhook"`
	} `json:"data"`
}

// Create registers a new webhook endpoint. The returned CreateWebhookResult
// includes the signing secret; store it securely as it will not be returned again.
func (c *WebhooksClient) Create(ctx context.Context, p *CreateWebhookParams) (*CreateWebhookResult, error) {
	resp, err := internal.Do[createWebhookResponse](ctx, c.cfg, "POST", "/webhooks", p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Webhook, nil
}

// Get retrieves a webhook by ID.
func (c *WebhooksClient) Get(ctx context.Context, id string) (*Webhook, error) {
	resp, err := internal.Do[webhookResponse](ctx, c.cfg, "GET", "/webhooks/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Webhook, nil
}

// Update updates an existing webhook configuration.
func (c *WebhooksClient) Update(ctx context.Context, id string, p *UpdateWebhookParams) (*Webhook, error) {
	resp, err := internal.Do[webhookResponse](ctx, c.cfg, "PUT", "/webhooks/"+id, p, false)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Webhook, nil
}
