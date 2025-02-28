// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestCreatingCollector(t *testing.T) {
	// Initialize the client
	options := ClientOptions{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	fmt.Printf("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	// Attempt to create a collector
	collector := CreateCollectorSettings{
		Name:           "Test Collector",
		DeploymentType: "OVA",
	}

	response, err := client.CreateCollector(collector)
	if err != nil {
		t.Errorf("Error creating collector: %s", err)
	}

	// Log the response
	if response != nil {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("Error marshaling server response: %s\n", err)
		}

		// Attempt to pretty-print the JSON
		var prettyResponse bytes.Buffer
		if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
			fmt.Printf("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
		} else {
			fmt.Println("Failed to pretty-print JSON.")
		}
	} else {
		fmt.Println("No response received from server.")
	}
}

func TestGettingCollectors(t *testing.T) {
	// Initialize the client
	options := ClientOptions{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	fmt.Printf("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	// Attempt to get all collectors
	response, err := client.GetCollectors()
	if err != nil {
		t.Errorf("Error getting sites: %s", err)
	}

	// Log the response
	if response != nil {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("Error marshaling server response: %s\n", err)
		}

		// Attempt to pretty-print the JSON
		var prettyResponse bytes.Buffer
		if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
			fmt.Printf("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
		} else {
			fmt.Println("Failed to pretty-print JSON.")
		}
	} else {
		fmt.Println("No response received from server.")
	}
}

func TestGettingCollectorByID(t *testing.T) {
	// Initialize the client
	options := ClientOptions{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	fmt.Printf("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	// Attempt to get all collectors
	response, err := client.GetCollectorByID("8153")
	if err != nil {
		t.Errorf("Error getting sites: %s", err)
	}

	// Log the response
	if response != nil {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("Error marshaling server response: %s\n", err)
		}

		// Attempt to pretty-print the JSON
		var prettyResponse bytes.Buffer
		if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
			fmt.Printf("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
		} else {
			fmt.Println("Failed to pretty-print JSON.")
		}
	} else {
		fmt.Println("No response received from server.")
	}
}

func TestUpdatingCollector(t *testing.T) {
	// Initialize the client
	options := ClientOptions{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	fmt.Printf("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	// Attempt to update a collector
	collector := UpdateCollectorSettings{
		Name:           "Test Collector",
		DeploymentType: "OVA",
	}

	response, err := client.UpdateCollector("8158", collector)
	if err != nil {
		t.Errorf("Error updating collector: %s", err)
	}

	// Log the response
	if response != nil {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("Error marshaling server response: %s\n", err)
		}

		// Attempt to pretty-print the JSON
		var prettyResponse bytes.Buffer
		if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
			fmt.Printf("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
		} else {
			fmt.Println("Failed to pretty-print JSON.")
		}
	} else {
		fmt.Println("No response received from server.")
	}
}

func TestDeletingCollector(t *testing.T) {
	// Initialize the client
	options := ClientOptions{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	fmt.Printf("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	// Attempt to delete test user
	success, err := client.DeleteCollector("8158")
	if err != nil {
		t.Errorf("Error deleting user: %s", err)
	}

	// Log the response
	if !success {
		fmt.Printf("Failed to delete user.")
	} else {
		fmt.Printf("Successfully deleted user.")
	}
}
