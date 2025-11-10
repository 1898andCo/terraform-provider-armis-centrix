// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

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
	resp.TypeName = req.ProviderTypeName + "_lists"
}

// Schema defines the schema for the lists data source.
func (d *listsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Armis list information.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Optional filter by list name.",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "Optional filter by list type.",
				Optional:    true,
			},
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
						"list_type": schema.StringAttribute{
							Description: "The type of the list.",
							Computed:    true,
						},
						"creation_time": schema.StringAttribute{
							Description: "Creation time of the list.",
							Computed:    true,
						},
						"last_update_time": schema.StringAttribute{
							Description: "Last update time of the list.",
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
	Name  types.String `tfsdk:"name"`
	Type  types.String `tfsdk:"type"`
	Lists []listModel  `tfsdk:"lists"`
}

// listModel maps the list schema data.
type listModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	LastUpdatedBy  types.String `tfsdk:"last_updated_by"`
	CreatedBy      types.String `tfsdk:"created_by"`
	ListType       types.String `tfsdk:"list_type"`
	CreationTime   types.String `tfsdk:"creation_time"`
	LastUpdateTime types.String `tfsdk:"last_update_time"`
}

// Read refreshes the Terraform state with the latest data.
func (d *listsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state listsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch all lists
	apiLists, err := d.client.GetLists(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Armis Lists",
			err.Error(),
		)
		return
	}

	// Apply optional filters
	filtered := apiLists
	if !state.Name.IsNull() {
		name := state.Name.ValueString()
		var tmp []armis.ListSettings
		for _, l := range filtered {
			if l.ListName == name {
				tmp = append(tmp, l)
			}
		}
		filtered = tmp
	}
	if !state.Type.IsNull() {
		t := state.Type.ValueString()
		var tmp []armis.ListSettings
		for _, l := range filtered {
			if l.ListType == t {
				tmp = append(tmp, l)
			}
		}
		filtered = tmp
	}

	// Map response body to model
	var lists []listModel
	for _, l := range filtered {
		lists = append(lists, listModel{
			ID:             types.StringValue(strconv.Itoa(l.ListID)),
			Name:           types.StringValue(l.ListName),
			Description:    types.StringValue(l.Description),
			LastUpdatedBy:  types.StringValue(l.LastUpdatedBy),
			CreatedBy:      types.StringValue(l.CreatedBy),
			ListType:       types.StringValue(l.ListType),
			CreationTime:   types.StringValue(l.CreationTime),
			LastUpdateTime: types.StringValue(l.LastUpdateTime),
		})
	}

	// Save data into Terraform state
	state.Lists = lists
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
