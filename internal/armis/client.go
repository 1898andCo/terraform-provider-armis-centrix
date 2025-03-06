// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	armisApiUrl     = "https://api.armis.com"
	armisApiVersion = "v1"
)

type Client struct {
	ApiUrl                string
	ApiKey                string
	ApiVersion            string
	AccessToken           string
	AccessTokenExpiration time.Time

	HTTPClient *http.Client
}

// NewClient returns a new Armis client and authenticates to the Armis API endpoint.
func NewClient(options Client) (*Client, error) {
	apiUrl := armisApiUrl
	apiVersion := armisApiVersion

	if apiUrl != "" {
		apiUrl = options.ApiUrl
	}

	if options.ApiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	if options.ApiVersion != "" {
		apiVersion = options.ApiVersion
	}

	client := &Client{
		ApiUrl:     apiUrl,
		ApiVersion: apiVersion,
		HTTPClient: http.DefaultClient,
	}

	// Authenticate and get the access token
	if options.ApiKey != "" {
		err := client.Authenticate(options.ApiKey)
		if err != nil {
			return nil, fmt.Errorf("failed to authenticate during client initialization: %w", err)
		}
	}

	return client, nil
}

// NewRequest creates a new HTTP request with the access token header.
func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.ApiUrl+path, body)
	if err != nil {
		return nil, err
	}

	if c.AccessToken != "" {
		req.Header.Set("Authorization", c.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}

	return req, nil
}

// doRequest sends an HTTP request and handles the response.
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return body, nil
	}

	return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
}
