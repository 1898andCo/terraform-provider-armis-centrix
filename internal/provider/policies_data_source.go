// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/1898andCo/terraform-provider-armis-centrix/armis"
	u "github.com/1898andCo/terraform-provider-armis-centrix/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &policiesDataSource{}
	_ datasource.DataSourceWithConfigure = &policiesDataSource{}
)

type policiesDataSource struct {
	client *armis.Client
}

func PoliciesDataSource() datasource.DataSource {
	return &policiesDataSource{}
}

func (d *policiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *policiesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policies"
}

func (d *policiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Armis policies. Filter by policy_id or match_prefix to limit the results.",
		Attributes: map[string]schema.Attribute{
			"policy_id": schema.StringAttribute{
				Optional:    true,
				Description: "An optional policy identifier used to filter the retrieved policy information.",
			},
			"match_prefix": schema.StringAttribute{
				Optional:    true,
				Description: "Optional prefix to match policy names. When provided, only policies with names starting with this prefix are returned.",
			},
			"policies": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A computed list of Armis policies.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier of the policy.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the policy.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Description of the policy.",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Indicates whether the policy is enabled.",
						},
						"rule_type": schema.StringAttribute{
							Computed:    true,
							Description: "Rule type associated with the policy.",
						},
						"labels": schema.ListAttribute{
							Computed:    true,
							Description: "Labels assigned to the policy.",
							ElementType: types.StringType,
						},
						"mitre_attack_labels": schema.ListNestedAttribute{
							Computed:    true,
							Description: "MITRE ATT&CK labels associated with the policy.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"matrix": schema.StringAttribute{
										Computed:    true,
										Description: "MITRE matrix name.",
									},
									"sub_technique": schema.StringAttribute{
										Computed:    true,
										Description: "MITRE sub-technique identifier.",
									},
									"tactic": schema.StringAttribute{
										Computed:    true,
										Description: "MITRE tactic name.",
									},
									"technique": schema.StringAttribute{
										Computed:    true,
										Description: "MITRE technique name.",
									},
								},
							},
						},
						"actions": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Actions configured on the policy.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Computed:    true,
										Description: "Type of the action.",
									},
									"params": schema.SingleNestedAttribute{
										Computed:    true,
										Description: "Parameters of the action.",
										Attributes: map[string]schema.Attribute{
											"severity": schema.StringAttribute{
												Computed:    true,
												Description: "Severity assigned to the action.",
											},
											"title": schema.StringAttribute{
												Computed:    true,
												Description: "Title associated with the action.",
											},
											"type": schema.StringAttribute{
												Computed:    true,
												Description: "Category of the action.",
											},
											"endpoint": schema.StringAttribute{
												Computed:    true,
												Description: "Endpoint targeted by the action.",
											},
											"tags": schema.ListAttribute{
												Computed:    true,
												Description: "Tags applied to the action.",
												ElementType: types.StringType,
											},
											"consolidation": schema.SingleNestedAttribute{
												Computed:    true,
												Description: "Consolidation settings for the action.",
												Attributes: map[string]schema.Attribute{
													"amount": schema.Int64Attribute{
														Computed:    true,
														Description: "Consolidation amount.",
													},
													"unit": schema.StringAttribute{
														Computed:    true,
														Description: "Consolidation unit.",
													},
												},
											},
										},
									},
								},
							},
						},
						"rules": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Rules evaluated by the policy.",
							Attributes: map[string]schema.Attribute{
								"and": schema.ListAttribute{
									Computed:    true,
									Description: "AND conditions configured for the policy.",
									ElementType: types.StringType,
								},
								"or": schema.ListAttribute{
									Computed:    true,
									Description: "OR conditions configured for the policy.",
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *policiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config u.PoliciesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var policies []u.PolicyDataSourcePolicyModel

	if !config.PolicyID.IsNull() && config.PolicyID.ValueString() != "" {
		policy, err := d.client.GetPolicy(ctx, config.PolicyID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Armis Policy",
				err.Error(),
			)
			return
		}

		model := u.BuildPolicyDataSourceModelFromGet(policy, config.PolicyID.ValueString())
		if u.ShouldIncludePolicy(model, config.MatchPrefix) {
			policies = append(policies, model)
		}
	} else {
		allPolicies, err := d.client.GetAllPolicies(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Armis Policies",
				err.Error(),
			)
			return
		}

		for _, policy := range allPolicies {
			model := u.BuildPolicyDataSourceModelFromSingle(policy)
			if u.ShouldIncludePolicy(model, config.MatchPrefix) {
				policies = append(policies, model)
			}
		}
	}

	config.Policies = policies

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
