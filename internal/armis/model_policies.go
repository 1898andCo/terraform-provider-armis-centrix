// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

// Struct to match the entire API response for creating policies.
type CreatePolicyAPIResponse struct {
	Data    PolicyId `json:"data"`
	Success bool     `json:"success,omitempty"`
}

// Struct to match response for retrieving policies by ID.
type GetPolicyAPIResponse struct {
	Data    GetPolicySettings `json:"data"`
	Success bool              `json:"success,omitempty"`
}

type UpdatePolicyAPIResponse struct {
	Data    UpdatePolicySettings `json:"data"`
	Success bool                 `json:"success,omitempty"`
}

// Structs for deleting policies.
type DeletePolicyAPIResponse struct {
	Success bool `json:"success"`
}

type PolicyId struct {
	ID int `json:"id"`
}

// Consolidation represents the consolidation parameters.
type Consolidation struct {
	Amount int    `json:"amount,omitempty"`
	Unit   string `json:"unit,omitempty"`
}

// Params represents the parameters of an action.
type Params struct {
	Consolidation Consolidation `json:"consolidation,omitempty"`
	Severity      string        `json:"severity,omitempty"`
	Title         string        `json:"title,omitempty"`
	Type          string        `json:"type,omitempty"`
	Endpoint      string        `json:"endpoint,omitempty"`
	Tags          []string      `json:"tags,omitempty"`
}

// Action represents an individual action.
type Action struct {
	Params Params `json:"params,omitempty"`
	Type   string `json:"type,omitempty"`
}

// Rules represents the rules configuration for the policy.
type Rules struct {
	And []any `json:"and,omitempty"`
	Or  []any `json:"or,omitempty"`
}

// PolicySettings represents the main JSON structure.
type PolicySettings struct {
	Actions           []Action `json:"actions,omitempty"`
	Description       string   `json:"description,omitempty"`
	IsEnabled         bool     `json:"isEnabled,omitempty"`
	Labels            []string `json:"labels,omitempty"`
	MitreAttackLabels []string `json:"mitreAttackLabels,omitempty"`
	Name              string   `json:"name,omitempty"`
	RuleType          string   `json:"ruleType,omitempty"`
	Rules             Rules    `json:"rules,omitempty"`
}

// The API returns a separate response after updates that breaks the MITRE labels out.
type UpdatePolicySettings struct {
	Actions           []Action           `json:"actions,omitempty"`
	Description       string             `json:"description,omitempty"`
	IsEnabled         bool               `json:"isEnabled,omitempty"`
	Labels            []string           `json:"labels,omitempty"`
	MitreAttackLabels []MitreAttackLabel `json:"mitreAttackLabels,omitempty"`
	Name              string             `json:"name,omitempty"`
	RuleType          string             `json:"ruleType,omitempty"`
	Rules             Rules              `json:"rules,omitempty"`
}

// The API returns a separate response after updates that breaks the MITRE labels out.
type GetPolicySettings struct {
	Actions           []Action           `json:"actions,omitempty"`
	Description       string             `json:"description,omitempty"`
	IsEnabled         bool               `json:"isEnabled,omitempty"`
	Labels            []string           `json:"labels,omitempty"`
	MitreAttackLabels []MitreAttackLabel `json:"mitreAttackLabels,omitempty"`
	Name              string             `json:"name,omitempty"`
	RuleType          string             `json:"ruleType,omitempty"`
	Rules             Rules              `json:"rules,omitempty"`
}

type MitreAttackLabel struct {
	Matrix       string `json:"matrix"`
	SubTechnique string `json:"subTechnique"`
	Tactic       string `json:"tactic"`
	Technique    string `json:"technique"`
}
