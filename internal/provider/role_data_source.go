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
							"read": schema.SingleNestedAttribute{
								Description: "Permission to read alerts.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all read permissions are enabled.",
										Computed:    true,
									},
								},
							},
						},
					},
					"device": schema.SingleNestedAttribute{
						Description: "Permissions for managing devices.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all device permissions are enabled.",
								Computed:    true,
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing devices.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all device management permissions are enabled.",
										Computed:    true,
									},
									"create": schema.BoolAttribute{
										Description: "Permission to create devices.",
										Computed:    true,
									},
									"delete": schema.BoolAttribute{
										Description: "Permission to delete devices.",
										Computed:    true,
									},
									"edit": schema.BoolAttribute{
										Description: "Permission to edit devices.",
										Computed:    true,
									},
									"enforce": schema.SingleNestedAttribute{
										Description: "Permissions for enforcing device policies.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all enforce permissions are enabled.",
												Computed:    true,
											},
											"create": schema.BoolAttribute{
												Description: "Permission to create enforcement policies.",
												Computed:    true,
											},
											"delete": schema.BoolAttribute{
												Description: "Permission to delete enforcement policies.",
												Computed:    true,
											},
										},
									},
									"merge": schema.BoolAttribute{
										Description: "Permission to merge devices.",
										Computed:    true,
									},
									"request_deleted_data": schema.BoolAttribute{
										Description: "Permission to request deleted data.",
										Computed:    true,
									},
									"tags": schema.BoolAttribute{
										Description: "Permission to manage device tags.",
										Computed:    true,
									},
								},
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read devices.",
								Computed:    true,
							},
						},
					},
					"policy": schema.SingleNestedAttribute{
						Description: "Permissions for managing policies.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all policy permissions are enabled.",
								Computed:    true,
							},
							"manage": schema.BoolAttribute{
								Description: "Permission to manage policies.",
								Computed:    true,
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read policies.",
								Computed:    true,
							},
						},
					},
					"report": schema.SingleNestedAttribute{
						Description: "Permissions for managing reports.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all report permissions are enabled.",
								Computed:    true,
							},
							"export": schema.BoolAttribute{
								Description: "Permission to export reports.",
								Computed:    true,
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing reports.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all report management permissions are enabled.",
										Computed:    true,
									},
									"create": schema.BoolAttribute{
										Description: "Permission to create reports.",
										Computed:    true,
									},
									"delete": schema.BoolAttribute{
										Description: "Permission to delete reports.",
										Computed:    true,
									},
									"edit": schema.BoolAttribute{
										Description: "Permission to edit reports.",
										Computed:    true,
									},
								},
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read reports.",
								Computed:    true,
							},
						},
					},
					"risk_factor": schema.SingleNestedAttribute{
						Description: "Permissions for managing risk factors.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all risk factor permissions are enabled.",
								Computed:    true,
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing risk factors.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all risk factor management permissions are enabled.",
										Computed:    true,
									},
									"customization": schema.SingleNestedAttribute{
										Description: "Permissions for customizing risk factors.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all customization permissions are enabled.",
												Computed:    true,
											},
											"create": schema.BoolAttribute{
												Description: "Permission to create customizations.",
												Computed:    true,
											},
											"disable": schema.BoolAttribute{
												Description: "Permission to disable customizations.",
												Computed:    true,
											},
											"edit": schema.BoolAttribute{
												Description: "Permission to edit customizations.",
												Computed:    true,
											},
										},
									},
									"status": schema.SingleNestedAttribute{
										Description: "Permissions for managing risk factor status.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all status permissions are enabled.",
												Computed:    true,
											},
											"ignore": schema.BoolAttribute{
												Description: "Permission to ignore risk factors.",
												Computed:    true,
											},
											"resolve": schema.BoolAttribute{
												Description: "Permission to resolve risk factors.",
												Computed:    true,
											},
										},
									},
								},
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read risk factors.",
								Computed:    true,
							},
						},
					},
					"settings": schema.SingleNestedAttribute{
						Description: "Permissions for managing settings.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all settings permissions are enabled.",
								Computed:    true,
							},
							"audit_log": schema.BoolAttribute{
								Description: "Permission to access audit logs.",
								Computed:    true,
							},
							"boundary": schema.SingleNestedAttribute{
								Description: "Permissions for managing boundaries.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all boundary permissions are enabled.",
										Computed:    true,
									},
									"manage": schema.SingleNestedAttribute{
										Description: "Permissions for managing boundaries.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all boundary management permissions are enabled.",
												Computed:    true,
											},
											"create": schema.BoolAttribute{
												Description: "Permission to create boundaries.",
												Computed:    true,
											},
											"delete": schema.BoolAttribute{
												Description: "Permission to delete boundaries.",
												Computed:    true,
											},
											"edit": schema.BoolAttribute{
												Description: "Permission to edit boundaries.",
												Computed:    true,
											},
										},
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read boundaries.",
										Computed:    true,
									},
								},
							},
							"business_impact": schema.SingleNestedAttribute{
								Description: "Permissions for managing business impact.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all business impact permissions are enabled.",
										Computed:    true,
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage business impact.",
										Computed:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read business impact.",
										Computed:    true,
									},
								},
							},
							"collector": schema.SingleNestedAttribute{
								Description: "Permissions for managing collectors.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all collector permissions are enabled.",
										Computed:    true,
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage collectors.",
										Computed:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read collectors.",
										Computed:    true,
									},
								},
							},
							"custom_properties": schema.SingleNestedAttribute{
								Description: "Permissions for managing custom properties.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all custom properties permissions are enabled.",
										Computed:    true,
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage custom properties.",
										Computed:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read custom properties.",
										Computed:    true,
									},
								},
							},
							"integration": schema.SingleNestedAttribute{
								Description: "Permissions for managing integrations.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all integration permissions are enabled.",
										Computed:    true,
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage integrations.",
										Computed:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read integrations.",
										Computed:    true,
									},
								},
							},
							"internal_ips": schema.SingleNestedAttribute{
								Description: "Permissions for managing internal IPs.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all internal IPs permissions are enabled.",
										Computed:    true,
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage internal IPs.",
										Computed:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read internal IPs.",
										Computed:    true,
									},
								},
							},
							"notifications": schema.SingleNestedAttribute{
								Description: "Permissions for managing notifications.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all notifications permissions are enabled.",
										Computed:    true,
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage notifications.",
										Computed:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read notifications.",
										Computed:    true,
									},
								},
							},
							"oidc": schema.SingleNestedAttribute{
								Description: "Permissions for managing OIDC.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all OIDC permissions are enabled.",
										Computed:    true,
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage OIDC.",
										Computed:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read OIDC.",
										Computed:    true,
									},
								},
							},
							"saml": schema.SingleNestedAttribute{
								Description: "Permissions for managing SAML.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all SAML permissions are enabled.",
										Computed:    true,
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage SAML.",
										Computed:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read SAML.",
										Computed:    true,
									},
								},
							},
							"secret_key": schema.BoolAttribute{
								Description: "Permission to access secret keys.",
								Computed:    true,
							},
							"security_settings": schema.BoolAttribute{
								Description: "Permission to access security settings.",
								Computed:    true,
							},
							"sites_and_sensors": schema.SingleNestedAttribute{
								Description: "Permissions for managing sites and sensors.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all sites and sensors permissions are enabled.",
										Computed:    true,
									},
									"manage": schema.SingleNestedAttribute{
										Description: "Permissions for managing sites and sensors.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all manage permissions are enabled.",
												Computed:    true,
											},
											"sensors": schema.BoolAttribute{
												Description: "Permission to manage sensors.",
												Computed:    true,
											},
											"sites": schema.BoolAttribute{
												Description: "Permission to manage sites.",
												Computed:    true,
											},
										},
									},
									"read": schema.SingleNestedAttribute{
										Description: "Permission to read sites and sensors.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all read permissions are enabled.",
												Computed:    true,
											},
										},
									},
								},
							},
							"users_and_roles": schema.SingleNestedAttribute{
								Description: "Permissions for managing users and roles.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all users and roles permissions are enabled.",
										Computed:    true,
									},
									"manage": schema.SingleNestedAttribute{
										Description: "Permissions for managing users and roles.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all manage permissions are enabled.",
												Computed:    true,
											},
											"roles": schema.SingleNestedAttribute{
												Description: "Permissions for managing roles.",
												Computed:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all role permissions are enabled.",
														Computed:    true,
													},
													"create": schema.BoolAttribute{
														Description: "Permission to create roles.",
														Computed:    true,
													},
													"delete": schema.BoolAttribute{
														Description: "Permission to delete roles.",
														Computed:    true,
													},
													"edit": schema.BoolAttribute{
														Description: "Permission to edit roles.",
														Computed:    true,
													},
												},
											},
											"users": schema.SingleNestedAttribute{
												Description: "Permissions for managing users.",
												Computed:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all user permissions are enabled.",
														Computed:    true,
													},
													"create": schema.BoolAttribute{
														Description: "Permission to create users.",
														Computed:    true,
													},
													"delete": schema.BoolAttribute{
														Description: "Permission to delete users.",
														Computed:    true,
													},
													"edit": schema.BoolAttribute{
														Description: "Permission to edit users.",
														Computed:    true,
													},
												},
											},
										},
									},
									"read": schema.SingleNestedAttribute{
										Description: "Permission to read users and roles.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all read permissions are enabled.",
												Computed:    true,
											},
										},
									},
								},
							},
						},
					},
					"user": schema.SingleNestedAttribute{
						Description: "Permissions for managing users.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all user permissions are enabled.",
								Computed:    true,
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing users.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all user management permissions are enabled.",
										Computed:    true,
									},
									"upsert": schema.BoolAttribute{
										Description: "Permission to upsert users.",
										Computed:    true,
									},
								},
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read users.",
								Computed:    true,
							},
						},
					},
					"vulnerability": schema.SingleNestedAttribute{
						Description: "Permissions for managing vulnerabilities.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all vulnerability permissions are enabled.",
								Computed:    true,
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing vulnerabilities.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all vulnerability management permissions are enabled.",
										Computed:    true,
									},
									"ignore": schema.BoolAttribute{
										Description: "Permission to ignore vulnerabilities.",
										Computed:    true,
									},
									"resolve": schema.BoolAttribute{
										Description: "Permission to resolve vulnerabilities.",
										Computed:    true,
									},
									"write": schema.BoolAttribute{
										Description: "Permission to write vulnerabilities.",
										Computed:    true,
									},
								},
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read vulnerabilities.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

