package flowgladtest

// fixtureTimestamp is a representative epoch-millisecond timestamp used in test fixtures.
const fixtureTimestamp = 1_700_000_000_000

// CustomerFixture returns a pre-built customer JSON response envelope suitable
// for passing to RespondWith.
func CustomerFixture(id, name, email, externalID string) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"customer": map[string]any{
				"id":             id,
				"createdAt":      fixtureTimestamp,
				"updatedAt":      fixtureTimestamp,
				"livemode":       false,
				"organizationId": "org_test",
				"name":           name,
				"email":          email,
				"externalId":     externalID,
				"archived":       false,
			},
		},
	}
}

// SubscriptionFixture returns a pre-built subscription JSON response envelope.
func SubscriptionFixture(id, customerID, priceID, status string) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"subscription": map[string]any{
				"id":             id,
				"createdAt":      fixtureTimestamp,
				"updatedAt":      fixtureTimestamp,
				"livemode":       false,
				"organizationId": "org_test",
				"customerId":     customerID,
				"priceId":        priceID,
				"status":         status,
			},
		},
	}
}

// CheckoutSessionFixture returns a pre-built checkout session JSON response envelope.
func CheckoutSessionFixture(id, url, status string) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"checkoutSession": map[string]any{
				"id":             id,
				"createdAt":      fixtureTimestamp,
				"updatedAt":      fixtureTimestamp,
				"livemode":       false,
				"organizationId": "org_test",
				"url":            url,
				"status":         status,
			},
		},
	}
}

// CustomerListFixture returns a paginated customer list response envelope.
func CustomerListFixture(customers []map[string]any, hasMore bool, nextCursor string) map[string]any {
	return map[string]any{
		"data":          customers,
		"hasMore":       hasMore,
		"nextCursor":    nextCursor,
		"currentCursor": "",
		"total":         len(customers),
	}
}

// ErrorFixture returns a Flowglad-shaped error response body.
func ErrorFixture(code, message string) map[string]any {
	return map[string]any{
		"code":    code,
		"message": message,
		"issues":  []any{},
	}
}
