// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

// Package provider implements the Armis Terraform provider
// for managing Armis Centrix security platform resources.
// It provides the necessary configuration, resources, and data sources
// to interact with the Armis API.
package provider

import (
	"context"
	"os"

	armis "github.com/1898andCo/terraform-provider-armis-centrix/internal/armis"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure ArmisProvider satisfies various provider interfaces.
var (
	_ provider.Provider = &ArmisProvider{}
)

// ArmisProvider defines the provider implementation.
type ArmisProvider struct {
	version string
}

// ArmisProviderModel describes the provider data model.
type ArmisProviderModel struct {
	APIUrl types.String `tfsdk:"api_url"`
	APIKey types.String `tfsdk:"api_key"`
}

func (p *ArmisProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	// This will remain the same across Centrix and Apex implementations similar to the google and google-beta provider
	resp.TypeName = "armis"
	resp.Version = p.version
}

func (p *ArmisProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Armis Centrix Terraform Provider allows you to manage and automate Armis Centrix security platform resources",
		Attributes: map[string]schema.Attribute{
			// These attributes are optional because they can be set via environment variables
			"api_url": schema.StringAttribute{
				MarkdownDescription: "URL endpoint for the Armis API.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API Key for the Armis API.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *ArmisProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	tflog.Info(ctx, "Configuring the Armis client")

	var config ArmisProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle unknown values
	if config.APIUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Unknown Armis API URL",
			"The provider cannot create the Armis API client as there is an unknown configuration value for the API URL. "+
				"Set this value or use the ARMIS_API_URL environment variable.",
		)
	}

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Armis API Key",
			"The provider cannot create the Armis API client as there is an unknown configuration value for the API key. "+
				"Set this value or use the ARMIS_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default to environment variables if configuration is not provided
	apiKey := os.Getenv("ARMIS_API_KEY")
	apiUrl := os.Getenv("ARMIS_API_URL")

	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	if !config.APIUrl.IsNull() {
		apiUrl = config.APIUrl.ValueString()
	}

	// Handle missing values
	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Armis API Key",
			"The provider cannot create the Armis API client as there is a missing API key. "+
				"Set this value in the configuration or use the ARMIS_API_KEY environment variable.",
		)
	}

	if apiUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Missing Armis API URL",
			"The provider cannot create the Armis API client as there is a missing API URL. "+
				"Set this value in the configuration or use the ARMIS_API_URL environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Debug log the configuration
	ctx = tflog.SetField(ctx, "armis_api_url", apiUrl)
	ctx = tflog.SetField(ctx, "armis_api_key", apiKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "armis_api_key")

	tflog.Debug(ctx, "Creating the Armis API client")

	// Create the Armis client
	client, err := armis.NewClient(
		apiKey,
		armis.WithAPIURL(apiUrl),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Armis API Client",
			"An error occurred when creating the Armis API client: "+err.Error(),
		)
		return
	}

	// Make the client available to data sources and resources
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Armis API client created", map[string]any{"success": true})
}

func (p *ArmisProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		RoleResource,
		UserResource,
		CollectorResource,
		PolicyResource,
	}
}

func (p *ArmisProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		RoleDataSource,
		SiteDataSource,
		UserDataSource,
		CollectorDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ArmisProvider{
			version: version,
		}
	}
}
