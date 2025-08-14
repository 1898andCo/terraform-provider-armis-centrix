// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	armis "github.com/1898andCo/terraform-provider-armis-centrix/internal/armis"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &policyResource{}
	_ resource.ResourceWithConfigure = &policyResource{}
)

type policyResource struct {
	client *armis.Client
}

func PolicyResource() resource.Resource {
	return &policyResource{}
}

// Configure adds the provider configured client to the resource.
func (r *policyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*armis.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *armis.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *policyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

// Schema defines the schema for the policy resource.
func (r *policyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
		Provides an Armis policy

		The resource provisions a policy with the ability to define rules, parameters, and settings.
		`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "The ID of the policy.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The full name of the policy.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the policy.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the policy is enabled.",
			},
			"labels": schema.ListAttribute{
				Optional:    true,
				Description: "A list of labels to apply to the policy.",
				ElementType: types.StringType,
			},
			"mitre_attack_labels": schema.ListAttribute{
				Optional:    true,
				Description: "A list of MITRE ATT&CK labels to apply to the policy.",
				ElementType: types.StringType,
			},
			"rule_type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of rule to apply to the policy.",
				Validators: []validator.String{
					// Must be uppercase
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Z_]+$`),
						"must contain only uppercase letters and underscores",
					),
					// Must be Activity, IP Connection, Device or Vulnerability
					stringvalidator.OneOf("ACTIVITY", "IP_CONNECTION", "DEVICE", "VULNERABILITY"),
				},
			},
			"actions": schema.ListNestedAttribute{
				Optional:    true,
				Description: "A list of actions to apply to the policy.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Optional:    true,
							Description: "The type of action to apply to the policy.",
						},
						"params": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "The parameters for the action.",
							Attributes: map[string]schema.Attribute{
								"severity": schema.StringAttribute{
									Optional:    true,
									Description: "The severity of the action.",
								},
								"title": schema.StringAttribute{
									Optional:    true,
									Description: "The title of the action.",
								},
								"type": schema.StringAttribute{
									Optional:    true,
									Description: "The type of the action.",
									Validators: []validator.String{
										stringvalidator.OneOf("Network Performance", "Security - Other", "Security - Risk", "Security - Threat"),
									},
								},
								"endpoint": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "Endpoints to apply this action to.",
									Default:     stringdefault.StaticString("ALL"),
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"tags": schema.ListAttribute{
									Optional:    true,
									Description: "Tags to apply to the action.",
									ElementType: types.StringType,
								},
								"consolidation": schema.SingleNestedAttribute{
									Optional:    true,
									Description: "The consolidation settings for the action.",
									Attributes: map[string]schema.Attribute{
										"amount": schema.Int64Attribute{
											Optional:    true,
											Description: "The amount of time to consolidate the action.",
										},
										"unit": schema.StringAttribute{
											Optional:    true,
											Description: "The unit of time to consolidate the action.",
										},
									},
								},
							},
						},
					},
				},
			},
			"rules": schema.SingleNestedAttribute{
				Required:    true,
				Description: "The rules to apply to the policy.",
				Attributes: map[string]schema.Attribute{
					"and": schema.ListAttribute{
						Optional:    true,
						Description: "A list of AND rules to apply to the policy.",
						ElementType: types.StringType,
					},
					"or": schema.ListAttribute{
						Optional:    true,
						Description: "A list of OR rules to apply to the policy.",
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

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
	Rules             RulesModel   `tfsdk:"rules"`
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

// BuildPolicySettings converts a PolicyResourceModel to armis.PolicySettings.
func BuildPolicySettings(model *PolicyResourceModel) (armis.PolicySettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	policy := armis.PolicySettings{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsEnabled:   model.IsEnabled.ValueBool(),
		RuleType:    model.RuleType.ValueString(),
		Labels:      ConvertListToStringSlice(model.Labels),
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

// responseToPolicyFromGet converts armis.GetPolicySettings to PolicyResourceModel.
func responseToPolicyFromGet(ctx context.Context, policy armis.GetPolicySettings) *PolicyResourceModel {
	tflog.Debug(ctx, "Processing policy", map[string]any{
		"policy_name":    policy.Name,
		"policy_type":    policy.RuleType,
		"policy_enabled": policy.IsEnabled,
	})

	// TODO: Handle MitreAttackLabels properly
	emptyMitreLabels, _ := types.ListValue(types.StringType, []attr.Value{})

	result := &PolicyResourceModel{
		Name:              types.StringValue(policy.Name),
		Description:       types.StringValue(policy.Description),
		IsEnabled:         types.BoolValue(policy.IsEnabled),
		Labels:            ConvertStringSliceToList(policy.Labels),
		MitreAttackLabels: emptyMitreLabels,
		RuleType:          types.StringValue(policy.RuleType),
		Actions:           ConvertActionsToList(policy.Actions),
		Rules: RulesModel{
			And: ConvertSliceToList(policy.Rules.And),
			Or:  ConvertSliceToList(policy.Rules.Or),
		},
	}

	return result
}

// responseToPolicyFromUpdate converts armis.UpdatePolicySettings to PolicyResourceModel.
func responseToPolicyFromUpdate(ctx context.Context, policy armis.UpdatePolicySettings) *PolicyResourceModel {
	tflog.Debug(ctx, "Processing updated policy", map[string]any{
		"policy_name":    policy.Name,
		"policy_type":    policy.RuleType,
		"policy_enabled": policy.IsEnabled,
	})

	// TODO: Handle MitreAttackLabels properly
	emptyMitreLabels, _ := types.ListValue(types.StringType, []attr.Value{})

	result := &PolicyResourceModel{
		Name:              types.StringValue(policy.Name),
		Description:       types.StringValue(policy.Description),
		IsEnabled:         types.BoolValue(policy.IsEnabled),
		Labels:            ConvertStringSliceToList(policy.Labels),
		MitreAttackLabels: emptyMitreLabels,
		RuleType:          types.StringValue(policy.RuleType),
		Actions:           ConvertActionsToList(policy.Actions),
		Rules: RulesModel{
			And: ConvertSliceToList(policy.Rules.And),
			Or:  ConvertSliceToList(policy.Rules.Or),
		},
	}

	return result
}

// Create decodes the plan into a model, converts it to an Armis
// PolicySettings payload, invokes r.client.CreatePolicy, stores the returned
// policy ID in state, and writes the updated state back—aborting early whenever
// diagnostics report an error.
func (r *policyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, diags := BuildPolicySettings(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createResp, err := r.client.CreatePolicy(ctx, policy)
	if err != nil {
		resp.Diagnostics.AddError("Error creating policy", fmt.Sprintf("API error: %v", err))
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(createResp.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *policyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getResp, err := r.client.GetPolicy(ctx, state.ID.ValueString())
	if err != nil {
		// Handle 404 Not Found by removing resource from state
		if strings.Contains(err.Error(), "status: 404") {
			tflog.Warn(ctx, "Policy not found, removing from state", map[string]any{
				"policy_id": state.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading policy", err.Error())
		return
	}

	// Update state with the retrieved policy data
	result := responseToPolicyFromGet(ctx, getResp)

	// Preserve the MitreAttackLabels from state
	result.MitreAttackLabels = state.MitreAttackLabels

	// Preserve the ID from state
	result.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

// Update loads plan and state, maps the plan to an Armis PolicySettings
// payload, calls r.client.UpdatePolicy with the existing ID, and writes the
// (unchanged-ID) state back—bailing out on any diagnostics or API error.
func (r *policyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, diags := BuildPolicySettings(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateResp, err := r.client.UpdatePolicy(ctx, policy, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating policy", fmt.Sprintf("API error: %v", err))
		return
	}

	// Update the plan with the response data
	result := responseToPolicyFromUpdate(ctx, updateResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the ID
	result.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

// Delete removes a policy of the provided ID.
func (r *policyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PolicyResourceModel
	if diags := req.State.Get(ctx, &state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	deleteResp, err := r.client.DeletePolicy(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting policy", err.Error())
		return
	}

	if !deleteResp {
		resp.Diagnostics.AddError("Error deleting policy", "Delete operation was not successful")
	}
}
