// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	log "github.com/charmbracelet/log"
)

func TestGettingSites(t *testing.T) {
	// Initialize the client
	options := Client{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	log.Info("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	// Attempt to get all sites
	response, err := client.GetSites()
	if err != nil {
		t.Errorf("Error getting sites: %s", err)
	}

	// Log the response
	if response != nil {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Info("Error marshaling server response: %s\n", err)
		}

		// Attempt to pretty-print the JSON
		var prettyResponse bytes.Buffer
		if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
			log.Info("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
		} else {
			log.Info("Failed to pretty-print JSON.")
		}
	} else {
		log.Info("No response received from server.")
	}
}
