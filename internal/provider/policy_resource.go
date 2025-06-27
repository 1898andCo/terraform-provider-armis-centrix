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
						"must contain only uppercase letters",
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
									Description: "Endpoints to apply this action to.",
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

// policyResourceModel maps the resource schema data.
type policyResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	IsEnabled         types.Bool   `tfsdk:"enabled"`
	Labels            types.List   `tfsdk:"labels"`
	MitreAttackLabels types.List   `tfsdk:"mitre_attack_labels"`
	RuleType          types.String `tfsdk:"rule_type"`
	Actions           types.List   `tfsdk:"actions"`
	Rules             rulesModel   `tfsdk:"rules"`
}

// actionModel maps the action schema data.
type actionModel struct {
	Type   types.String `tfsdk:"type"`
	Params types.Object `tfsdk:"params"`
}

// paramsModel maps the params schema data.
type paramsModel struct {
	Severity      types.String `tfsdk:"severity"`
	Title         types.String `tfsdk:"title"`
	Type          types.String `tfsdk:"type"`
	Endpoint      types.String `tfsdk:"endpoint"`
	Tags          types.List   `tfsdk:"tags"`
	Consolidation types.Object `tfsdk:"consolidation"`
}

// consolidationModel maps the consolidation schema data.
type consolidationModel struct {
	Amount types.Int64  `tfsdk:"amount"`
	Unit   types.String `tfsdk:"unit"`
}

// rules maps the rules schema data.
type rulesModel struct {
	And types.List `tfsdk:"and"`
	Or  types.List `tfsdk:"or"`
}

// extractPolicyFromPlan transforms a typed Terraform plan (policyResourceModel) into
// an armis.PolicySettings request for the Armis API. It pulls each collection field
// out of the plan with ElementsAs, accumulating any element-conversion diagnostics
// and aborting early if errors occur. When successful it maps scalars, labels,
// MITRE ATT&CK tags, actions, and AND/OR rule slices (normalising RuleType to
// upper-case and converting rule slices to []interface{}) into the API struct and
// returns it alongside any collected diagnostics.
func extractPolicyFromPlan(ctx context.Context, plan *policyResourceModel) (armis.PolicySettings, diag.Diagnostics) {
	var diags diag.Diagnostics
	var mitreAttackLabels, labels []string
	var actions []actionModel
	var andRules, orRules []string

	if d := plan.MitreAttackLabels.ElementsAs(ctx, &mitreAttackLabels, false); d.HasError() {
		diags.Append(d...)
	}
	if d := plan.Labels.ElementsAs(ctx, &labels, false); d.HasError() {
		diags.Append(d...)
	}
	if d := plan.Actions.ElementsAs(ctx, &actions, false); d.HasError() {
		diags.Append(d...)
	}
	if d := plan.Rules.And.ElementsAs(ctx, &andRules, false); d.HasError() {
		diags.Append(d...)
	}
	if d := plan.Rules.Or.ElementsAs(ctx, &orRules, false); d.HasError() {
		diags.Append(d...)
	}

	if diags.HasError() {
		return armis.PolicySettings{}, diags
	}

	return armis.PolicySettings{
		Name:              plan.Name.ValueString(),
		Description:       plan.Description.ValueString(),
		IsEnabled:         plan.IsEnabled.ValueBool(),
		RuleType:          strings.ToUpper(plan.RuleType.ValueString()),
		Labels:            labels,
		MitreAttackLabels: mitreAttackLabels,
		Actions:           convertActionsToAPI(actions),
		Rules: armis.Rules{
			And: convertStringSliceToInterface(andRules),
			Or:  convertStringSliceToInterface(orRules),
		},
	}, diags
}

