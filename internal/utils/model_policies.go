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
	Rules             types.Object `tfsdk:"rules"`
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
