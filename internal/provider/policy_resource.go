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
	u "github.com/1898andCo/terraform-provider-armis-centrix/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &policyResource{}
	_ resource.ResourceWithConfigure   = &policyResource{}
	_ resource.ResourceWithImportState = &policyResource{}
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
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
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

// Create decodes the plan into a model, converts it to an Armis
// PolicySettings payload, invokes r.client.CreatePolicy, stores the returned
// policy ID in state, and writes the updated state back—aborting early whenever
// diagnostics report an error.
func (r *policyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan u.PolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, diags := u.BuildPolicySettings(&plan)
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
	var state u.PolicyResourceModel
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
	result := u.ResponseToPolicyFromGet(ctx, getResp)

	// Preserve the ID and MITRE labels from state
	result.ID = state.ID
	result.MitreAttackLabels = state.MitreAttackLabels
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

// Update loads plan and state, maps the plan to an Armis PolicySettings
// payload, calls r.client.UpdatePolicy with the existing ID, and writes the
// (unchanged-ID) state back—bailing out on any diagnostics or API error.
func (r *policyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state u.PolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, diags := u.BuildPolicySettings(&plan)
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
	result := u.ResponseToPolicyFromUpdate(ctx, updateResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the ID and MITRE labels
	result.ID = state.ID
	result.MitreAttackLabels = plan.MitreAttackLabels
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

// Delete removes a policy of the provided ID.
func (r *policyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state u.PolicyResourceModel
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

func (r *policyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
