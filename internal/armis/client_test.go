// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	// Initialize the client using the environment variables
	options := Client{
		ApiUrl:     os.Getenv("ARMIS_API_URL"),
		ApiKey:     os.Getenv("ARMIS_API_KEY"),
		APIVersion: "v1",
	}

	client, err := NewClient(options)

	// Assertions
	assert.NoError(t, err, "Expected no error when initializing client")
	assert.NotNil(t, client, "Client should not be nil")
	assert.Equal(t, "https://lab-1898andco.armis.com", client.ApiUrl, "Client API URL should match the environment variable")
	assert.Equal(t, "v1", client.APIVersion, "Client API version should match the environment variable")
	assert.NotEmpty(t, client.AccessToken, "Client should have an access token after authentication")
}
