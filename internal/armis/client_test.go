package armis

import (
	"context"
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	// Basic sanity checks.
	if client.apiURL == "" {
		t.Fatalf("client apiURL should not be empty")
	}
	if client.apiKey == "" {
		t.Fatalf("client apiKey should not be empty")
	}

	// Verify that the helper sets an auth token by attempting to build a request.
	req, err := client.newRequest(context.Background(), http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatalf("newRequest failed: %v", err)
	}

	if got := req.Header.Get("Authorization"); got == "" {
		t.Fatalf("expected Authorization header to be set, got empty string")
	}
}
