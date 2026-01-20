// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/1898andCo/armis-sdk-go/armis"
	u "github.com/1898andCo/terraform-provider-armis-centrix/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	resp.TypeName = req.ProviderTypeName + "_roles"
}

// Schema defines the schema for the roles data source.
func (d *rolesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an Armis role",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the role. When omitted, roles can be filtered using prefix options.",
			},
			"role_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the role.",
			},
			"match_prefix": schema.StringAttribute{
				Optional:    true,
				Description: "Optional prefix to match role names. When provided, only roles with names starting with this prefix are returned.",
			},
			"exclude_prefix": schema.StringAttribute{
				Optional:    true,
				Description: "Optional prefix to exclude role names. When provided, roles with names starting with this prefix are not returned.",
			},
			"roles": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A computed list of Armis roles matching the supplied filters.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role_id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for the role.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the role.",
						},
						"vipr_role": schema.BoolAttribute{
							Computed:    true,
							Description: "Indicates if the role is a VIPR-specific role.",
						},
					},
				},
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
									"read": schema.BoolAttribute{
										Description: "Permission to read sites and sensors.",
										Computed:    true,
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
									"read": schema.BoolAttribute{
										Description: "Permission to read users and roles.",
										Computed:    true,
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

// Read refreshes the Terraform state with the latest data.
func (d *rolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config u.RoleDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var roles []u.RoleDataSourceSummaryModel

	if !config.Name.IsNull() && config.Name.ValueString() != "" {
		matchPrefix := config.MatchPrefix
		excludePrefix := config.ExcludePrefix
		role, err := d.client.GetRoleByName(ctx, config.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Armis Role",
				err.Error(),
			)
			return
		}

		summary := u.BuildRoleDataSourceSummaryModel(role)
		if u.ShouldIncludeRole(summary, matchPrefix) && !u.ShouldExcludeRole(summary, excludePrefix) {
			roles = append(roles, summary)
			config = u.BuildRoleDataSourceModel(role)
			config.MatchPrefix = matchPrefix
			config.ExcludePrefix = excludePrefix
		}
	} else {
		allRoles, err := d.client.GetRoles(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Armis Roles",
				err.Error(),
			)
			return
		}

		for idx := range allRoles {
			summary := u.BuildRoleDataSourceSummaryModel(&allRoles[idx])
			if u.ShouldIncludeRole(summary, config.MatchPrefix) && !u.ShouldExcludeRole(summary, config.ExcludePrefix) {
				roles = append(roles, summary)
			}
		}
	}

	config.Roles = roles

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
