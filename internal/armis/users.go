// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

// GetUsers fetches all users from Armis.
func (c *Client) GetUsers() ([]UserSettings, error) {
	// Create a new request
	req, err := c.newRequest("GET", fmt.Sprintf("/api/%s/users/", c.ApiVersion), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetUsers: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	// Parse the response
	var response UserApiResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse users response: %w", err)
	}

	return response.Data.Users, nil
}

// GetUser fetches a user from Armis using the user's ID or email.
func (c *Client) GetUser(userId string) (*UserSettings, error) {
	if userId == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	// URL encode the user ID
	encodedUserId := url.QueryEscape(userId)

	// Create a new request
	req, err := c.newRequest("GET", fmt.Sprintf("/api/%s/users/%s/", c.ApiVersion, encodedUserId), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Get User: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Parse the response
	var response GetUserApiResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API error: response indicates failure: %+v", response)
	}

	return &response.Data, nil
}

// CreateUser creates a new user in Armis.
func (c *Client) CreateUser(user UserSettings) (*UserSettings, error) {
	if user.Name == "" {
		return nil, fmt.Errorf("user name cannot be empty")
	}

	if user.Email == "" {
		return nil, fmt.Errorf("user email cannot be empty")
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	req, err := c.newRequest("POST", fmt.Sprintf("/api/%s/users/", c.ApiVersion), bytes.NewReader(userData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for CreateUser: %w", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user %q: %w", user.Name, err)
	}

	var apiResponse CreateUserApiResponse
	if err := json.Unmarshal(res, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("API error: response indicates failure: %+v", apiResponse)
	}

	// Return the parsed user settings directly
	return &apiResponse.Data, nil
}

// UpdateUser updates a user in Armis.
func (c *Client) UpdateUser(user UserSettings, userId string) (*UserSettings, error) {
	if user.Name == "" {
		return nil, fmt.Errorf("user name cannot be empty")
	}

	if user.Email == "" {
		return nil, fmt.Errorf("user email cannot be empty")
	}

	if userId == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	// URL encode the user ID
	encodedUserId := url.QueryEscape(userId)

	req, err := c.newRequest("PATCH", fmt.Sprintf("/api/%s/users/%s/", c.ApiVersion, encodedUserId), bytes.NewReader(userData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for UpdateUser: %w", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user %q: %w", user.Name, err)
	}

	var apiResponse UpdateUserApiResponse
	if err := json.Unmarshal(res, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("API error: response indicates failure: %+v", apiResponse)
	}

	// Return the parsed user settings directly
	return &apiResponse.Data, nil
}

// DeleteUser deletes a user from Armis.
func (c *Client) DeleteUser(userId string) (bool, error) {
	if userId == "" {
		return false, fmt.Errorf("user ID cannot be empty")
	}

	// URL encode the user ID
	encodedUserId := url.QueryEscape(userId)

	// Create a new request
	req, err := c.newRequest("DELETE", fmt.Sprintf("/api/%s/users/%s/", c.ApiVersion, encodedUserId), nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request for DeleteUser: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return false, fmt.Errorf("failed to delete user: %w", err)
	}

	// Parse the response
	var response DeleteUserApiResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return false, fmt.Errorf("failed to parse user response: %w", err)
	}

	return response.Success, nil
}
