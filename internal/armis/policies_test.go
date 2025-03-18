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

func TestCreatingPolicy(t *testing.T) {
	// Initialize the client with environment variables
	options := Client{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	log.Info("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	// Create a new policy with necessary fields
	newPolicy := PolicySettings{
		Name:              "Test Policy",
		Description:       "This is a test policy",
		IsEnabled:         false,
		Labels:            []string{"Security"},
		MitreAttackLabels: []string{"Enterprise.TA0009.T1056.001", "Enterprise.TA0009.T1056.004"},
		RuleType:          "ACTIVITY",
		Actions: []Action{
			{
				Type: "alert",
				Params: Params{
					Severity: "high",
					Title:    "Test Security Alert",
					Type:     "Security - Threat",
					Consolidation: Consolidation{
						Amount: 1,
						Unit:   "Days",
					},
				},
			},
		},
		Rules: Rules{
			And: []any{
				"protocol:BMS",
				Rules{
					Or: []any{
						"content:(iPhone)",
						"content:(Android)",
					},
				},
			},
		},
	}

	// Pretty print with indentation the policy before sending
	newPolicyJSON, err := json.Marshal(newPolicy)
	if err != nil {
		t.Fatalf("Error marshaling policy: %s", err)
	}
	log.Info("\n=== Policy to Create ===\n%s\n", string(newPolicyJSON))

	// Call CreatePolicy to create the policy
	response, err := client.CreatePolicy(newPolicy)
	if err != nil {
		t.Errorf("Error creating policy: %s", err)
		return
	}

	// Log the response
	if response != nil {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("Error marshaling server response: %s", err)
		}

		// Attempt to pretty-print the JSON for better readability
		var prettyResponse bytes.Buffer
		if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
			log.Info("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
		} else {
			t.Log("Failed to pretty-print JSON.")
		}
	} else {
		t.Log("No response received from server.")
	}
}

func TestCreatingTagPolicy(t *testing.T) {
	// Initialize the client with environment variables
	options := Client{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	log.Info("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	// Create a new policy with necessary fields
	newPolicy := PolicySettings{
		Name:              "Test Tag Policy",
		Description:       "This is a test tag policy",
		IsEnabled:         false,
		Labels:            []string{"Security"},
		MitreAttackLabels: []string{"Enterprise.TA0009.T1056.001", "Enterprise.TA0009.T1056.004"},
		RuleType:          "ACTIVITY",
		Actions: []Action{
			{
				Type: "tag",
				Params: Params{
					Endpoint: "ALL",
					Tags:     []string{"Agent and Scanner Gaps"},
				},
			},
		},
		Rules: Rules{
			And: []any{
				"protocol:BMS",
				Rules{
					Or: []any{
						"content:(iPhone)",
						"content:(Android)",
					},
				},
			},
		},
	}

	// Pretty print with indentation the policy before sending
	newPolicyJSON, err := json.Marshal(newPolicy)
	if err != nil {
		t.Fatalf("Error marshaling policy: %s", err)
	}
	log.Info("\n=== Policy to Create ===\n%s\n", string(newPolicyJSON))

	// Call CreatePolicy to create the policy
	response, err := client.CreatePolicy(newPolicy)
	if err != nil {
		t.Errorf("Error creating policy: %s", err)
		return
	}

	// Log the response
	if response != nil {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("Error marshaling server response: %s", err)
		}

		// Attempt to pretty-print the JSON for better readability
		var prettyResponse bytes.Buffer
		if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
			log.Info("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
		} else {
			t.Log("Failed to pretty-print JSON.")
		}
	} else {
		t.Log("No response received from server.")
	}
}

func TestGettingPolicy(t *testing.T) {
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

	// Attempt to get policy
	response, err := client.GetPolicy("76884")
	if err != nil {
		t.Errorf("Error getting policy: %s", err)
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

func TestUpdatingPolicy(t *testing.T) {
	// Initialize the client with environment variables
	options := Client{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	log.Info("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	// Update a policy with necessary fields
	updatedPolicy := PolicySettings{
		Name:              "Test Policy",
		Description:       "This is an updated test policy",
		IsEnabled:         true,
		Labels:            []string{"Security"},
		MitreAttackLabels: []string{"Enterprise.TA0009.T1056.001", "Enterprise.TA0009.T1056.004"},
		RuleType:          "ACTIVITY",
		Actions: []Action{
			{
				Type: "alert",
				Params: Params{
					Severity: "high",
					Title:    "Test Security Alert",
					Type:     "Security - Threat",
					Consolidation: Consolidation{
						Amount: 1,
						Unit:   "Days",
					},
				},
			},
		},
		Rules: Rules{
			And: []any{
				"protocol:BMS",
				Rules{
					Or: []any{
						"content:(iPhone)",
						"content:(Android)",
					},
				},
			},
		},
	}

	// Pretty print with indentation the policy before sending
	updatedPolicyJSON, err := json.Marshal(updatedPolicy)
	if err != nil {
		t.Fatalf("Error marshaling policy: %s", err)
	}
	log.Info("\n=== Policy to Update ===\n%s\n", string(updatedPolicyJSON))

	// Call UpdatePolicy to update the policy
	response, err := client.UpdatePolicy(updatedPolicy, "76700")
	if err != nil {
		t.Errorf("Error updating policy: %s", err)
		return
	}

	// Log the response
	if response != nil {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("Error marshaling server response: %s", err)
		}

		// Attempt to pretty-print the JSON for better readability
		var prettyResponse bytes.Buffer
		if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
			log.Info("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
		} else {
			t.Log("Failed to pretty-print JSON.")
		}
	} else {
		t.Log("No response received from server.")
	}
}

func TestDeletingPolicy(t *testing.T) {
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

	// Attempt to delete test policy
	success, err := client.DeletePolicy("76700")
	if err != nil {
		t.Errorf("Error deleting policy: %s", err)
	}

	// Log the response
	if !success {
		log.Info("Failed to delete policy.")
	} else {
		log.Info("Successfully deleted policy.")
	}
}
