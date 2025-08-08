// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	armis "github.com/1898andCo/terraform-provider-armis-centrix/internal/armis"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &roleResource{}
	_ resource.ResourceWithConfigure = &roleResource{}
)

type roleResource struct {
	client *armis.Client
}

func RoleResource() resource.Resource {
	return &roleResource{}
}

// Configure adds the provider configured client to the resource.
func (r *roleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Metadata returns the resource type name.
func (r *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema defines the schema for the role resource.
func (r *roleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
		Provides an Armis Role.

		This resource configures permissions, including advanced and alert permissions, for a role.
		The nested permissions follow a parent-child Boolean relationship:
		- If a parent option is True, all its nested options must also be True.
		- If any nested option is False, the parent option cannot be True.
		`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:      true,
				Description:   "The name of the role.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "Unique identifier for the role.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"permissions": schema.SingleNestedAttribute{
				Description: "Permissions associated with the role, categorized by feature and capability.",
				Required:    true, // Permissions must be defined for every role.
				Attributes: map[string]schema.Attribute{
					"advanced_permissions": schema.SingleNestedAttribute{
						Description: "Advanced permissions for managing sensitive data and configurations.",
						Optional:    true,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all advanced permissions are enabled.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"behavioral": schema.SingleNestedAttribute{
								Description: "Behavioral permissions for specific system entities.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all behavioral permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"application_name": schema.SingleNestedAttribute{
										Description: "Permission to access application names.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all application name permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"host_name": schema.SingleNestedAttribute{
										Description: "Permission to access host names.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all host name permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"service_name": schema.SingleNestedAttribute{
										Description: "Permission to access service names.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all service name permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
								},
							},
							"device": schema.SingleNestedAttribute{
								Description: "Permissions for managing device-related data.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all device-related permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"device_names": schema.SingleNestedAttribute{
										Description: "Permission to access device names.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all device name permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"ip_addresses": schema.SingleNestedAttribute{
										Description: "Permission to access device IP addresses.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all IP address permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"mac_addresses": schema.SingleNestedAttribute{
										Description: "Permission to access device MAC addresses.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all MAC address permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"phone_numbers": schema.SingleNestedAttribute{
										Description: "Permission to access device phone numbers.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all phone number permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
								},
							},
						},
					},
					"alert": schema.SingleNestedAttribute{
						Description: "Permissions for managing alerts.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all alert permissions are enabled.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing alerts.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all alert management permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"resolve": schema.SingleNestedAttribute{
										Description: "Permission to resolve alerts.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all alert resolve permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"suppress": schema.SingleNestedAttribute{
										Description: "Permission to suppress alerts.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all alert suppress permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"whitelist_devices": schema.SingleNestedAttribute{
										Description: "Permission to whitelist devices in alerts.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all whitelist device permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
								},
							},
							"read": schema.SingleNestedAttribute{
								Description: "Permission to read alerts.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all alert read permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
								},
							},
						},
					},
					"device": schema.SingleNestedAttribute{
						Description: "Permissions for managing devices.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all device permissions are enabled.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing devices.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all device management permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"create": schema.SingleNestedAttribute{
										Description: "Permission to create devices.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all device create permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"delete": schema.SingleNestedAttribute{
										Description: "Permission to delete devices.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all device delete permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"edit": schema.SingleNestedAttribute{
										Description: "Permission to edit devices.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all device edit permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"enforce": schema.SingleNestedAttribute{
										Description: "Permission to enforce device policies.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all device enforcement permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
											"create": schema.SingleNestedAttribute{
												Description: "Permission to create device enforcement rules.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all device enforcement create permissions are enabled.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
												},
											},
											"delete": schema.SingleNestedAttribute{
												Description: "Permission to delete device enforcement rules.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all device enforcement delete permissions are enabled.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
												},
											},
										},
									},
									"merge": schema.SingleNestedAttribute{
										Description: "Permission to merge devices.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all device merge permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"request_deleted_data": schema.SingleNestedAttribute{
										Description: "Permission to request deleted device data.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all request deleted data permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"tags": schema.SingleNestedAttribute{
										Description: "Permission to manage device tags.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all device tag permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
								},
							},
							"read": schema.SingleNestedAttribute{
								Description: "Permission to read device information.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all device read permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
								},
							},
						},
					},
					"policy": schema.SingleNestedAttribute{
						Description: "Permissions for managing policies.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all policy permissions are enabled.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permission to manage policies.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all policy management permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
								},
							},
							"read": schema.SingleNestedAttribute{
								Description: "Permission to read policies.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all policy read permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
								},
							},
						},
					},
					"report": schema.SingleNestedAttribute{
						Description: "Permissions for managing reports.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all report permissions are enabled.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"export": schema.SingleNestedAttribute{
								Description: "Permission to export reports.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all report export permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
								},
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing reports.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all report management permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"create": schema.SingleNestedAttribute{
										Description: "Permission to create reports.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all report create permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"delete": schema.SingleNestedAttribute{
										Description: "Permission to delete reports.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all report delete permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"edit": schema.SingleNestedAttribute{
										Description: "Permission to edit reports.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all report edit permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
								},
							},
							"read": schema.SingleNestedAttribute{
								Description: "Permission to read reports.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all report read permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
								},
							},
						},
					},
					"risk_factor": schema.SingleNestedAttribute{
						Description: "Permissions for managing risk factors.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all risk factor permissions are enabled.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing risk factors.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all risk factor management permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"customization": schema.SingleNestedAttribute{
										Description: "Permission to customize risk factors.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all risk factor customization permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
											"create": schema.SingleNestedAttribute{
												Description: "Permission to create risk factor customizations.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all risk factor customization create permissions are enabled.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
												},
											},
											"disable": schema.SingleNestedAttribute{
												Description: "Permission to disable risk factor customizations.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all risk factor customization disable permissions are enabled.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
												},
											},
											"edit": schema.SingleNestedAttribute{
												Description: "Permission to edit risk factor customizations.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all risk factor customization edit permissions are enabled.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
												},
											},
										},
									},
									"status": schema.SingleNestedAttribute{
										Description: "Permission to manage risk factor status.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all risk factor status permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
											"ignore": schema.SingleNestedAttribute{
												Description: "Permission to ignore risk factors.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all risk factor ignore permissions are enabled.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
												},
											},
											"resolve": schema.SingleNestedAttribute{
												Description: "Permission to resolve risk factors.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all risk factor resolve permissions are enabled.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
												},
											},
										},
									},
								},
							},
							"read": schema.SingleNestedAttribute{
								Description: "Permission to read risk factors.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all risk factor read permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
								},
							},
						},
					},
					"settings": schema.SingleNestedAttribute{
						Description: "Permissions for managing system settings.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all settings permissions are enabled.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"audit_log": schema.SingleNestedAttribute{
								Description: "Permission to access audit logs.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all audit log permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
								},
							},
							"boundary": schema.SingleNestedAttribute{
								Description: "Permissions for managing boundaries.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all boundary permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"manage": schema.SingleNestedAttribute{
										Description: "Permissions for managing boundaries.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all boundary management permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
											"create": schema.SingleNestedAttribute{
												Description: "Permission to create boundaries.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all boundary create permissions are enabled.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
												},
											},
											"delete": schema.SingleNestedAttribute{
												Description: "Permission to delete boundaries.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all boundary delete permissions are enabled.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
												},
											},
											"edit": schema.SingleNestedAttribute{
												Description: "Permission to edit boundaries.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all boundary edit permissions are enabled.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
												},
											},
										},
									},
									"read": schema.SingleNestedAttribute{
										Description: "Permission to read boundaries.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all boundary read permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
								},
							},
							"business_impact": schema.SingleNestedAttribute{
								Description: "Permissions for managing business impact settings.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all business impact permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"manage": schema.SingleNestedAttribute{
										Description: "Permission to manage business impact settings.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all business impact management permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"read": schema.SingleNestedAttribute{
										Description: "Permission to read business impact settings.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all business impact read permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
								},
							},
							"collector": schema.SingleNestedAttribute{
								Description: "Permissions for managing collectors.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all collector permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"manage": schema.SingleNestedAttribute{
										Description: "Permission to manage collectors.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all collector management permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"read": schema.SingleNestedAttribute{
										Description: "Permission to read collector information.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all collector read permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
								},
							},
							"custom_properties": schema.SingleNestedAttribute{
								Description: "Permissions for managing custom properties.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all custom properties permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"manage": schema.SingleNestedAttribute{
										Description: "Permission to manage custom properties.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all custom properties management permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"read": schema.SingleNestedAttribute{
										Description: "Permission to read custom properties.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all custom properties read permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
								},
							},
							"integration": schema.SingleNestedAttribute{
								Description: "Permissions for managing integrations.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all integration permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"manage": schema.SingleNestedAttribute{
										Description: "Permission to manage integrations.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all integration management permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"read": schema.SingleNestedAttribute{
										Description: "Permission to read integration information.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all integration read permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
								},
							},
							"internal_ips": schema.SingleNestedAttribute{
								Description: "Permissions for managing internal IP addresses.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all internal IP permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"manage": schema.SingleNestedAttribute{
										Description: "Permission to manage internal IP addresses.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all internal IP management permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"read": schema.SingleNestedAttribute{
										Description: "Permission to read internal IP addresses.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all internal IP read permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
								},
							},
							"notifications": schema.SingleNestedAttribute{
								Description: "Permissions for managing notifications.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all notification permissions are enabled.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"manage": schema.SingleNestedAttribute{
										Description: "Permission to manage notifications.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all notification management permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
									"read": schema.SingleNestedAttribute{
										Description: "Permission to read notification settings.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all notification read permissions are enabled.",
												Optional:    true,
												Computed:    true,
												Default:     booldefault.StaticBool(false),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// RoleResourceModel maps the RoleSettings schema data. The full PermissionsModel lives in role_data_source.go.
type RoleResourceModel struct {
	Name        types.String      `tfsdk:"name"`
	Permissions *PermissionsModel `tfsdk:"permissions"`
	ID          types.String      `tfsdk:"id"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RoleResourceModel
	tflog.Info(ctx, "Starting role creation")

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating role with provided plan", map[string]any{"name": plan.Name.ValueString()})

	if plan.Permissions == nil {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Permissions block is required but not provided.",
		)
		return
	}

	if plan.Permissions.AdvancedPermissions == nil {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Advanced permissions block is required but not provided.",
		)
		return
	}

	role := mapPlanToRoleSettings(plan)

	// Call API to create the role
	success, err := r.client.CreateRole(ctx, role)
	if err != nil || !success {
		resp.Diagnostics.AddError(
			"Error creating role",
			fmt.Sprintf("Failed to create role %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	// Fetch the created role to get its ID and other attributes
	createdRole, err := r.client.GetRoleByName(ctx, plan.Name.ValueString())
	if err != nil || createdRole == nil {
		resp.Diagnostics.AddError(
			"Error fetching created role",
			fmt.Sprintf("Role %q was created but could not fetch details: %s", plan.Name.ValueString(), err),
		)
		return
	}

	// Update the Terraform state with the created role's details
	plan.ID = types.StringValue(strconv.Itoa(createdRole.ID))

	tflog.Info(ctx, "Setting Terraform state for created role", map[string]any{"role_id": plan.ID.ValueString()})
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read reads the role's current state.
func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoleResourceModel
	tflog.Info(ctx, "Reading role state")

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the role by ID
	tflog.Debug(ctx, "Fetching role by ID", map[string]any{"role_id": state.ID.ValueString()})
	role, err := r.client.GetRoleByID(ctx, state.ID.ValueString())
	if err != nil || role == nil {
		resp.Diagnostics.AddError(
			"Error reading role",
			fmt.Sprintf("Failed to fetch role with ID %q: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Update the state with refreshed role details
	roleState := BuildRoleResourceModel(role)

	tflog.Debug(ctx, "Setting refreshed state for role", map[string]any{"role_id": state.ID.ValueString()})
	diags = resp.State.Set(ctx, &roleState)
	resp.Diagnostics.Append(diags...)
}

// Update updates the role.
func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating role")

	// Retrieve values from the plan
	var plan RoleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the current state to get the role ID
	var state RoleResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that the role ID is available
	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Error Updating Role",
			"The role ID is missing from the state. This is required to update the role.",
		)
		return
	}

	// Map the plan to role settings for the update
	role := mapPlanToRoleSettings(plan)

	// Update the role in the API
	tflog.Debug(ctx, "Sending update request to Armis API", map[string]any{"role_id": state.ID.ValueString()})
	_, err := r.client.UpdateRole(ctx, role, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Role",
			fmt.Sprintf("Failed to update role with ID %q: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Fetch the updated role details
	updatedRole, err := r.client.GetRoleByID(ctx, state.ID.ValueString())
	if err != nil || updatedRole == nil {
		resp.Diagnostics.AddError(
			"Error Fetching Updated Role",
			fmt.Sprintf("The role with ID %q was updated, but its details could not be fetched: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Map the updated role details back to the plan
	plan = BuildRoleResourceModel(updatedRole)

	// Save the updated state
	tflog.Info(ctx, "Setting updated state for role", map[string]any{"role_id": state.ID.ValueString()})
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the role.
func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RoleResourceModel
	tflog.Info(ctx, "Deleting role")

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the role in the API
	tflog.Debug(ctx, "Deleting role in Armis", map[string]any{"role_id": state.ID.ValueString()})
	success, err := r.client.DeleteRole(ctx, state.ID.ValueString())
	if err != nil || !success {
		resp.Diagnostics.AddError(
			"Error deleting role",
			fmt.Sprintf("Failed to delete role %q: %s", state.ID.ValueString(), err),
		)
		return
	}

	tflog.Info(ctx, "Role deleted successfully", map[string]any{"role_id": state.ID.ValueString()})
}

// mapPlanToRoleSettings converts a RoleResourceModel to an Armis RoleSettings object.
func mapPlanToRoleSettings(plan RoleResourceModel) armis.RoleSettings {
	var permissions armis.Permissions

	if plan.Permissions != nil {
		if plan.Permissions.AdvancedPermissions != nil {
			permissions.AdvancedPermissions = armis.AdvancedPermissions{
				All: plan.Permissions.AdvancedPermissions.All.ValueBool(),
			}

			if plan.Permissions.AdvancedPermissions.Behavioral != nil {
				permissions.AdvancedPermissions.Behavioral = armis.Behavioral{
					All: plan.Permissions.AdvancedPermissions.Behavioral.All.ValueBool(),
					ApplicationName: armis.Permission{
						All: plan.Permissions.AdvancedPermissions.Behavioral.ApplicationName.ValueBool(),
					},
					HostName: armis.Permission{
						All: plan.Permissions.AdvancedPermissions.Behavioral.HostName.ValueBool(),
					},
					ServiceName: armis.Permission{
						All: plan.Permissions.AdvancedPermissions.Behavioral.ServiceName.ValueBool(),
					},
				}
			}
		}

		if plan.Permissions.Alert != nil {
			permissions.Alert = armis.Alert{
				All: plan.Permissions.Alert.All.ValueBool(),
			}
		}
	}

	return armis.RoleSettings{
		Name:        plan.Name.ValueString(),
		Permissions: permissions,
	}
}

func BuildRoleResourceModel(role *armis.RoleSettings) RoleResourceModel {
	return RoleResourceModel{
		ID:   types.StringValue(fmt.Sprintf("%d", role.ID)),
		Name: types.StringValue(role.Name),
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
				Read: types.BoolValue(role.Permissions.Alert.Read.All),
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
					All:    types.BoolValue(role.Permissions.Report.Manage.All),
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
					Read: types.BoolValue(role.Permissions.Settings.SitesAndSensors.Read.All),
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
