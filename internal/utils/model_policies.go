// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package utils

import "github.com/hashicorp/terraform-plugin-framework/types"

// PolicyResourceModel maps the resource schema data.
type PolicyResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	IsEnabled         types.Bool   `tfsdk:"enabled"`
	Labels            types.List   `tfsdk:"labels"`
	MitreAttackLabels types.List   `tfsdk:"mitre_attack_labels"`
	RuleType          types.String `tfsdk:"rule_type"`
	Actions           types.List   `tfsdk:"actions"`
	Rules             *RulesModel  `tfsdk:"rules"`
}

// ActionModel maps the action schema data.
type ActionModel struct {
	Type   types.String `tfsdk:"type"`
	Params types.Object `tfsdk:"params"`
}

// ParamsModel maps the params schema data.
type ParamsModel struct {
	Severity      types.String `tfsdk:"severity"`
	Title         types.String `tfsdk:"title"`
	Type          types.String `tfsdk:"type"`
	Endpoint      types.String `tfsdk:"endpoint"`
	Tags          types.List   `tfsdk:"tags"`
	Consolidation types.Object `tfsdk:"consolidation"`
}

// ConsolidationModel maps the consolidation schema data.
type ConsolidationModel struct {
	Amount types.Int64  `tfsdk:"amount"`
	Unit   types.String `tfsdk:"unit"`
}

// RulesModel maps the rules schema data.
type RulesModel struct {
	And types.List `tfsdk:"and"`
	Or  types.List `tfsdk:"or"`
}

type PoliciesDataSourceModel struct {
	PolicyID      types.String                  `tfsdk:"policy_id"`
	MatchPrefix   types.String                  `tfsdk:"match_prefix"`
	ExcludePrefix types.String                  `tfsdk:"exclude_prefix"`
	Policies      []PolicyDataSourcePolicyModel `tfsdk:"policies"`
}

type PolicyDataSourcePolicyModel struct {
	ID                types.String                      `tfsdk:"id"`
	Name              types.String                      `tfsdk:"name"`
	Description       types.String                      `tfsdk:"description"`
	IsEnabled         types.Bool                        `tfsdk:"enabled"`
	RuleType          types.String                      `tfsdk:"rule_type"`
	Labels            []types.String                    `tfsdk:"labels"`
	MitreAttackLabels []PolicyDataSourceMitreLabelModel `tfsdk:"mitre_attack_labels"`
	Actions           []PolicyDataSourceActionModel     `tfsdk:"actions"`
	Rules             PolicyDataSourceRulesModel        `tfsdk:"rules"`
}

type PolicyDataSourceMitreLabelModel struct {
	Matrix       types.String `tfsdk:"matrix"`
	SubTechnique types.String `tfsdk:"sub_technique"`
	Tactic       types.String `tfsdk:"tactic"`
	Technique    types.String `tfsdk:"technique"`
}

type PolicyDataSourceActionModel struct {
	Type   types.String                `tfsdk:"type"`
	Params PolicyDataSourceParamsModel `tfsdk:"params"`
}

type PolicyDataSourceParamsModel struct {
	Severity      types.String                       `tfsdk:"severity"`
	Title         types.String                       `tfsdk:"title"`
	Type          types.String                       `tfsdk:"type"`
	Endpoint      types.String                       `tfsdk:"endpoint"`
	Tags          []types.String                     `tfsdk:"tags"`
	Consolidation PolicyDataSourceConsolidationModel `tfsdk:"consolidation"`
}

type PolicyDataSourceConsolidationModel struct {
	Amount types.Int64  `tfsdk:"amount"`
	Unit   types.String `tfsdk:"unit"`
}

type PolicyDataSourceRulesModel struct {
	And []types.String `tfsdk:"and"`
	Or  []types.String `tfsdk:"or"`
}
