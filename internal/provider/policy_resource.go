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

	newPolicy, err := r.client.CreatePolicy(policy)
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

	policy, err := r.client.GetPolicy(state.ID.ValueString())
	if err != nil {
		// Handle 404 Not Found by removing resource from state
		if strings.Contains(err.Error(), "status: 404") {
			tflog.Warn(ctx, "Policy not found, removing from state", map[string]any{
				"policy_id": state.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading policy", fmt.Sprintf("Failed to fetch policy: %v", err))
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
	state.Name = types.StringValue(policy.Name)
	state.Description = types.StringValue(policy.Description)
	state.IsEnabled = types.BoolValue(policy.IsEnabled)
	state.RuleType = types.StringValue(policy.RuleType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

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

	_, err := r.client.UpdatePolicy(policy, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating policy", fmt.Sprintf("API error: %v", err))
		return
	}

	plan.ID = state.ID

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *policyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyResourceModel
	if diags := req.State.Get(ctx, &state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if success, err := r.client.DeletePolicy(state.ID.ValueString()); err != nil || !success {
		resp.Diagnostics.AddError("Error deleting policy", err.Error())
	}
}

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
