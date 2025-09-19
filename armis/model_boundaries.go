// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

// BoundaryAPIResponse struct to match the entire API response.
type BoundaryAPIResponse struct {
	Data    Site `json:"data"`
	Success bool `json:"success,omitempty"`
}

// Boundaries represents a collection of boundaries.
type Boundaries struct {
	Count      int                `json:"count,omitempty"`
	Next       *int               `json:"next"` // Handle null as a pointer
	Prev       int                `json:"prev"`
	Boundaries []BoundarySettings `json:"boundaries"`
}

// BoundarySettings struct for individual boundary settings.
type BoundarySettings struct {
	ID            string  `json:"id,omitempty"`
	AffectedSites string  `json:"affectedSites,omitempty"`
	Name          string  `json:"name,omitempty"`
	RuleAQL       RuleAQL `json:"ruleAql,omitempty"`
}

// RuleAQL struct for individual rule AQL settings.
type RuleAQL struct {
	And []string `json:"and,omitempty"`
	Or  []string `json:"or,omitempty"`
}