// RoleDataSourceModel defines the structure for the role data source model.
type RoleDataSourceModel struct {
	Name        types.String      `tfsdk:"name"`
	Permissions *PermissionsModel `tfsdk:"permissions"`
	ID          types.String      `tfsdk:"role_id"`
	ViprRole    types.Bool        `tfsdk:"vipr_role"`
}

// AllPermissionsModel defines the structure for all permissions.
type AllPermissionsModel struct {
	All types.Bool `tfsdk:"all"`
}

// PermissionsModel defines the structure for permissions.
type PermissionsModel struct {
	AdvancedPermissions *AdvancedPermissionsModel `tfsdk:"advanced_permissions"`
	Alert               *AlertModel               `tfsdk:"alert"`
	Device              *DeviceModel              `tfsdk:"device"`
	Policy              *PolicyModel              `tfsdk:"policy"`
	Report              *ReportModel              `tfsdk:"report"`
	RiskFactor          *RiskFactorModel          `tfsdk:"risk_factor"`
	Settings            *SettingsModel            `tfsdk:"settings"`
	User                *UserModel                `tfsdk:"user"`
	Vulnerability       *VulnerabilityModel       `tfsdk:"vulnerability"`
}

// AdvancedPermissionsModel defines the structure for advanced permissions.
type AdvancedPermissionsModel struct {
	All        types.Bool           `tfsdk:"all"`
	Behavioral *BehavioralModel     `tfsdk:"behavioral"`
	Device     *DeviceAdvancedModel `tfsdk:"device"`
}

