// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

// Package sweep provides utility functions for configuring and running
// test sweepers for Armis resources in Terraform.
package sweep

import (
	"fmt"
	"log"
	"os"

	"github.com/1898andCo/terraform-provider-armis-centrix/internal/armis"
)

// ConfigureSweeperClient initializes an Armis client using environment variables
// It returns the client and an error if initialization fails.
func ConfigureSweeperClient(name string) (*armis.Client, error) {
	// Get configuration from environment variables
	apiKey := os.Getenv("ARMIS_API_KEY")
	if apiKey == "" {
		log.Printf("[INFO] Skipping %s sweeper - ARMIS_API_KEY not set", name)
		return nil, armis.ErrGetKey
	}

	apiURL := os.Getenv("ARMIS_API_URL")
	if apiURL == "" {
		log.Printf("[INFO] Skipping %s sweeper - ARMIS_API_URL not set", name)
		return nil, armis.ErrGetURL
	}

	// Initialize the client
	client, err := armis.NewClient(
		apiKey,
		armis.WithAPIURL(apiURL),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating Armis client: %w", err)
	}

	return client, nil
}
