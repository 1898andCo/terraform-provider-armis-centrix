// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Authenticate retrieves a temporary access token using the API key.
func (c *Client) Authenticate(apiKey string) error {
	if apiKey == "" {
		return errors.New("API key is required")
	}

	// Check if we already have a valid access token
	if c.AccessToken != "" && c.AccessTokenExpiration.Before(time.Now()) {
		return nil
	}

	// Prepare the form data for the request
	form := url.Values{}
	form.Set("secret_key", apiKey)

	// Create the POST request
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/%s/access_token/", c.ApiUrl, c.ApiVersion),
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create authentication request: %w", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Send the request and get the response
	body, err := c.doRequest(req)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Parse the response JSON
	var authResponse AuthResponse
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return fmt.Errorf("failed to parse authentication response: %w", err)
	}

	// Check if the response indicates success
	if !authResponse.Success {
		return errors.New("authentication failed: API response indicates failure")
	}

	// Store the access token and expiration in the client
	c.AccessToken = authResponse.Data.AccessToken
	c.AccessTokenExpiration, err = time.Parse(time.RFC3339Nano, authResponse.Data.ExpirationUtc)

	// Expire 5 minutes early to ensure we don't ever get an invalid token
	c.AccessTokenExpiration = c.AccessTokenExpiration.Add(-time.Minute * 5)
	if err != nil {
		return fmt.Errorf("failed to parse expiration time: %w", err)
	}

	return nil
}
