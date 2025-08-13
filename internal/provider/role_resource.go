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

func (r *roleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Armis role",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the role.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the role.",
			},
			"permissions": schema.SingleNestedAttribute{
				Description: "Permissions associated with the role.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"advanced_permissions": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Advanced permissions for the role.",
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Optional:    true,
								Description: "Indicates if the role has all advanced permissions.",
							},
							"behavioral": schema.SingleNestedAttribute{
								Optional:    true,
								Description: "Behavioral permissions for the role.",
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Optional:    true,
										Description: "Indicates if the role has all behavioral permissions.",
									},
									"application_name": schema.BoolAttribute{
										Optional:    true,
										Description: "Permission for application names.",
									},
									"host_name": schema.BoolAttribute{
										Optional:    true,
										Description: "Permission for host names.",
									},
									"service_name": schema.BoolAttribute{
										Optional:    true,
										Description: "Permission for service names.",
									},
								},
							},
							"device": schema.SingleNestedAttribute{
								Optional:    true,
								Description: "Device-related permissions.",
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Optional:    true,
										Description: "Indicates if the role has all device permissions.",
									},
									"device_names": schema.BoolAttribute{
										Optional:    true,
										Description: "Permission for device names.",
									},
									"ip_addresses": schema.BoolAttribute{
										Optional:    true,
										Description: "Permission for IP addresses.",
									},
									"mac_addresses": schema.BoolAttribute{
										Optional:    true,
										Description: "Permission for MAC addresses.",
									},
									"phone_numbers": schema.BoolAttribute{
										Optional:    true,
										Description: "Permission for phone numbers.",
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
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing alerts.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all alert management permissions are enabled.",
										Optional:    true,
									},
									"resolve": schema.BoolAttribute{
										Description: "Permission to resolve alerts.",
										Optional:    true,
									},
									"suppress": schema.BoolAttribute{
										Description: "Permission to suppress alerts.",
										Optional:    true,
									},
									"whitelist_devices": schema.BoolAttribute{
										Description: "Permission to whitelist devices in alerts.",
										Optional:    true,
									},
								},
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read alerts.",
								Optional:    true,
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
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing devices.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all device management permissions are enabled.",
										Optional:    true,
									},
									"create": schema.BoolAttribute{
										Description: "Permission to create devices.",
										Optional:    true,
									},
									"delete": schema.BoolAttribute{
										Description: "Permission to delete devices.",
										Optional:    true,
									},
									"edit": schema.BoolAttribute{
										Description: "Permission to edit devices.",
										Optional:    true,
									},
									"enforce": schema.SingleNestedAttribute{
										Description: "Permissions for enforcing device policies.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all enforce permissions are enabled.",
												Optional:    true,
											},
											"create": schema.BoolAttribute{
												Description: "Permission to create enforcement policies.",
												Optional:    true,
											},
											"delete": schema.BoolAttribute{
												Description: "Permission to delete enforcement policies.",
												Optional:    true,
											},
										},
									},
									"merge": schema.BoolAttribute{
										Description: "Permission to merge devices.",
										Optional:    true,
									},
									"request_deleted_data": schema.BoolAttribute{
										Description: "Permission to request deleted data.",
										Optional:    true,
									},
									"tags": schema.BoolAttribute{
										Description: "Permission to manage device tags.",
										Optional:    true,
									},
								},
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read devices.",
								Optional:    true,
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
							},
							"manage": schema.BoolAttribute{
								Description: "Permission to manage policies.",
								Optional:    true,
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read policies.",
								Optional:    true,
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
							},
							"export": schema.BoolAttribute{
								Description: "Permission to export reports.",
								Optional:    true,
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing reports.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all report management permissions are enabled.",
										Optional:    true,
									},
									"create": schema.BoolAttribute{
										Description: "Permission to create reports.",
										Optional:    true,
									},
									"delete": schema.BoolAttribute{
										Description: "Permission to delete reports.",
										Optional:    true,
									},
									"edit": schema.BoolAttribute{
										Description: "Permission to edit reports.",
										Optional:    true,
									},
								},
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read reports.",
								Optional:    true,
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
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing risk factors.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all risk factor management permissions are enabled.",
										Optional:    true,
									},
									"customization": schema.SingleNestedAttribute{
										Description: "Permissions for customizing risk factors.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all customization permissions are enabled.",
												Optional:    true,
											},
											"create": schema.BoolAttribute{
												Description: "Permission to create customizations.",
												Optional:    true,
											},
											"disable": schema.BoolAttribute{
												Description: "Permission to disable customizations.",
												Optional:    true,
											},
											"edit": schema.BoolAttribute{
												Description: "Permission to edit customizations.",
												Optional:    true,
											},
										},
									},
									"status": schema.SingleNestedAttribute{
										Description: "Permissions for managing risk factor status.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all status permissions are enabled.",
												Optional:    true,
											},
											"ignore": schema.BoolAttribute{
												Description: "Permission to ignore risk factors.",
												Optional:    true,
											},
											"resolve": schema.BoolAttribute{
												Description: "Permission to resolve risk factors.",
												Optional:    true,
											},
										},
									},
								},
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read risk factors.",
								Optional:    true,
							},
						},
					},
					"settings": schema.SingleNestedAttribute{
						Description: "Permissions for managing settings.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all settings permissions are enabled.",
								Optional:    true,
							},
							"audit_log": schema.BoolAttribute{
								Description: "Permission to access audit logs.",
								Optional:    true,
							},
							"boundary": schema.SingleNestedAttribute{
								Description: "Permissions for managing boundaries.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all boundary permissions are enabled.",
										Optional:    true,
									},
									"manage": schema.SingleNestedAttribute{
										Description: "Permissions for managing boundaries.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all boundary management permissions are enabled.",
												Optional:    true,
											},
											"create": schema.BoolAttribute{
												Description: "Permission to create boundaries.",
												Optional:    true,
											},
											"delete": schema.BoolAttribute{
												Description: "Permission to delete boundaries.",
												Optional:    true,
											},
											"edit": schema.BoolAttribute{
												Description: "Permission to edit boundaries.",
												Optional:    true,
											},
										},
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read boundaries.",
										Optional:    true,
									},
								},
							},
							"business_impact": schema.SingleNestedAttribute{
								Description: "Permissions for managing business impact.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all business impact permissions are enabled.",
										Optional:    true,
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage business impact.",
										Optional:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read business impact.",
										Optional:    true,
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
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage collectors.",
										Optional:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read collectors.",
										Optional:    true,
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
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage custom properties.",
										Optional:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read custom properties.",
										Optional:    true,
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
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage integrations.",
										Optional:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read integrations.",
										Optional:    true,
									},
								},
							},
							"internal_ips": schema.SingleNestedAttribute{
								Description: "Permissions for managing internal IPs.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all internal IPs permissions are enabled.",
										Optional:    true,
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage internal IPs.",
										Optional:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read internal IPs.",
										Optional:    true,
									},
								},
							},
							"notifications": schema.SingleNestedAttribute{
								Description: "Permissions for managing notifications.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all notifications permissions are enabled.",
										Optional:    true,
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage notifications.",
										Optional:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read notifications.",
										Optional:    true,
									},
								},
							},
							"oidc": schema.SingleNestedAttribute{
								Description: "Permissions for managing OIDC.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all OIDC permissions are enabled.",
										Optional:    true,
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage OIDC.",
										Optional:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read OIDC.",
										Optional:    true,
									},
								},
							},
							"saml": schema.SingleNestedAttribute{
								Description: "Permissions for managing SAML.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all SAML permissions are enabled.",
										Optional:    true,
									},
									"manage": schema.BoolAttribute{
										Description: "Permission to manage SAML.",
										Optional:    true,
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read SAML.",
										Optional:    true,
									},
								},
							},
							"secret_key": schema.BoolAttribute{
								Description: "Permission to access secret keys.",
								Optional:    true,
							},
							"security_settings": schema.BoolAttribute{
								Description: "Permission to access security settings.",
								Optional:    true,
							},
							"sites_and_sensors": schema.SingleNestedAttribute{
								Description: "Permissions for managing sites and sensors.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all sites and sensors permissions are enabled.",
										Optional:    true,
									},
									"manage": schema.SingleNestedAttribute{
										Description: "Permissions for managing sites and sensors.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all manage permissions are enabled.",
												Optional:    true,
											},
											"sensors": schema.BoolAttribute{
												Description: "Permission to manage sensors.",
												Optional:    true,
											},
											"sites": schema.BoolAttribute{
												Description: "Permission to manage sites.",
												Optional:    true,
											},
										},
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read sites and sensors.",
										Optional:    true,
									},
								},
							},
							"users_and_roles": schema.SingleNestedAttribute{
								Description: "Permissions for managing users and roles.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all users and roles permissions are enabled.",
										Optional:    true,
									},
									"manage": schema.SingleNestedAttribute{
										Description: "Permissions for managing users and roles.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all manage permissions are enabled.",
												Optional:    true,
											},
											"roles": schema.SingleNestedAttribute{
												Description: "Permissions for managing roles.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all role permissions are enabled.",
														Optional:    true,
													},
													"create": schema.BoolAttribute{
														Description: "Permission to create roles.",
														Optional:    true,
													},
													"delete": schema.BoolAttribute{
														Description: "Permission to delete roles.",
														Optional:    true,
													},
													"edit": schema.BoolAttribute{
														Description: "Permission to edit roles.",
														Optional:    true,
													},
												},
											},
											"users": schema.SingleNestedAttribute{
												Description: "Permissions for managing users.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"all": schema.BoolAttribute{
														Description: "Indicates if all user permissions are enabled.",
														Optional:    true,
													},
													"create": schema.BoolAttribute{
														Description: "Permission to create users.",
														Optional:    true,
													},
													"delete": schema.BoolAttribute{
														Description: "Permission to delete users.",
														Optional:    true,
													},
													"edit": schema.BoolAttribute{
														Description: "Permission to edit users.",
														Optional:    true,
													},
												},
											},
										},
									},
									"read": schema.BoolAttribute{
										Description: "Permission to read users and roles.",
										Optional:    true,
									},
								},
							},
						},
					},
					"user": schema.SingleNestedAttribute{
						Description: "Permissions for managing users.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all user permissions are enabled.",
								Optional:    true,
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing users.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all user management permissions are enabled.",
										Optional:    true,
									},
									"upsert": schema.BoolAttribute{
										Description: "Permission to upsert users.",
										Optional:    true,
									},
								},
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read users.",
								Optional:    true,
							},
						},
					},
					"vulnerability": schema.SingleNestedAttribute{
						Description: "Permissions for managing vulnerabilities.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"all": schema.BoolAttribute{
								Description: "Indicates if all vulnerability permissions are enabled.",
								Optional:    true,
							},
							"manage": schema.SingleNestedAttribute{
								Description: "Permissions for managing vulnerabilities.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all vulnerability management permissions are enabled.",
										Optional:    true,
									},
									"ignore": schema.BoolAttribute{
										Description: "Permission to ignore vulnerabilities.",
										Optional:    true,
									},
									"resolve": schema.BoolAttribute{
										Description: "Permission to resolve vulnerabilities.",
										Optional:    true,
									},
									"write": schema.BoolAttribute{
										Description: "Permission to write vulnerabilities.",
										Optional:    true,
									},
								},
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read vulnerabilities.",
								Optional:    true,
							},
						},
					},
				},
			},
		},
	}
}

// RoleResourceModel maps the RoleSettings schema data.
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

	// Define the role to create
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
	mapRoleSettingsToPlan(role, &state)
	tflog.Debug(ctx, "Setting refreshed state for role", map[string]any{"role_id": state.ID.ValueString()})
	diags = resp.State.Set(ctx, state)
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
	mapRoleSettingsToPlan(updatedRole, &plan)

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

// Utility Functions

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

// mapRoleSettingsToPlan updates a RoleResourceModel with data from a RoleSettings object.
func mapRoleSettingsToPlan(role *armis.RoleSettings, plan *RoleResourceModel) {
	plan.Name = types.StringValue(role.Name)
	plan.ID = types.StringValue(strconv.Itoa(role.ID))
}
