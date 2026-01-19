// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/1898andCo/armis-sdk-go/armis"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &boundariesDataSource{}
	_ datasource.DataSourceWithConfigure = &boundariesDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *boundariesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

// BoundaryDataSource is a helper function to simplify the provider implementation.
func BoundaryDataSource() datasource.DataSource {
	return &boundariesDataSource{}
}

// boundariesDataSource is the data source implementation.
type boundariesDataSource struct {
	client *armis.Client
}

// Metadata returns the data source type name.
func (d *boundariesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_boundary"
}

// Schema defines the schema for the boundaries data source.
func (d *boundariesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Armis boundary information. If a boundary ID is provided, the data source will retrieve information about the specified boundary",
		Attributes: map[string]schema.Attribute{
			"boundary_id": schema.StringAttribute{
				Description: "An optional boundary ID used to filter the retrieved boundary information. If specified, only the boundary matching this ID will be returned.",
				Optional:    true,
			},
			"boundaries": schema.ListNestedAttribute{
				Description: "A computed list of boundaries. Each object in the list contains detailed information about a boundary.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "A unique identifier for the boundary.",
							Computed:    true,
						},
						"affected_sites": schema.StringAttribute{
							Description: "The sites affected by this boundary.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the boundary.",
							Computed:    true,
						},
						"rule_aql": schema.SingleNestedAttribute{
							Description: "The AQL rule configuration for the boundary.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"and": schema.ListAttribute{
									Description: "List of AND conditions in the AQL rule.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"or": schema.ListAttribute{
									Description: "List of OR conditions in the AQL rule.",
									ElementType: types.StringType,
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// boundariesDataSourceModel maps the data source schema data.
type boundariesDataSourceModel struct {
	BoundaryID types.String    `tfsdk:"boundary_id"`
	Boundaries []boundaryModel `tfsdk:"boundaries"`
}

// boundaryModel maps the boundary schema data.
type boundaryModel struct {
	ID            types.String `tfsdk:"id"`
	AffectedSites types.String `tfsdk:"affected_sites"`
	Name          types.String `tfsdk:"name"`
	RuleAQL       ruleAQLModel `tfsdk:"rule_aql"`
}

// ruleAQLModel maps the rule AQL schema data.
type ruleAQLModel struct {
	And []types.String `tfsdk:"and"`
	Or  []types.String `tfsdk:"or"`
}

// Read refreshes the Terraform state with the latest data.
func (d *boundariesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config boundariesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var boundaries []boundaryModel

	if !config.BoundaryID.IsNull() {
		// Fetch a specific boundary by ID
		boundary, err := d.client.GetBoundaryByID(ctx, config.BoundaryID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Armis Boundary",
				err.Error(),
			)
			return
		}

		// Map response body to model
		boundaryState := boundaryModel{
			ID:            types.StringValue(fmt.Sprintf("%d", boundary.ID)),
			AffectedSites: types.StringValue(boundary.AffectedSites),
			Name:          types.StringValue(boundary.Name),
		}

		// Map rule AQL
		var andConditions []types.String
		for _, condition := range boundary.RuleAQL.And {
			andConditions = append(andConditions, types.StringValue(condition))
		}

		var orConditions []types.String
		for _, condition := range boundary.RuleAQL.Or {
			orConditions = append(orConditions, types.StringValue(condition))
		}

		boundaryState.RuleAQL = ruleAQLModel{
			And: andConditions,
			Or:  orConditions,
		}

		boundaries = append(boundaries, boundaryState)
	} else {
		// Fetch all boundaries
		allBoundaries, err := d.client.GetBoundaries(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Armis Boundaries",
				err.Error(),
			)
			return
		}

		// Map response body to model
		for _, boundary := range allBoundaries {
			boundaryState := boundaryModel{
				ID:            types.StringValue(fmt.Sprintf("%d", boundary.ID)),
				AffectedSites: types.StringValue(boundary.AffectedSites),
				Name:          types.StringValue(boundary.Name),
			}

			// Map rule AQL
			var andConditions []types.String
			for _, condition := range boundary.RuleAQL.And {
				andConditions = append(andConditions, types.StringValue(condition))
			}

			var orConditions []types.String
			for _, condition := range boundary.RuleAQL.Or {
				orConditions = append(orConditions, types.StringValue(condition))
			}

			boundaryState.RuleAQL = ruleAQLModel{
				And: andConditions,
				Or:  orConditions,
			}

			boundaries = append(boundaries, boundaryState)
		}
	}

	// Save data into Terraform state
	config.Boundaries = boundaries
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