// BehavioralModel defines the structure for behavioral permissions.
type BehavioralModel struct {
	All             types.Bool `tfsdk:"all"`
	ApplicationName types.Bool `tfsdk:"application_name"`
	HostName        types.Bool `tfsdk:"host_name"`
	ServiceName     types.Bool `tfsdk:"service_name"`
}

// DeviceAdvancedModel defines the structure for device advanced permissions.
type DeviceAdvancedModel struct {
	All          types.Bool `tfsdk:"all"`
	DeviceNames  types.Bool `tfsdk:"device_names"`
	IPAddresses  types.Bool `tfsdk:"ip_addresses"`
	MACAddresses types.Bool `tfsdk:"mac_addresses"`
	PhoneNumbers types.Bool `tfsdk:"phone_numbers"`
}

// AlertModel defines the structure for alert permissions.
type AlertModel struct {
	All    types.Bool           `tfsdk:"all"`
	Manage *ManageAlertsModel   `tfsdk:"manage"`
	Read   *AllPermissionsModel `tfsdk:"read"`
}

// ManageAlertsModel defines the structure for manage permissions.
type ManageAlertsModel struct {
	All              types.Bool `tfsdk:"all"`
	Resolve          types.Bool `tfsdk:"resolve"`
	Suppress         types.Bool `tfsdk:"suppress"`
	WhitelistDevices types.Bool `tfsdk:"whitelist_devices"`
}

