// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

// SearchAPIResponse represents the response returned by the search endpoint.
type SearchAPIResponse struct {
	Data    SearchData `json:"data"`
	Success bool       `json:"success"`
}

// SearchData contains pagination metadata and the list of search results.
type SearchData struct {
	Count   int            `json:"count"`
	Next    *int           `json:"next"`
	Prev    *int           `json:"prev"`
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}

// SearchResult represents a single search hit returned by Armis.
type SearchResult struct {
	ActivityUUIDs        []string         `json:"activityUUIDs,omitempty"`
	AffectedDevicesCount int              `json:"affectedDevicesCount,omitempty"`
	AlertID              int              `json:"alertId,omitempty"`
	Classification       string           `json:"classification,omitempty"`
	ConnectionIDs        []string         `json:"connectionIds,omitempty"`
	Description          string           `json:"description,omitempty"`
	DestinationEndpoints []SearchEndpoint `json:"destinationEndpoints,omitempty"`
	DeviceIDs            []int            `json:"deviceIds,omitempty"`
	LastAlertUpdateTime  string           `json:"lastAlertUpdateTime,omitempty"`
	MitreAttackLabels    []string         `json:"mitreAttackLabels,omitempty"`
	PolicyID             string           `json:"policyId,omitempty"`
	PolicyLabels         []string         `json:"policyLabels,omitempty"`
	PolicyTitle          string           `json:"policyTitle,omitempty"`
	Severity             string           `json:"severity,omitempty"`
	SourceEndpoints      []SearchEndpoint `json:"sourceEndpoints,omitempty"`
	Status               string           `json:"status,omitempty"`
	StatusChangeTime     string           `json:"statusChangeTime,omitempty"`
	Time                 string           `json:"time,omitempty"`
	Title                string           `json:"title,omitempty"`
	Type                 string           `json:"type,omitempty"`
}

// SearchEndpoint represents an endpoint referenced in a search result.
type SearchEndpoint struct {
	ID   string `json:"id,omitempty"`
	IP   string `json:"ip,omitempty"`
	Name string `json:"name,omitempty"`
}
