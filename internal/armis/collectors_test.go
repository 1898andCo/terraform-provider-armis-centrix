// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"

	log "github.com/charmbracelet/log"
)

// integrationClient builds a live API client or skips the test suite if the
// required environment variables are not present.
func integrationClient(t *testing.T) *Client {
	apiURL := os.Getenv("ARMIS_API_URL")
	apiKey := os.Getenv("ARMIS_API_KEY")
	if apiURL == "" || apiKey == "" {
		t.Skip("ARMIS_API_URL and ARMIS_API_KEY must be set for integration tests")
	}

	client, err := NewClient(apiKey, WithAPIURL(apiURL))
	if err != nil {
		t.Fatalf("create client: %v", err)
	}
	return client
}

func prettyPrint(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Info("marshal response: %v", err)
		return
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, data, "", "  "); err == nil {
		log.Info("\n=== Response ===\n%s", buf.String())
	}
}

func TestCreatingCollector(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	ctx := context.Background()
	payload := CreateCollectorSettings{
		Name:           "Test Collector",
		DeploymentType: "OVA",
	}

	res, err := client.CreateCollector(ctx, payload)
	if err != nil {
		t.Fatalf("create collector: %v", err)
	}
	prettyPrint(res)
}

func TestGettingCollectors(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	res, err := client.GetCollectors(context.Background())
	if err != nil {
		t.Fatalf("get collectors: %v", err)
	}
	prettyPrint(res)
}

func TestGettingCollectorByID(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	const id = "8153" // TODO: make dynamic or env-driven
	res, err := client.GetCollectorByID(context.Background(), id)
	if err != nil {
		t.Fatalf("get collector by id: %v", err)
	}
	prettyPrint(res)
}

func TestUpdatingCollector(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	payload := UpdateCollectorSettings{
		Name:           "Test Collector Updated",
		DeploymentType: "OVA",
	}

	const id = "8158"
	res, err := client.UpdateCollector(context.Background(), id, payload)
	if err != nil {
		t.Fatalf("update collector: %v", err)
	}
	prettyPrint(res)
}

func TestDeletingCollector(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	const id = "8158"
	ok, err := client.DeleteCollector(context.Background(), id)
	if err != nil {
		t.Fatalf("delete collector: %v", err)
	}
	if !ok {
		t.Fatalf("collector %s not deleted", id)
	}
	log.Info("collector %s deleted", id)
}