// DeviceModel defines the structure for device permissions.
type DeviceModel struct {
	All    types.Bool           `tfsdk:"all"`
	Manage *ManageDeviceModel   `tfsdk:"manage"`
	Read   *AllPermissionsModel `tfsdk:"read"`
}

// ManageDeviceModel defines the structure for device management permissions.
type ManageDeviceModel struct {
	All                types.Bool    `tfsdk:"all"`
	Create             types.Bool    `tfsdk:"create"`
	Delete             types.Bool    `tfsdk:"delete"`
	Edit               types.Bool    `tfsdk:"edit"`
	Enforce            *EnforceModel `tfsdk:"enforce"`
	Merge              types.Bool    `tfsdk:"merge"`
	RequestDeletedData types.Bool    `tfsdk:"request_deleted_data"`
	Tags               types.Bool    `tfsdk:"tags"`
}

// EnforceModel defines the structure for enforce permissions.
type EnforceModel struct {
	All    types.Bool `tfsdk:"all"`
	Create types.Bool `tfsdk:"create"`
	Delete types.Bool `tfsdk:"delete"`
}

// PolicyModel defines the structure for policy permissions.
type PolicyModel struct {
	All    types.Bool `tfsdk:"all"`
	Manage types.Bool `tfsdk:"manage"`
	Read   types.Bool `tfsdk:"read"`
}

// ReportModel defines the structure for report permissions.
type ReportModel struct {
	All    types.Bool         `tfsdk:"all"`
	Export types.Bool         `tfsdk:"export"`
	Manage *ManageReportModel `tfsdk:"manage"`
	Read   types.Bool         `tfsdk:"read"`
}

// ManageReportModel defines the structure for report management permissions.
type ManageReportModel struct {
	All    types.Bool `tfsdk:"all"`
	Create types.Bool `tfsdk:"create"`
	Delete types.Bool `tfsdk:"delete"`
	Edit   types.Bool `tfsdk:"edit"`
}

// RiskFactorModel defines the structure for risk factor permissions.
type RiskFactorModel struct {
	All    types.Bool       `tfsdk:"all"`
	Manage *ManageRiskModel `tfsdk:"manage"`
	Read   types.Bool       `tfsdk:"read"`
}

// ManageRiskModel defines the structure for risk management permissions.
type ManageRiskModel struct {
	All           types.Bool          `tfsdk:"all"`
	Customization *CustomizationModel `tfsdk:"customization"`
	Status        *StatusModel        `tfsdk:"status"`
}

// CustomizationModel defines the structure for customization permissions.
type CustomizationModel struct {
	All     types.Bool `tfsdk:"all"`
	Create  types.Bool `tfsdk:"create"`
	Disable types.Bool `tfsdk:"disable"`
	Edit    types.Bool `tfsdk:"edit"`
}

// StatusModel defines the structure for status permissions.
type StatusModel struct {
	All     types.Bool `tfsdk:"all"`
	Ignore  types.Bool `tfsdk:"ignore"`
	Resolve types.Bool `tfsdk:"resolve"`
}

// SettingsModel defines the structure for settings permissions.
type SettingsModel struct {
	All              types.Bool            `tfsdk:"all"`
	AuditLog         types.Bool            `tfsdk:"audit_log"`
	Boundary         *BoundaryModel        `tfsdk:"boundary"`
	BusinessImpact   *ManageAndReadModel   `tfsdk:"business_impact"`
	Collector        *ManageAndReadModel   `tfsdk:"collector"`
	CustomProperties *ManageAndReadModel   `tfsdk:"custom_properties"`
	Integration      *ManageAndReadModel   `tfsdk:"integration"`
	InternalIps      *ManageAndReadModel   `tfsdk:"internal_ips"`
	Notifications    *ManageAndReadModel   `tfsdk:"notifications"`
	OIDC             *ManageAndReadModel   `tfsdk:"oidc"`
	SAML             *ManageAndReadModel   `tfsdk:"saml"`
	SecretKey        types.Bool            `tfsdk:"secret_key"`
	SecuritySettings types.Bool            `tfsdk:"security_settings"`
	SitesAndSensors  *SitesAndSensorsModel `tfsdk:"sites_and_sensors"`
	UsersAndRoles    *UsersAndRolesModel   `tfsdk:"users_and_roles"`
}