// Create decodes the plan into a model, converts it to an Armis
// PolicySettings payload, invokes r.client.CreatePolicy, stores the returned
// policy ID in state, and writes the updated state back—aborting early whenever
// diagnostics report an error.
func (r *policyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan policyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, diags := extractPolicyFromPlan(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newPolicy, err := r.client.CreatePolicy(ctx, policy)
	if err != nil {
		resp.Diagnostics.AddError("Error creating policy", fmt.Sprintf("API error: %v", err))
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(newPolicy.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *policyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.client.GetPolicy(ctx, state.ID.ValueString())
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

	if policy == nil {
		resp.State.RemoveResource(ctx)
		tflog.Warn(ctx, "Policy is nil, removing from state", map[string]any{
			"policy_id": state.ID.ValueString(),
		})
		return
	}

	// Update state with the retrieved policy data
	result, diags := responseToPolicy(ctx, policy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

// Update loads plan and state, maps the plan to an Armis PolicySettings
// payload, calls r.client.UpdatePolicy with the existing ID, and writes the
// (unchanged-ID) state back—bailing out on any diagnostics or API error.
func (r *policyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state policyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, diags := extractPolicyFromPlan(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdatePolicy(ctx, policy, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating policy", fmt.Sprintf("API error: %v", err))
		return
	}

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Deletes a policy of the provided ID.
func (r *policyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyResourceModel
	if diags := req.State.Get(ctx, &state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if success, err := r.client.DeletePolicy(ctx, state.ID.ValueString()); err != nil || !success {
		resp.Diagnostics.AddError("Error deleting policy", err.Error())
	}
}

// convertActionsToAPI transforms internal actionModel entries into Armis
// Action structs: it copies the type, converts non-null/non-unknown params
// with convertParamsToAPI (skipping any that fail extraction), logs the result
// for debugging, and returns the assembled slice.
func convertActionsToAPI(actions []actionModel) []armis.Action {
	var apiActions []armis.Action
	for _, a := range actions {
		var params armis.Params
		if !a.Params.IsNull() && !a.Params.IsUnknown() {
			var p paramsModel
			// Ensure correct extraction of params
			err := a.Params.As(context.TODO(), &p, basetypes.ObjectAsOptions{})
			if err != nil {
				// Skip if error occurs
				continue
			}
			params = convertParamsToAPI(p)
		}

		apiActions = append(apiActions, armis.Action{
			Type:   a.Type.ValueString(),
			Params: params,
		})
	}

	// Debug Log: Check extracted actions
	tflog.Debug(context.TODO(), "Converted Actions for API", map[string]any{
		"actions": apiActions,
	})

	return apiActions
}

// convertParamsToAPI converts Terraform params model to Armis API params.
func convertParamsToAPI(p paramsModel) armis.Params {
	var consolidation armis.Consolidation
	if !p.Consolidation.IsNull() && !p.Consolidation.IsUnknown() {
		var c consolidationModel
		p.Consolidation.As(context.TODO(), &c, basetypes.ObjectAsOptions{})
		consolidation = convertConsolidationToAPI(c)
	}

	// Populate tags
	var tags []string
	if !p.Tags.IsNull() && !p.Tags.IsUnknown() {
		p.Tags.ElementsAs(context.TODO(), &tags, false)
	}

	return armis.Params{
		Severity:      p.Severity.ValueString(),
		Title:         p.Title.ValueString(),
		Type:          p.Type.ValueString(),
		Endpoint:      p.Endpoint.ValueString(),
		Tags:          tags,
		Consolidation: consolidation,
	}
}

// convertConsolidationToAPI converts Terraform consolidation model to Armis API consolidation.
func convertConsolidationToAPI(c consolidationModel) armis.Consolidation {
	return armis.Consolidation{
		Amount: int(c.Amount.ValueInt64()),
		Unit:   c.Unit.ValueString(),
	}
}

func convertStringSliceToInterface(elements []string) []any {
	interfaces := make([]any, len(elements))
	for i, v := range elements {
		interfaces[i] = v
	}
	return interfaces
}

func convertSliceToStringSlice(in []any) []string {
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = fmt.Sprint(v)
	}
	return out
}

func convertToStringOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

func convertToIntOrNull(i int64) types.Int64 {
	if i == 0 {
		return types.Int64Null()
	}
	return types.Int64Value(i)
}

var consolidationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"amount": types.Int64Type,
		"unit":   types.StringType,
	},
}

var paramsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"severity":      types.StringType,
		"title":         types.StringType,
		"type":          types.StringType,
		"endpoint":      types.StringType,
		"tags":          types.ListType{ElemType: types.StringType}, // schema says list, not set
		"consolidation": consolidationObjectType,
	},
}

var actionObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type":   types.StringType,
		"params": paramsObjectType,
	},
}

// actionsListFromAPI converts Armis []Action into a Terraform list(object)
// while normalising empty/zero values to null using convertToStringOrNull and IntOrNull.
func actionsListFromAPI(
	ctx context.Context,
	api []armis.Action,
) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var elems []actionModel

	for _, a := range api {
		tagList, d := types.ListValueFrom(ctx, types.StringType, a.Params.Tags)
		diags.Append(d...)

		amt := convertToIntOrNull(int64(a.Params.Consolidation.Amount))
		unit := convertToStringOrNull(a.Params.Consolidation.Unit)
		var consObj types.Object
		if amt.IsNull() && unit.IsNull() {
			consObj = types.ObjectNull(consolidationObjectType.AttrTypes)
		} else {
			consObj, d = types.ObjectValueFrom(ctx,
				consolidationObjectType.AttrTypes,
				consolidationModel{Amount: amt, Unit: unit},
			)
		}
		diags.Append(d...)

		paramsObj, d := types.ObjectValueFrom(ctx, paramsObjectType.AttrTypes, paramsModel{
			Severity:      convertToStringOrNull(a.Params.Severity),
			Title:         convertToStringOrNull(a.Params.Title),
			Type:          convertToStringOrNull(a.Params.Type),
			Endpoint:      convertToStringOrNull(a.Params.Endpoint),
			Tags:          tagList,
			Consolidation: consObj,
		})
		diags.Append(d...)

		elems = append(elems, actionModel{
			Type:   convertToStringOrNull(a.Type),
			Params: paramsObj,
		})
	}

	listVal, d := types.ListValueFrom(ctx, actionObjectType, elems)
	diags.Append(d...)

	return listVal, diags
}

