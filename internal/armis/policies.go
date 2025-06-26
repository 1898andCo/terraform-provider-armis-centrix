// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// CreatePolicy creates a new policy in Armis.
func (c *Client) CreatePolicy(ctx context.Context, policy PolicySettings) (*PolicyId, error) {
	if policy.Name == "" {
		return nil, fmt.Errorf("%w", ErrPolicyName)
	}

	// Ensure policy rules exist
	if len(policy.Rules.And) == 0 && len(policy.Rules.Or) == 0 {
		return nil, fmt.Errorf("%w", ErrPolicyRules)
	}

	// Policy description must be less than 500 characters
	if len(policy.Description) > 500 {
		return nil, fmt.Errorf("%w", ErrPolicyDescription)
	}

	// Rule type must be ACTIVITY, IP CONNECTION, DEVICE, or VULNERABILITY
	if policy.RuleType != "ACTIVITY" && policy.RuleType != "IP_CONNECTION" && policy.RuleType != "DEVICE" && policy.RuleType != "VULNERABILITY" {
		return nil, fmt.Errorf("%w", ErrPolicyRuleType)
	}

	policyData, err := json.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal policy data: %w", err)
	}

	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/api/%s/policies/", c.apiVersion), bytes.NewReader(policyData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for CreatePolicy: %w", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create policy %q: %w", policy.Name, err)
	}

	var apiResponse CreatePolicyApiResponse
	if err := json.Unmarshal(res, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse policy response: %w", err)
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("%w:%v", ErrHTTPResponse, apiResponse)
	}

	// Return the parsed policy settings directly
	return &apiResponse.Data, nil
}

// GetPolicy fetches a policy from Armis using the policy's ID.
func (c *Client) GetPolicy(ctx context.Context, policyID string) (*GetPolicySettings, error) {
	if policyID == "" {
		return nil, fmt.Errorf("%w", ErrPolicyID)
	}

	// URL encode the policy ID
	encodedPolicyID := url.QueryEscape(policyID)

	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/policies/%s/", c.apiVersion, encodedPolicyID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Get Policy: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch policy: %w", err)
	}

	// Parse the response
	var response GetPolicyApiResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse policy response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("%w:%v", ErrHTTPResponse, response)
	}

	return &response.Data, nil
}

// UpdatePolicy updates a policy in Armis.
func (c *Client) UpdatePolicy(ctx context.Context, policy PolicySettings, policyID string) (*UpdatePolicySettings, error) {
	if policy.Name == "" {
		return nil, fmt.Errorf("%w", ErrPolicyName)
	}

	if policyID == "" {
		return nil, fmt.Errorf("%w", ErrPolicyID)
	}

	policyData, err := json.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal policy data: %w", err)
	}

	// URL encode the policy ID
	encodedPolicyID := url.QueryEscape(policyID)

	req, err := c.newRequest(ctx, "PATCH", fmt.Sprintf("/api/%s/policies/%s/", c.apiVersion, encodedPolicyID), bytes.NewReader(policyData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for UpdatePolicy: %w", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update policy %q: %w", policy.Name, err)
	}

	var apiResponse UpdatePolicyApiResponse
	if err := json.Unmarshal(res, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse policy response: %w", err)
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("%w:%v", ErrHTTPResponse, apiResponse)
	}

	// Return the parsed policy settings directly
	return &apiResponse.Data, nil
}

// DeletePolicy deletes a policy from Armis.
func (c *Client) DeletePolicy(ctx context.Context, policyID string) (bool, error) {
	if policyID == "" {
		return false, fmt.Errorf("%w", ErrPolicyID)
	}

	// URL encode the policy ID
	encodedPolicyID := url.QueryEscape(policyID)

	// Create a new request
	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("/api/%s/policies/%s/", c.apiVersion, encodedPolicyID), nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request for DeletePolicy: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return false, fmt.Errorf("failed to delete policy: %w", err)
	}

	// Parse the response
	var response DeletePolicyApiResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return false, fmt.Errorf("failed to parse policy response: %w", err)
	}

	return response.Success, nil
}

// Ensure the policy rules exist.
func (r Rules) Len() int {
	return len(r.And)
}
