// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"

	"github.com/1898andCo/terraform-provider-armis-centrix/armis"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ConvertModelToRules converts RulesModel to armis.Rules.
func ConvertModelToRules(model RulesModel) (armis.Rules, diag.Diagnostics) {
	var diags diag.Diagnostics
	rules := armis.Rules{}

	// Convert AND rules
	if !model.And.IsNull() && !model.And.IsUnknown() {
		andElements := model.And.Elements()
		rules.And = make([]any, 0, len(andElements))
		for _, elem := range andElements {
			if strVal, ok := elem.(types.String); ok && !strVal.IsNull() {
				rules.And = append(rules.And, strVal.ValueString())
			}
		}
	}

	// Convert OR rules
	if !model.Or.IsNull() && !model.Or.IsUnknown() {
		orElements := model.Or.Elements()
		rules.Or = make([]any, 0, len(orElements))
		for _, elem := range orElements {
			if strVal, ok := elem.(types.String); ok && !strVal.IsNull() {
				rules.Or = append(rules.Or, strVal.ValueString())
			}
		}
	}

	return rules, diags
}

// ConvertListToStringSlice converts a types.List to []string.
func ConvertListToStringSlice(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	elements := list.Elements()
	result := make([]string, 0, len(elements))

	for _, elem := range elements {
		if strVal, ok := elem.(types.String); ok && !strVal.IsNull() {
			result = append(result, strVal.ValueString())
		}
	}

	return result
}

// ConvertStringSliceToList converts []string to types.List.
func ConvertStringSliceToList(slice []string) types.List {
	if slice == nil {
		return types.ListNull(types.StringType)
	}

	elements := make([]attr.Value, len(slice))
	for i, s := range slice {
		elements[i] = types.StringValue(s)
	}

	listValue, _ := types.ListValue(types.StringType, elements)
	return listValue
}

// ConvertSliceToList converts []any to types.List.
func ConvertSliceToList(input []any) types.List {
	if input == nil {
		return types.ListNull(types.StringType)
	}

	elements := make([]attr.Value, 0, len(input))
	for _, item := range input {
		if str, ok := item.(string); ok {
			elements = append(elements, types.StringValue(str))
		}
	}

	listValue, _ := types.ListValue(types.StringType, elements)
	return listValue
}

// ConvertListToActions converts a types.List to []armis.Action.
func ConvertListToActions(list types.List) ([]armis.Action, diag.Diagnostics) {
	var diags diag.Diagnostics

	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	actionModels, extractDiags := extractActionModels(list)
	diags.Append(extractDiags...)
	if diags.HasError() {
		return nil, diags
	}

	actions := make([]armis.Action, 0, len(actionModels))
	for _, model := range actionModels {
		action, actionDiags := convertActionModel(model)
		diags.Append(actionDiags...)
		if diags.HasError() {
			return nil, diags
		}

		actions = append(actions, action)
	}

	return actions, diags
}

func extractActionModels(list types.List) ([]ActionModel, diag.Diagnostics) {
	var (
		models []ActionModel
		diags  diag.Diagnostics
	)

	diags.Append(list.ElementsAs(context.Background(), &models, false)...) //nolint:contextcheck
	return models, diags
}

func convertActionModel(model ActionModel) (armis.Action, diag.Diagnostics) {
	var diags diag.Diagnostics

	action := armis.Action{}
	if value, ok := stringValue(model.Type); ok {
		action.Type = value
	}

	params, hasParams, paramsDiags := paramsFromObject(model.Params)
	diags.Append(paramsDiags...)
	if diags.HasError() {
		return armis.Action{}, diags
	}

	if hasParams {
		action.Params = params
	}

	return action, diags
}

func paramsFromObject(obj types.Object) (armis.Params, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if obj.IsNull() || obj.IsUnknown() {
		return armis.Params{}, false, diags
	}

	var model ParamsModel
	diags.Append(obj.As(context.Background(), &model, basetypes.ObjectAsOptions{})...) //nolint:contextcheck
	if diags.HasError() {
		return armis.Params{}, false, diags
	}

	params := armis.Params{}
	hasParams := false

	if value, ok := stringValue(model.Severity); ok {
		params.Severity = value
		hasParams = true
	}
	if value, ok := stringValue(model.Title); ok {
		params.Title = value
		hasParams = true
	}
	if value, ok := stringValue(model.Type); ok {
		params.Type = value
		hasParams = true
	}
	if value, ok := stringValue(model.Endpoint); ok {
		params.Endpoint = value
		hasParams = true
	}
	if !model.Tags.IsNull() && !model.Tags.IsUnknown() {
		params.Tags = ConvertListToStringSlice(model.Tags)
		if len(params.Tags) > 0 {
			hasParams = true
		}
	}

	consolidation, hasConsolidation, consolidationDiags := consolidationFromObject(model.Consolidation)
	diags.Append(consolidationDiags...)
	if hasConsolidation {
		params.Consolidation = consolidation
		hasParams = true
	}

	return params, hasParams, diags
}

