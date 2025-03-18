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
	_ datasource.DataSource              = &rolesDataSource{}
	_ datasource.DataSourceWithConfigure = &rolesDataSource{}
)

// rolesDataSource is the data source implementation.
type rolesDataSource struct {
	client *armis.Client
}

// Configure adds the provider configured client to the data source.
func (d *rolesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// RoleDataSource is a helper function to simplify the provider implementation.
func RoleDataSource() datasource.DataSource {
	return &rolesDataSource{}
}

// Metadata returns the data source type name.
func (d *rolesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema defines the schema for the roles data source.
func (d *rolesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an Armis role",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the role.",
			},
			"role_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the role.",
			},
			"vipr_role": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates if the role is a VIPR-specific role.",
			},
			"permissions": schema.SingleNestedAttribute{
				Description: "Permissions associated with the role.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"advanced_permissions": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Advanced permissions for the role.",
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Computed:    true,
								Description: "Indicates if the role has all advanced permissions.",
							},
							"behavioral": schema.SingleNestedAttribute{
								Computed:    true,
								Description: "Behavioral permissions for the role.",
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Computed:    true,
										Description: "Indicates if the role has all behavioral permissions.",
									},
									"application_name": schema.BoolAttribute{
										Computed:    true,
										Description: "Permission for application names.",
									},
									"host_name": schema.BoolAttribute{
										Computed:    true,
										Description: "Permission for host names.",
									},
									"service_name": schema.BoolAttribute{
										Computed:    true,
										Description: "Permission for service names.",
									},
								},
							},
							"device": schema.SingleNestedAttribute{
								Computed:    true,
								Description: "Device-related permissions.",
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Computed:    true,
										Description: "Indicates if the role has all device permissions.",
									},
									"device_names": schema.BoolAttribute{
										Computed:    true,
										Description: "Permission for device names.",
									},
									"ip_addresses": schema.BoolAttribute{
										Computed:    true,
										Description: "Permission for IP addresses.",
									},
									"mac_addresses": schema.BoolAttribute{
										Computed:    true,
										Description: "Permission for MAC addresses.",
									},
									"phone_numbers": schema.BoolAttribute{
										Computed:    true,
										Description: "Permission for phone numbers.",
									},
								},
							},
						},
					},
					"alert": schema.SingleNestedAttribute{
						Description: "Permissions for managing alerts.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all alert permissions are enabled.",
								Computed:    true,
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing alerts.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all alert management permissions are enabled.",
										Computed:    true,
									},
									"resolve": schema.BoolAttribute{
										Description: "Permission to resolve alerts.",
										Computed:    true,
									},
									"suppress": schema.BoolAttribute{
										Description: "Permission to suppress alerts.",
										Computed:    true,
									},
									"whitelist_devices": schema.BoolAttribute{
										Description: "Permission to whitelist devices in alerts.",
										Computed:    true,
									},
								},
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read alerts.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

// RoleSettingsModel maps the RoleSettings schema data.
type RoleDataSourceModel struct {
	Name        types.String      `tfsdk:"name"`
	Permissions *PermissionsModel `tfsdk:"permissions"`
	ID          types.String      `tfsdk:"role_id"`
	ViprRole    types.Bool        `tfsdk:"vipr_role"`
}

// Read refreshes the Terraform state with the latest data.
func (d *rolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config RoleDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Name.IsNull() || config.Name.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Role Name",
			"The 'name' attribute is required to fetch a specific role.",
		)
		return
	}

	role, err := d.client.GetRoleByName(config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Armis Role",
			fmt.Sprintf("Error fetching role: %s", err.Error()),
		)
		return
	}

	roleState := RoleDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("%d", role.ID)),
		Name:     types.StringValue(role.Name),
		ViprRole: types.BoolValue(role.ViprRole),
		Permissions: &PermissionsModel{
			AdvancedPermissions: &AdvancedPermissionsModel{
				All: types.BoolValue(role.Permissions.AdvancedPermissions.All),
				Behavioral: &BehavioralModel{
					All:             types.BoolValue(role.Permissions.AdvancedPermissions.Behavioral.All),
					ApplicationName: types.BoolValue(role.Permissions.AdvancedPermissions.Behavioral.ApplicationName.All),
					HostName:        types.BoolValue(role.Permissions.AdvancedPermissions.Behavioral.HostName.All),
					ServiceName:     types.BoolValue(role.Permissions.AdvancedPermissions.Behavioral.ServiceName.All),
				},
				Device: &DeviceAdvancedModel{
					All:          types.BoolValue(role.Permissions.AdvancedPermissions.Device.All),
					DeviceNames:  types.BoolValue(role.Permissions.AdvancedPermissions.Device.DeviceNames.All),
					IPAddresses:  types.BoolValue(role.Permissions.AdvancedPermissions.Device.IPAddresses.All),
					MACAddresses: types.BoolValue(role.Permissions.AdvancedPermissions.Device.MACAddresses.All),
					PhoneNumbers: types.BoolValue(role.Permissions.AdvancedPermissions.Device.PhoneNumbers.All),
				},
			},
			Alert: &AlertModel{
				All: types.BoolValue(role.Permissions.Alert.All),
				Manage: &ManageModel{
					All:              types.BoolValue(role.Permissions.Alert.Manage.All),
					Resolve:          types.BoolValue(role.Permissions.Alert.Manage.Resolve.All),
					Suppress:         types.BoolValue(role.Permissions.Alert.Manage.Suppress.All),
					WhitelistDevices: types.BoolValue(role.Permissions.Alert.Manage.WhitelistDevices.All),
				},
				Read: types.BoolValue(role.Permissions.Alert.Read.All),
			},
		},
	}

	// Set the state with the fetched role
	resp.Diagnostics.Append(resp.State.Set(ctx, &roleState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
