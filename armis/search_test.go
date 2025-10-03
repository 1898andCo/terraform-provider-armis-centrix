// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestGetSearch(t *testing.T) {
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

func TestGetSearchRequiresAQL(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.GetSearch(context.Background(), " ", false, false); !errors.Is(err, ErrSearchAQL) {
		t.Fatalf("expected ErrSearchAQL, got %v", err)
	}
}