func consolidationFromObject(obj types.Object) (armis.Consolidation, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if obj.IsNull() || obj.IsUnknown() {
		return armis.Consolidation{}, false, diags
	}

	var model ConsolidationModel
	diags.Append(obj.As(context.Background(), &model, basetypes.ObjectAsOptions{})...) //nolint:contextcheck
	if diags.HasError() {
		return armis.Consolidation{}, false, diags
	}

	consolidation := armis.Consolidation{}
	hasValue := false

	if value, ok := intValue(model.Amount); ok {
		consolidation.Amount = value
		hasValue = true
	}
	if value, ok := stringValue(model.Unit); ok {
		consolidation.Unit = value
		hasValue = true
	}

	return consolidation, hasValue, diags
}

func stringValue(value types.String) (string, bool) {
	if value.IsNull() || value.IsUnknown() {
		return "", false
	}

	return value.ValueString(), true
}

func intValue(value types.Int64) (int, bool) {
	if value.IsNull() || value.IsUnknown() {
		return 0, false
	}

	return int(value.ValueInt64()), true
}

// ConvertActionsToList converts []armis.Action to types.List.
func ConvertActionsToList(actions []armis.Action) types.List {
	if actions == nil {
		actionObjType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type": types.StringType,
				"params": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"severity": types.StringType,
						"title":    types.StringType,
						"type":     types.StringType,
						"endpoint": types.StringType,
						"tags":     types.ListType{ElemType: types.StringType},
						"consolidation": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"amount": types.Int64Type,
								"unit":   types.StringType,
							},
						},
					},
				},
			},
		}
		return types.ListNull(actionObjType)
	}

	elements := make([]attr.Value, 0, len(actions))
	for _, action := range actions {
		// Convert params
		var paramsObj types.Object

		severityVal := types.StringNull()
		if action.Params.Severity != "" {
			severityVal = types.StringValue(action.Params.Severity)
		}

		titleVal := types.StringNull()
		if action.Params.Title != "" {
			titleVal = types.StringValue(action.Params.Title)
		}

		typeVal := types.StringNull()
		if action.Params.Type != "" {
			typeVal = types.StringValue(action.Params.Type)
		}

		endpointVal := types.StringNull()
		if action.Params.Endpoint != "" {
			endpointVal = types.StringValue(action.Params.Endpoint)
		}

		tagsVal := ConvertStringSliceToList(action.Params.Tags)

		// Convert consolidation
		consolidationObj := types.ObjectNull(map[string]attr.Type{
			"amount": types.Int64Type,
			"unit":   types.StringType,
		})
		if action.Params.Consolidation.Amount != 0 || action.Params.Consolidation.Unit != "" {
			consolidationAttrs := map[string]attr.Value{
				"amount": types.Int64Value(int64(action.Params.Consolidation.Amount)),
				"unit":   types.StringValue(action.Params.Consolidation.Unit),
			}
			consolidationObj, _ = types.ObjectValue(map[string]attr.Type{
				"amount": types.Int64Type,
				"unit":   types.StringType,
			}, consolidationAttrs)
		}

		hasParams := !severityVal.IsNull() || !titleVal.IsNull() || !typeVal.IsNull() ||
			!endpointVal.IsNull() || !tagsVal.IsNull() || !consolidationObj.IsNull()

		if hasParams {
			paramsAttrs := map[string]attr.Value{
				"severity":      severityVal,
				"title":         titleVal,
				"type":          typeVal,
				"endpoint":      endpointVal,
				"tags":          tagsVal,
				"consolidation": consolidationObj,
			}

			paramsObj, _ = types.ObjectValue(map[string]attr.Type{
				"severity": types.StringType,
				"title":    types.StringType,
				"type":     types.StringType,
				"endpoint": types.StringType,
				"tags":     types.ListType{ElemType: types.StringType},
				"consolidation": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"amount": types.Int64Type,
						"unit":   types.StringType,
					},
				},
			}, paramsAttrs)
		} else {
			paramsObj = types.ObjectNull(map[string]attr.Type{
				"severity": types.StringType,
				"title":    types.StringType,
				"type":     types.StringType,
				"endpoint": types.StringType,
				"tags":     types.ListType{ElemType: types.StringType},
				"consolidation": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"amount": types.Int64Type,
						"unit":   types.StringType,
					},
				},
			})
		}

		actionTypeVal := types.StringNull()
		if action.Type != "" {
			actionTypeVal = types.StringValue(action.Type)
		}

		actionAttrs := map[string]attr.Value{
			"type":   actionTypeVal,
			"params": paramsObj,
		}

		actionObj, _ := types.ObjectValue(map[string]attr.Type{
			"type": types.StringType,
			"params": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"severity": types.StringType,
					"title":    types.StringType,
					"type":     types.StringType,
					"endpoint": types.StringType,
					"tags":     types.ListType{ElemType: types.StringType},
					"consolidation": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"amount": types.Int64Type,
							"unit":   types.StringType,
						},
					},
				},
			},
		}, actionAttrs)

		elements = append(elements, actionObj)
	}

	listValue, _ := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type": types.StringType,
			"params": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"severity": types.StringType,
					"title":    types.StringType,
					"type":     types.StringType,
					"endpoint": types.StringType,
					"tags":     types.ListType{ElemType: types.StringType},
					"consolidation": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"amount": types.Int64Type,
							"unit":   types.StringType,
						},
					},
				},
			},
		},
	}, elements)

	return listValue
}

