// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/1898andCo/terraform-provider-armis-centrix/armis"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &listsDataSource{}
	_ datasource.DataSourceWithConfigure = &listsDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *listsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// ListsDataSource is a helper function to simplify the provider implementation.
func ListsDataSource() datasource.DataSource {
	return &listsDataSource{}
}

// listsDataSource is the data source implementation.
type listsDataSource struct {
	client *armis.Client
}

// Metadata returns the data source type name.
func (d *listsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_list"
}

/*
type ListSettings struct {
	CreatedBy      string `json:"created_by"`
	CreationTime   string `json:"creation_time"`
	Description    string `json:"description"`
	LastUpdateTime string `json:"last_update_time"`
	LastUpdatedBy  string `json:"last_updated_by"`
	ListID         int    `json:"list_id"`
	ListName       string `json:"list_name"`
	ListType       string `json:"list_type"`
}
*/

// Schema defines the schema for the lists data source.
func (d *listsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Armis list information.",
		Attributes: map[string]schema.Attribute{
			"lists": schema.ListNestedAttribute{
				Description: "A computed list of lists. Each object in the list contains detailed information about a list.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "A unique identifier for the list.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the list.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the list.",
							Computed:    true,
						},
						"created_by": schema.StringAttribute{
							Description: "The user who created the list.",
							Computed:    true,
						},
						"last_updated_by": schema.StringAttribute{
							Description: "The user who last updated the list.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// listsDataSourceModel maps the data source schema data.
type listsDataSourceModel struct {
	lists []listModel `tfsdk:"lists"`
}

// listModel maps the list schema data.
type listModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	LastUpdatedBy types.String `tfsdk:"last_updated_by"`
	CreatedBy     types.String `tfsdk:"created_by"`
}

// Read refreshes the Terraform state with the latest data.
func (d *listsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config listsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var lists []boundaryModel

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

		lists = append(lists, boundaryState)
	} else {
		// Fetch all lists
		alllists, err := d.client.Getlists(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Armis lists",
				err.Error(),
			)
			return
		}

		// Map response body to model
		for _, boundary := range alllists {
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

			lists = append(lists, boundaryState)
		}
	}

	// Save data into Terraform state
	config.lists = lists
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
