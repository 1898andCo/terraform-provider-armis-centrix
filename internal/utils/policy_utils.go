// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"

	"github.com/1898andCo/terraform-provider-armis-centrix/internal/armis"
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

	var actionModels []ActionModel
	diags.Append(list.ElementsAs(context.Background(), &actionModels, false)...)
	if diags.HasError() {
		return nil, diags
	}

	actions := make([]armis.Action, 0, len(actionModels))
	for _, am := range actionModels {
		action := armis.Action{
			Type: am.Type.ValueString(),
		}

		// Convert params if present
		if !am.Params.IsNull() && !am.Params.IsUnknown() {
			var paramsModel ParamsModel
			diags.Append(am.Params.As(context.Background(), &paramsModel, basetypes.ObjectAsOptions{})...)

			params := armis.Params{
				Severity: paramsModel.Severity.ValueString(),
				Title:    paramsModel.Title.ValueString(),
				Type:     paramsModel.Type.ValueString(),
				Endpoint: paramsModel.Endpoint.ValueString(),
				Tags:     ConvertListToStringSlice(paramsModel.Tags),
			}

			// Convert consolidation if present
			if !paramsModel.Consolidation.IsNull() && !paramsModel.Consolidation.IsUnknown() {
				var consolidationModel ConsolidationModel
				diags.Append(paramsModel.Consolidation.As(context.Background(), &consolidationModel, basetypes.ObjectAsOptions{})...)

				params.Consolidation = armis.Consolidation{
					Amount: int(consolidationModel.Amount.ValueInt64()),
					Unit:   consolidationModel.Unit.ValueString(),
				}
			}

			action.Params = params
		}

		actions = append(actions, action)
	}

	return actions, diags
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

		// Check if Params has any non-zero values
		hasParams := action.Params.Severity != "" || action.Params.Title != "" ||
			action.Params.Type != "" || len(action.Params.Tags) > 0 ||
			action.Params.Consolidation.Amount != 0 ||
			action.Params.Consolidation.Unit != ""

		// Always include endpoint in hasParams check to preserve its value
		// even if it's an empty string
		hasParams = hasParams || action.Params.Endpoint != ""

		if hasParams {
			// Convert consolidation
			var consolidationObj types.Object
			if action.Params.Consolidation.Amount != 0 || action.Params.Consolidation.Unit != "" {
				consolidationAttrs := map[string]attr.Value{
					"amount": types.Int64Value(int64(action.Params.Consolidation.Amount)),
					"unit":   types.StringValue(action.Params.Consolidation.Unit),
				}
				consolidationObj, _ = types.ObjectValue(map[string]attr.Type{
					"amount": types.Int64Type,
					"unit":   types.StringType,
				}, consolidationAttrs)
			} else {
				consolidationObj = types.ObjectNull(map[string]attr.Type{
					"amount": types.Int64Type,
					"unit":   types.StringType,
				})
			}

			paramsAttrs := map[string]attr.Value{
				"severity":      types.StringValue(action.Params.Severity),
				"title":         types.StringValue(action.Params.Title),
				"type":          types.StringValue(action.Params.Type),
				"endpoint":      types.StringValue(action.Params.Endpoint),
				"tags":          ConvertStringSliceToList(action.Params.Tags),
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

		actionAttrs := map[string]attr.Value{
			"type":   types.StringValue(action.Type),
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
		Rules: RulesModel{
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
		Rules: RulesModel{
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
	rules, ruleDiags := ConvertModelToRules(model.Rules)
	diags.Append(ruleDiags...)
	policy.Rules = rules

	return policy, diags
}