// ResponseToPolicyFromGet converts armis.GetPolicySettings to PolicyResourceModel.
func ResponseToPolicyFromGet(ctx context.Context, policy armis.GetPolicySettings) *PolicyResourceModel {
	tflog.Debug(ctx, "Processing policy", map[string]any{
		"policy_name":    policy.Name,
		"policy_type":    policy.RuleType,
		"policy_enabled": policy.IsEnabled,
	})

	result := &PolicyResourceModel{
		Name:        types.StringValue(policy.Name),
		Description: types.StringValue(policy.Description),
		IsEnabled:   types.BoolValue(policy.IsEnabled),
		Labels:      ConvertStringSliceToList(policy.Labels),
		RuleType:    types.StringValue(policy.RuleType),
		Actions:     ConvertActionsToList(policy.Actions),
		Rules: &RulesModel{
			And: ConvertSliceToList(policy.Rules.And),
			Or:  ConvertSliceToList(policy.Rules.Or),
		},
	}

	return result
}

// ResponseToPolicyFromUpdate converts armis.UpdatePolicySettings to PolicyResourceModel.
func ResponseToPolicyFromUpdate(ctx context.Context, policy armis.UpdatePolicySettings) *PolicyResourceModel {
	tflog.Debug(ctx, "Processing updated policy", map[string]any{
		"policy_name":    policy.Name,
		"policy_type":    policy.RuleType,
		"policy_enabled": policy.IsEnabled,
	})

	result := &PolicyResourceModel{
		Name:        types.StringValue(policy.Name),
		Description: types.StringValue(policy.Description),
		IsEnabled:   types.BoolValue(policy.IsEnabled),
		Labels:      ConvertStringSliceToList(policy.Labels),
		RuleType:    types.StringValue(policy.RuleType),
		Actions:     ConvertActionsToList(policy.Actions),
		Rules: &RulesModel{
			And: ConvertSliceToList(policy.Rules.And),
			Or:  ConvertSliceToList(policy.Rules.Or),
		},
	}

	return result
}

// BuildPolicySettings converts a PolicyResourceModel to armis.PolicySettings.
func BuildPolicySettings(model *PolicyResourceModel) (armis.PolicySettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	policy := armis.PolicySettings{
		Name:              model.Name.ValueString(),
		Description:       model.Description.ValueString(),
		IsEnabled:         model.IsEnabled.ValueBool(),
		RuleType:          model.RuleType.ValueString(),
		MitreAttackLabels: ConvertListToStringSlice(model.MitreAttackLabels),
		Labels:            ConvertListToStringSlice(model.Labels),
	}

	// Convert actions
	actions, actionDiags := ConvertListToActions(model.Actions)
	diags.Append(actionDiags...)
	policy.Actions = actions

	// Convert rules
	rules, ruleDiags := ConvertModelToRules(*model.Rules)
	diags.Append(ruleDiags...)
	policy.Rules = rules

	return policy, diags
}
