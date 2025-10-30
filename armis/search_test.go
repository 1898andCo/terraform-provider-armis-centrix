// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestGetAlertSearch(t *testing.T) {
	t.Parallel()

	const (
		aql           = "in:alerts status:Open"
		includeSample = true
		includeTotal  = true
	)

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/search/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}

			values := r.URL.Query()
			if got := values.Get("aql"); got != aql {
				t.Fatalf("unexpected aql: %q", got)
			}
			if got := values.Get("includeSample"); got != "true" {
				t.Fatalf("unexpected includeSample: %q", got)
			}
			if got := values.Get("includeTotal"); got != "true" {
				t.Fatalf("unexpected includeTotal: %q", got)
			}

			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"count": 1,
					"next":  nil,
					"prev":  nil,
					"total": 1,
					"results": []map[string]any{{
						"title":  "Example Alert",
						"status": "Open",
					}},
				},
			})
		},
	})
	defer cleanup()

	res, err := client.GetSearch(context.Background(), aql, includeSample, includeTotal)
	if err != nil {
		t.Fatalf("get search: %v", err)
	}
	if res.Total != 1 || len(res.Results) != 1 || res.Results[0].Title != "Example Alert" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestGetActivitySearch(t *testing.T) {
	t.Parallel()

	const (
		aql           = "in:activity status:Open"
		includeSample = true
		includeTotal  = true
	)

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/search/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}

			values := r.URL.Query()
			if got := values.Get("aql"); got != aql {
				t.Fatalf("unexpected aql: %q", got)
			}
			if got := values.Get("includeSample"); got != "true" {
				t.Fatalf("unexpected includeSample: %q", got)
			}
			if got := values.Get("includeTotal"); got != "true" {
				t.Fatalf("unexpected includeTotal: %q", got)
			}

			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"count": 1,
					"next":  nil,
					"prev":  nil,
					"total": 1,
					"results": []map[string]any{{
						"title":         "Example Activity",
						"activityUUIDs": []string{"11111111-2222-3333-4444-555555555555"},
						"time":          "2025-10-23T04:48:10.804405Z",
						"sourceEndpoints": []map[string]any{
							{
								"id":   97,
								"ip":   []string{"192.168.2.2"},
								"name": "192.168.2.2",
							},
						},
					}},
				},
			})
		},
	})
	defer cleanup()

	res, err := client.GetSearch(context.Background(), aql, includeSample, includeTotal)
	if err != nil {
		t.Fatalf("get search: %v", err)
	}
	if res.Total != 1 || len(res.Results) != 1 || res.Results[0].Title != "Example Activity" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestSearchEndpointIDUnmarshal(t *testing.T) {
	payload := []byte(`{"sourceEndpoints":[{"id":123},{"id":"456"}]}`)
	var res struct {
		Source []SearchEndpoint `json:"sourceEndpoints"`
	}
	if err := json.Unmarshal(payload, &res); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := res.Source[0].ID; got != SearchEndpointID("123") {
		t.Fatalf("expected 123, got %q", got)
	}
	if got := res.Source[1].ID; got != SearchEndpointID("456") {
		t.Fatalf("expected 456, got %q", got)
	}
}

func TestSearchEndpointIPsUnmarshal(t *testing.T) {
	t.Parallel()

	t.Run("slice input", func(t *testing.T) {
		payload := []byte(`{"sourceEndpoints":[{"ip":["10.0.0.1","fe80::1"]}]}`)
		var res struct {
			Source []SearchEndpoint `json:"sourceEndpoints"`
		}
		if err := json.Unmarshal(payload, &res); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := []string(res.Source[0].IP)
		if len(got) != 2 || got[0] != "10.0.0.1" || got[1] != "fe80::1" {
			t.Fatalf("unexpected ips: %#v", got)
		}
	})

	t.Run("string input", func(t *testing.T) {
		payload := []byte(`{"sourceEndpoints":[{"ip":"10.0.0.1"}]}`)
		var res struct {
			Source []SearchEndpoint `json:"sourceEndpoints"`
		}
		if err := json.Unmarshal(payload, &res); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := []string(res.Source[0].IP)
		if len(got) != 1 || got[0] != "10.0.0.1" {
			t.Fatalf("unexpected ips: %#v", got)
		}
	})
}

func TestGetSearchRequiresAQL(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.GetSearch(context.Background(), " ", false, false); !errors.Is(err, ErrSearchAQL) {
		t.Fatalf("expected ErrSearchAQL, got %v", err)
	}
}
