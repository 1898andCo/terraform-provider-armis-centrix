// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/1898andCo/armis-sdk-go/v2/armis"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &tagsDataSource{}
	_ datasource.DataSourceWithConfigure = &tagsDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *tagsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// TagsDataSource is a helper function to simplify the provider implementation.
func TagsDataSource() datasource.DataSource {
	return &tagsDataSource{}
}

// tagsDataSource is the data source implementation.
type tagsDataSource struct {
	client *armis.Client
}

// Metadata returns the data source type name.
func (d *tagsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tags"
}

// Schema defines the schema for the tags data source.
func (d *tagsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Armis tags. Use match_prefix and/or exclude_prefix to filter results.",
		Attributes: map[string]schema.Attribute{
			"match_prefix": schema.StringAttribute{
				Description: "Optional filter to include only tags starting with this prefix.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"exclude_prefix": schema.StringAttribute{
				Description: "Optional filter to exclude tags starting with this prefix.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"tags": schema.ListAttribute{
				Description: "A computed list of Armis tags.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// tagsDataSourceModel maps the data source schema data.
type tagsDataSourceModel struct {
	MatchPrefix   types.String `tfsdk:"match_prefix"`
	ExcludePrefix types.String `tfsdk:"exclude_prefix"`
	Tags          types.List   `tfsdk:"tags"`
}

// Read refreshes the Terraform state with the latest data.
func (d *tagsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tagsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch all tags
	apiTags, err := d.client.GetTags(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Armis Tags",
			err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "Fetched tags from Armis API", map[string]any{
		"total_count": len(apiTags),
	})

	// Apply optional filters
	filtered := apiTags
	if !state.MatchPrefix.IsNull() {
		prefix := state.MatchPrefix.ValueString()
		beforeCount := len(filtered)
		var tmp []string
		for _, tag := range filtered {
			if strings.HasPrefix(tag, prefix) {
				tmp = append(tmp, tag)
			}
		}
		filtered = tmp
		tflog.Debug(ctx, "Applied match_prefix filter", map[string]any{
			"prefix":       prefix,
			"before_count": beforeCount,
			"after_count":  len(filtered),
		})
	}
	if !state.ExcludePrefix.IsNull() {
		prefix := state.ExcludePrefix.ValueString()
		beforeCount := len(filtered)
		var tmp []string
		for _, tag := range filtered {
			if !strings.HasPrefix(tag, prefix) {
				tmp = append(tmp, tag)
			}
		}
		filtered = tmp
		tflog.Debug(ctx, "Applied exclude_prefix filter", map[string]any{
			"prefix":       prefix,
			"before_count": beforeCount,
			"after_count":  len(filtered),
		})
	}

	// Map response to model
	tagValues := make([]attr.Value, len(filtered))
	for i, tag := range filtered {
		tagValues[i] = types.StringValue(tag)
	}

	tagsList, diags := types.ListValue(types.StringType, tagValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	state.Tags = tagsList
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