// BoundaryModel defines the structure for boundary permissions.
type BoundaryModel struct {
	All    types.Bool           `tfsdk:"all"`
	Manage *ManageBoundaryModel `tfsdk:"manage"`
	Read   types.Bool           `tfsdk:"read"`
}

// ManageBoundaryModel defines the structure for managing boundary permissions.
type ManageBoundaryModel struct {
	All    types.Bool `tfsdk:"all"`
	Create types.Bool `tfsdk:"create"`
	Delete types.Bool `tfsdk:"delete"`
	Edit   types.Bool `tfsdk:"edit"`
}

// ManageAndReadModel defines the structure for manage and read permissions.
type ManageAndReadModel struct {
	All    types.Bool `tfsdk:"all"`
	Manage types.Bool `tfsdk:"manage"`
	Read   types.Bool `tfsdk:"read"`
}

// SitesAndSensorsModel defines the structure for sites and sensors permissions.
type SitesAndSensorsModel struct {
	All    types.Bool                  `tfsdk:"all"`
	Manage *ManageSitesAndSensorsModel `tfsdk:"manage"`
	Read   *AllPermissionsModel        `tfsdk:"read"`
}

// ManageSitesAndSensorsModel defines the structure for managing sites and sensors permissions.
type ManageSitesAndSensorsModel struct {
	All     types.Bool `tfsdk:"all"`
	Sensors types.Bool `tfsdk:"sensors"`
	Sites   types.Bool `tfsdk:"sites"`
}

// UsersAndRolesModel defines the structure for users and roles permissions.
type UsersAndRolesModel struct {
	All    types.Bool                `tfsdk:"all"`
	Manage *ManageUsersAndRolesModel `tfsdk:"manage"`
	Read   types.Bool                `tfsdk:"read"`
}

// ManageUsersAndRolesModel defines the structure for managing users and roles permissions.
type ManageUsersAndRolesModel struct {
	All   types.Bool        `tfsdk:"all"`
	Roles *ManageRolesModel `tfsdk:"roles"`
	Users *ManageUsersModel `tfsdk:"users"`
}

// ManageRolesModel defines the structure for managing roles permissions.
type ManageRolesModel struct {
	All    types.Bool `tfsdk:"all"`
	Create types.Bool `tfsdk:"create"`
	Delete types.Bool `tfsdk:"delete"`
	Edit   types.Bool `tfsdk:"edit"`
}

// ManageUsersModel defines the structure for managing users permissions.
type ManageUsersModel struct {
	All    types.Bool `tfsdk:"all"`
	Create types.Bool `tfsdk:"create"`
	Delete types.Bool `tfsdk:"delete"`
	Edit   types.Bool `tfsdk:"edit"`
}

// UserModel defines the structure for user permissions.
type UserModel struct {
	All    types.Bool       `tfsdk:"all"`
	Manage *ManageUserModel `tfsdk:"manage"`
	Read   types.Bool       `tfsdk:"read"`
}

// ManageUserModel defines the structure for managing user permissions.
type ManageUserModel struct {
	All    types.Bool `tfsdk:"all"`
	Upsert types.Bool `tfsdk:"upsert"`
}

// VulnerabilityModel defines the structure for vulnerability permissions.
type VulnerabilityModel struct {
	All    types.Bool                `tfsdk:"all"`
	Manage *ManageVulnerabilityModel `tfsdk:"manage"`
	Read   types.Bool                `tfsdk:"read"`
}

