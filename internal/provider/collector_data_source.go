// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/1898andCo/terraform-provider-armis-centrix/internal/armis"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &collectorsDataSource{}
	_ datasource.DataSourceWithConfigure = &collectorsDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *collectorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// CollectorDataSource is a helper function to simplify the provider implementation.
func CollectorDataSource() datasource.DataSource {
	return &collectorsDataSource{}
}

// collectorsDataSource is the data source implementation.
type collectorsDataSource struct {
	client *armis.Client
}

// Metadata returns the data source type name.
func (d *collectorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collectors"
}

// Schema defines the schema for the collectors data source.
func (d *collectorsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Armis collector information",
		Attributes: map[string]schema.Attribute{
			"collectors": schema.ListNestedAttribute{
				Description: "A list of Armis collectors. Each object contains detailed information about a collector.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"boot_time": schema.StringAttribute{
							Description: "The boot time of the collector, represented as a timestamp.",
							Computed:    true,
						},
						"city": schema.StringAttribute{
							Description: "The city where the collector is located.",
							Computed:    true,
						},
						"cluster_id": schema.NumberAttribute{
							Description: "The unique identifier of the cluster to which the collector belongs.",
							Computed:    true,
						},
						"collector_number": schema.NumberAttribute{
							Description: "A number that uniquely identifies the collector.",
							Computed:    true,
						},
						"country": schema.StringAttribute{
							Description: "The country where the collector is located.",
							Computed:    true,
						},
						"default_gateway": schema.StringAttribute{
							Description: "The default gateway associated with the collector.",
							Computed:    true,
						},
						"https_proxy_redacted": schema.StringAttribute{
							Description: "The HTTPS proxy configuration for the collector, with sensitive information redacted.",
							Computed:    true,
						},
						"ip_address": schema.StringAttribute{
							Description: "The IP address assigned to the collector.",
							Computed:    true,
						},
						"last_seen": schema.StringAttribute{
							Description: "The last time the collector was active, represented as a timestamp.",
							Computed:    true,
						},
						"mac_address": schema.StringAttribute{
							Description: "The MAC address of the collector's network interface.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the collector.",
							Computed:    true,
						},
						"namespace": schema.StringAttribute{
							Description: "The namespace within which the collector operates.",
							Computed:    true,
						},
						"product_serial": schema.StringAttribute{
							Description: "The serial number of the collector's hardware.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The current status of the collector (e.g., active, inactive).",
							Computed:    true,
						},
						"subnet": schema.StringAttribute{
							Description: "The subnet associated with the collector.",
							Computed:    true,
						},
						"system_vendor": schema.StringAttribute{
							Description: "The vendor of the system running the collector.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The type of collector, indicating its purpose or configuration.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// collectorsDataSourceModel maps the data source schema data.
type collectorsDataSourceModel struct {
	Collectors []collectorModel `tfsdk:"collectors"`
}

// collectorModel maps the collector schema data.
type collectorModel struct {
	BootTime         types.String `tfsdk:"boot_time"`
	City             types.String `tfsdk:"city"`
	ClusterID        types.Number `tfsdk:"cluster_id"`
	CollectorNumber  types.Number `tfsdk:"collector_number"`
	Country          types.String `tfsdk:"country"`
	DefaultGateway   types.String `tfsdk:"default_gateway"`
	HTTPSProxyRedact types.String `tfsdk:"https_proxy_redacted"`
	IPAddress        types.String `tfsdk:"ip_address"`
	LastSeen         types.String `tfsdk:"last_seen"`
	MACAddress       types.String `tfsdk:"mac_address"`
	Name             types.String `tfsdk:"name"`
	Namespace        types.String `tfsdk:"namespace"`
	ProductSerial    types.String `tfsdk:"product_serial"`
	Status           types.String `tfsdk:"status"`
	Subnet           types.String `tfsdk:"subnet"`
	SystemVendor     types.String `tfsdk:"system_vendor"`
	Type             types.String `tfsdk:"type"`
}

// Read refreshes the Terraform state with the latest data.
func (d *collectorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state collectorsDataSourceModel

	collectors, err := d.client.GetCollectors(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Armis Collectors",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, collector := range collectors {
		collectorState := collectorModel{
			BootTime:         types.StringValue(collector.BootTime),
			City:             types.StringValue(collector.City),
			ClusterID:        types.NumberValue(big.NewFloat(float64(collector.ClusterID))),
			CollectorNumber:  types.NumberValue(big.NewFloat(float64(collector.CollectorNumber))),
			Country:          types.StringValue(collector.Country),
			DefaultGateway:   types.StringValue(collector.DefaultGateway),
			HTTPSProxyRedact: types.StringValue(collector.HTTPSProxyRedacted),
			IPAddress:        types.StringValue(collector.IPAddress),
			LastSeen:         types.StringValue(collector.LastSeen),
			MACAddress:       types.StringValue(collector.MacAddress),
			Name:             types.StringValue(collector.Name),
			Namespace:        types.StringValue(collector.Namespace),
			ProductSerial:    types.StringValue(collector.ProductSerial),
			Status:           types.StringValue(collector.Status),
			Subnet:           types.StringValue(collector.Subnet),
			SystemVendor:     types.StringValue(collector.SystemVendor),
			Type:             types.StringValue(collector.Type),
		}

		state.Collectors = append(state.Collectors, collectorState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
