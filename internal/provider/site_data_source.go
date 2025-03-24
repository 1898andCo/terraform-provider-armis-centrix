// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/1898andCo/terraform-provider-armis-centrix/internal/armis"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sitesDataSource{}
	_ datasource.DataSourceWithConfigure = &sitesDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *sitesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// SiteDataSource is a helper function to simplify the provider implementation.
func SiteDataSource() datasource.DataSource {
	return &sitesDataSource{}
}

// sitesDataSource is the data source implementation.
type sitesDataSource struct {
	client *armis.Client
}

// Metadata returns the data source type name.
func (d *sitesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sites"
}

// Schema defines the schema for the sites data source.
func (d *sitesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Armis site information",
		Attributes: map[string]schema.Attribute{
			"sites": schema.ListNestedAttribute{
				Description: "A computed list of sites. Each object in the list contains detailed information about a site.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "A unique identifier for the site.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the site.",
							Computed:    true,
						},
						"latitude": schema.Float64Attribute{
							Description: "The latitude coordinate of the site.",
							Computed:    true,
						},
						"longitude": schema.Float64Attribute{
							Description: "The longitude coordinate of the site.",
							Computed:    true,
						},
						"location": schema.StringAttribute{
							Description: "The physical location or address of the site.",
							Computed:    true,
						},
						"tier": schema.StringAttribute{
							Description: "The tier classification of the site, which may represent its priority or categorization.",
							Computed:    true,
						},
						"user": schema.StringAttribute{
							Description: "The user associated with or responsible for the site.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// sitesDataSourceModel maps the data source schema data.
type sitesDataSourceModel struct {
	Sites []sitesModel `tfsdk:"sites"`
}

// sitesModel maps sites schema data.
type sitesModel struct {
	ID        types.String  `tfsdk:"id"`
	Name      types.String  `tfsdk:"name"`
	Latitude  types.Float64 `tfsdk:"latitude"`
	Longitude types.Float64 `tfsdk:"longitude"`
	Location  types.String  `tfsdk:"location"`
	Tier      types.String  `tfsdk:"tier"`
	User      types.String  `tfsdk:"user"`
}

// Read refreshes the Terraform state with the latest data.
func (d *sitesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sitesDataSourceModel

	sites, err := d.client.GetSites(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Armis Sites",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, site := range sites {
		siteState := sitesModel{
			ID:        types.StringValue(site.ID),
			Name:      types.StringValue(site.Name),
			Latitude:  types.Float64Value(site.Lat),
			Longitude: types.Float64Value(site.Lng),
			Location:  types.StringValue(site.Location),
			Tier:      types.StringValue(site.Tier),
			User:      types.StringValue(site.User),
		}

		state.Sites = append(state.Sites, siteState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