// ManageVulnerabilityModel defines the structure for managing vulnerabilities permissions.
type ManageVulnerabilityModel struct {
	All     types.Bool `tfsdk:"all"`
	Ignore  types.Bool `tfsdk:"ignore"`
	Resolve types.Bool `tfsdk:"resolve"`
	Write   types.Bool `tfsdk:"write"`
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

	role, err := d.client.GetRoleByName(ctx, config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Armis Role",
			fmt.Sprintf("Error fetching role: %s", err.Error()),
		)
		return
	}

	roleState := BuildRoleDataSourceModel(role)

	// Set the state with the fetched role
	resp.Diagnostics.Append(resp.State.Set(ctx, &roleState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func BuildRoleDataSourceModel(role *armis.RoleSettings) RoleDataSourceModel {
	return RoleDataSourceModel{
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
				Manage: &ManageAlertsModel{
					All:              types.BoolValue(role.Permissions.Alert.Manage.All),
					Resolve:          types.BoolValue(role.Permissions.Alert.Manage.Resolve.All),
					Suppress:         types.BoolValue(role.Permissions.Alert.Manage.Suppress.All),
					WhitelistDevices: types.BoolValue(role.Permissions.Alert.Manage.WhitelistDevices.All),
				},
				Read: &AllPermissionsModel{
					All: types.BoolValue(role.Permissions.Alert.Read.All),
				},
			},
			Policy: &PolicyModel{
				All:    types.BoolValue(role.Permissions.Policy.All),
				Manage: types.BoolValue(role.Permissions.Policy.Manage.All),
				Read:   types.BoolValue(role.Permissions.Policy.Read.All),
			},
			Report: &ReportModel{
				All:    types.BoolValue(role.Permissions.Report.All),
				Export: types.BoolValue(role.Permissions.Report.Export.All),
				Manage: &ManageReportModel{
					All:    types.BoolValue(role.Permissions.Report.All),
					Create: types.BoolValue(role.Permissions.Report.Manage.Create.All),
					Delete: types.BoolValue(role.Permissions.Report.Manage.Delete.All),
					Edit:   types.BoolValue(role.Permissions.Report.Manage.Edit.All),
				},
				Read: types.BoolValue(role.Permissions.Report.Read.All),
			},
			RiskFactor: &RiskFactorModel{
				All: types.BoolValue(role.Permissions.RiskFactor.All),
				Manage: &ManageRiskModel{
					All: types.BoolValue(role.Permissions.RiskFactor.Manage.All),
					Customization: &CustomizationModel{
						All:     types.BoolValue(role.Permissions.RiskFactor.Manage.Customization.All),
						Create:  types.BoolValue(role.Permissions.RiskFactor.Manage.Customization.Create.All),
						Disable: types.BoolValue(role.Permissions.RiskFactor.Manage.Customization.Disable.All),
						Edit:    types.BoolValue(role.Permissions.RiskFactor.Manage.Customization.Edit.All),
					},
					Status: &StatusModel{
						All:     types.BoolValue(role.Permissions.RiskFactor.Manage.Status.All),
						Ignore:  types.BoolValue(role.Permissions.RiskFactor.Manage.Status.Ignore.All),
						Resolve: types.BoolValue(role.Permissions.RiskFactor.Manage.Status.Resolve.All),
					},
				},
			},
			Settings: &SettingsModel{
				All:      types.BoolValue(role.Permissions.Settings.All),
				AuditLog: types.BoolValue(role.Permissions.Settings.AuditLog.All),
				Boundary: &BoundaryModel{
					All: types.BoolValue(role.Permissions.Settings.Boundary.All),
					Manage: &ManageBoundaryModel{
						All:    types.BoolValue(role.Permissions.Settings.Boundary.Manage.All),
						Create: types.BoolValue(role.Permissions.Settings.Boundary.Manage.Create.All),
						Delete: types.BoolValue(role.Permissions.Settings.Boundary.Manage.Delete.All),
						Edit:   types.BoolValue(role.Permissions.Settings.Boundary.Manage.Edit.All),
					},
					Read: types.BoolValue(role.Permissions.Settings.Boundary.Read.All),
				},
				BusinessImpact: &ManageAndReadModel{
					All:    types.BoolValue(role.Permissions.Settings.BusinessImpact.All),
					Manage: types.BoolValue(role.Permissions.Settings.BusinessImpact.Manage.All),
					Read:   types.BoolValue(role.Permissions.Settings.BusinessImpact.Read.All),
				},
				Collector: &ManageAndReadModel{
					All:    types.BoolValue(role.Permissions.Settings.Collector.All),
					Manage: types.BoolValue(role.Permissions.Settings.Collector.Manage.All),
					Read:   types.BoolValue(role.Permissions.Settings.Collector.Read.All),
				},
				CustomProperties: &ManageAndReadModel{
					All:    types.BoolValue(role.Permissions.Settings.CustomProperties.All),
					Manage: types.BoolValue(role.Permissions.Settings.CustomProperties.Manage.All),
					Read:   types.BoolValue(role.Permissions.Settings.CustomProperties.Read.All),
				},
				Integration: &ManageAndReadModel{
					All:    types.BoolValue(role.Permissions.Settings.Integration.All),
					Manage: types.BoolValue(role.Permissions.Settings.Integration.Manage.All),
					Read:   types.BoolValue(role.Permissions.Settings.Integration.Read.All),
				},
				InternalIps: &ManageAndReadModel{
					All:    types.BoolValue(role.Permissions.Settings.InternalIps.All),
					Manage: types.BoolValue(role.Permissions.Settings.InternalIps.Manage.All),
					Read:   types.BoolValue(role.Permissions.Settings.InternalIps.Read.All),
				},
				Notifications: &ManageAndReadModel{
					All:    types.BoolValue(role.Permissions.Settings.Notifications.All),
					Manage: types.BoolValue(role.Permissions.Settings.Notifications.Manage.All),
					Read:   types.BoolValue(role.Permissions.Settings.Notifications.Read.All),
				},
				OIDC: &ManageAndReadModel{
					All:    types.BoolValue(role.Permissions.Settings.OIDC.All),
					Manage: types.BoolValue(role.Permissions.Settings.OIDC.Manage.All),
					Read:   types.BoolValue(role.Permissions.Settings.OIDC.Read.All),
				},
				SAML: &ManageAndReadModel{
					All:    types.BoolValue(role.Permissions.Settings.SAML.All),
					Manage: types.BoolValue(role.Permissions.Settings.SAML.Manage.All),
					Read:   types.BoolValue(role.Permissions.Settings.SAML.Read.All),
				},
				SecretKey:        types.BoolValue(role.Permissions.Settings.SecretKey.All),
				SecuritySettings: types.BoolValue(role.Permissions.Settings.SecuritySettings.All),
				SitesAndSensors: &SitesAndSensorsModel{
					All: types.BoolValue(role.Permissions.Settings.SitesAndSensors.All),
					Manage: &ManageSitesAndSensorsModel{
						All:     types.BoolValue(role.Permissions.Settings.SitesAndSensors.Manage.All),
						Sensors: types.BoolValue(role.Permissions.Settings.SitesAndSensors.Manage.Sensors.All),
						Sites:   types.BoolValue(role.Permissions.Settings.SitesAndSensors.Manage.Sites.All),
					},
					Read: &AllPermissionsModel{
						All: types.BoolValue(role.Permissions.Settings.SitesAndSensors.Read.All),
					},
				},
				UsersAndRoles: &UsersAndRolesModel{
					All: types.BoolValue(role.Permissions.Settings.UsersAndRoles.All),
					Manage: &ManageUsersAndRolesModel{
						All: types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.All),
						Roles: &ManageRolesModel{
							All:    types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Roles.All),
							Create: types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Roles.Create.All),
							Delete: types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Roles.Delete.All),
							Edit:   types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Roles.Edit.All),
						},
						Users: &ManageUsersModel{
							All:    types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Users.All),
							Create: types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Users.Create.All),
							Delete: types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Users.Delete.All),
							Edit:   types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Users.Edit.All),
						},
					},
					Read: types.BoolValue(role.Permissions.Settings.UsersAndRoles.Read.All),
				},
			},
			User: &UserModel{
				All: types.BoolValue(role.Permissions.User.All),
				Manage: &ManageUserModel{
					All:    types.BoolValue(role.Permissions.User.Manage.All),
					Upsert: types.BoolValue(role.Permissions.User.Manage.Upsert.All),
				},
				Read: types.BoolValue(role.Permissions.User.Read.All),
			},
			Vulnerability: &VulnerabilityModel{
				All: types.BoolValue(role.Permissions.Vulnerability.All),
				Manage: &ManageVulnerabilityModel{
					All:     types.BoolValue(role.Permissions.Vulnerability.Manage.All),
					Ignore:  types.BoolValue(role.Permissions.Vulnerability.Manage.Ignore.All),
					Resolve: types.BoolValue(role.Permissions.Vulnerability.Manage.Resolve.All),
					Write:   types.BoolValue(role.Permissions.Vulnerability.Manage.Write.All),
				},
				Read: types.BoolValue(role.Permissions.Vulnerability.Read.All),
			},
		},
	}
}
