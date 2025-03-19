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

func TestGettingUsers(t *testing.T) {
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

	// Attempt to get all users
	ctx := context.Background()
	response, err := client.GetUsers(ctx)
	if err != nil {
		t.Errorf("Error getting users: %s", err)
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

func TestCreatingUser(t *testing.T) {
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

	// Create a new user with necessary fields
	newUser := UserSettings{
		Name:     "Test User",
		Username: "testuser",
		Email:    "test.user@1898andco.io",
		RoleAssignment: []RoleAssignment{
			{
				Name:  []string{"Read Only"},
				Sites: []string{"Lab"},
			},
		},
	}

	// Call CreateUser to create the user
	ctx := context.Background()
	response, err := client.CreateUser(ctx, newUser)
	if err != nil {
		t.Errorf("Error creating user: %s", err)
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

	if response.Name != "Test User" {
		t.Errorf("Expected user name to be 'Test User', got '%s'", response.Name)
	}
	if response.Email != "test.user@1898andco.io" {
		t.Errorf("Expected email to be 'test.user@1898andco.io', got '%s'", response.Email)
	}
}

func TestGettingUser(t *testing.T) {
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

	// Attempt to get a user by email
	ctx := context.Background()
	response, err := client.GetUser(ctx, "test.user@1898andco.io")
	if err != nil {
		t.Errorf("Error getting user: %s", err)
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

func TestUpdatingUser(t *testing.T) {
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

	// Update the user with new information
	updatedUser := UserSettings{
		Name: "Test User",
		// Updated field
		Username: "testupdateduser",
		Email:    "test.user@1898andco.io",
		RoleAssignment: []RoleAssignment{
			{
				// Updated field
				Name:  []string{"Admin"},
				Sites: []string{"Lab"},
			},
		},
	}

	// Update the user
	ctx := context.Background()
	response, err := client.UpdateUser(ctx, updatedUser, "test.user@1898andco.io")
	if err != nil {
		t.Errorf("Error updating user: %s", err)
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

	if response.Username != "testupdateduser" {
		t.Errorf("Expected user name to be 'testupdateduser', got '%s'", response.Username)
	}
	if response.RoleAssignment[0].Name[0] != "Admin" {
		t.Errorf("Expected role to be 'Admin', got '%s'", response.RoleAssignment[0].Name[0])
	}
	if response.Email != "test.user@1898andco.io" {
		t.Errorf("Expected email to be 'test.user@1898andco.io', got '%s'", response.Email)
	}
}

func TestDeletingUsers(t *testing.T) {
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

	// Attempt to delete test user
	ctx := context.Background()
	success, err := client.DeleteUser(ctx, "test.user@1898andco.io")
	if err != nil {
		t.Errorf("Error deleting user: %s", err)
	}

	// Log the response
	if !success {
		log.Error("Failed to delete user.")
	} else {
		log.Info("Successfully deleted user.")
	}
}
