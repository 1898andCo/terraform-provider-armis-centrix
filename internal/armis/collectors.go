// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *Client) GetCollectorByID(collectorId string) (*CollectorSettings, error) {
	if collectorId == "" {
		return nil, fmt.Errorf("collector ID cannot be empty")
	}

	// URL encide the collector ID
	encodedCollectorId := url.QueryEscape(collectorId)

	// Create a new request
	req, err := c.newRequest("GET", fmt.Sprintf("/api/%s/collectors/%s/", c.ApiVersion, encodedCollectorId), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetCollector: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collector: %w", err)
	}

	// Parse the response
	var response GetCollectorApiResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse collector response: %w", err)
	}

	return &response.Data, nil
}

func (c *Client) GetCollectors() ([]CollectorSettings, error) {
	// Create a new request
	req, err := c.newRequest("GET", fmt.Sprintf("/api/%s/collectors/", c.ApiVersion), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetCollectors: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collectors: %w", err)
	}

	// Parse the response
	var response CollectorApiResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse collectors response: %w", err)
	}

	return response.Data.Collectors, nil
}

func (c *Client) CreateCollector(collector CreateCollectorSettings) (*NewCollectorSettings, error) {
	if collector.Name == "" {
		return nil, fmt.Errorf("collector name is required")
	}

	if collector.DeploymentType == "" {
		return nil, fmt.Errorf("collector deployment type is required")
	}

	collectorData, err := json.Marshal(collector)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal collector data: %w", err)
	}

	// Create a new request
	req, err := c.newRequest("POST", fmt.Sprintf("/api/%s/collectors/", c.ApiVersion), bytes.NewBuffer(collectorData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for CreateCollector: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create collector: %w", err)
	}

	// Parse the response
	var response CreateCollectorApiResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse create collector response: %w", err)
	}

	return &response.Data, nil
}

func (c *Client) UpdateCollector(collectorId string, collector UpdateCollectorSettings) (*CollectorSettings, error) {
	if collectorId == "" {
		return nil, fmt.Errorf("collector ID is required")
	}

	if collector.Name == "" {
		return nil, fmt.Errorf("collector name is required")
	}

	if collector.DeploymentType == "" {
		return nil, fmt.Errorf("collector deployment type is required")
	}

	collectorData, err := json.Marshal(collector)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal collector data: %w", err)
	}

	// URL encode the collector ID
	encodedCollectorId := url.QueryEscape(collectorId)

	// Create a new request
	req, err := c.newRequest("PATCH", fmt.Sprintf("/api/%s/collectors/%s/", c.ApiVersion, encodedCollectorId), bytes.NewBuffer(collectorData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for UpdateCollector: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update collector: %w", err)
	}

	// Parse the response
	var response UpdateCollectorApiResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse update collector response: %w", err)
	}

	return &response.Data, nil
}

func (c *Client) DeleteCollector(collectorId string) (bool, error) {
	if collectorId == "" {
		return false, fmt.Errorf("collector ID cannot be empty")
	}

	// URL encode the collector ID
	encodedCollectorId := url.QueryEscape(collectorId)

	// Create a new request
	req, err := c.newRequest("DELETE", fmt.Sprintf("/api/%s/collectors/%s/", c.ApiVersion, encodedCollectorId), nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request for DeleteCollector: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return false, fmt.Errorf("failed to delete collector: %w", err)
	}

	// Parse the response
	var response DeleteCollectorApiResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return false, fmt.Errorf("failed to parse collector response: %w", err)
	}

	return response.Success, nil
}