// responseToPolicy converts the Armis API payload returned by GetPolicy
// into the Terraform state model, returning any conversion diagnostics.
func responseToPolicy(
	ctx context.Context,
	p *armis.GetPolicySettings,
) (policyResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	labels, d := types.SetValueFrom(ctx, types.StringType, p.Labels)
	diags.Append(d...)

	mitreLabels, d := types.SetValueFrom(ctx, types.StringType, p.MitreAttackLabels)
	diags.Append(d...)

	var andList types.List
	if len(p.Rules.And) == 0 {
		andList = types.ListNull(types.StringType)
	} else {
		andStrings := convertSliceToStringSlice(p.Rules.And)
		l, d2 := types.ListValueFrom(ctx, types.StringType, andStrings)
		diags.Append(d2...)
		andList = l
	}

	orStrings := convertSliceToStringSlice(p.Rules.Or)
	orList, d2 := types.ListValueFrom(ctx, types.StringType, orStrings)
	diags.Append(d2...)

	actionsVal, d := actionsListFromAPI(ctx, p.Actions)
	diags.Append(d...)

	model := policyResourceModel{
		Name:              types.StringValue(p.Name),
		Description:       types.StringValue(p.Description),
		IsEnabled:         types.BoolValue(p.IsEnabled),
		RuleType:          types.StringValue(p.RuleType),
		Labels:            types.List(labels),
		MitreAttackLabels: types.List(mitreLabels),
		Rules: rulesModel{
			And: andList,
			Or:  orList,
		},
		Actions: actionsVal,
	}

	return model, diags
}
