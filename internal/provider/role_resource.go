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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
							"read": schema.SingleNestedAttribute{
								Description: "Permission to read alerts.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all read permissions for alerts are enabled.",
										Optional:    true,
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
							"read": schema.SingleNestedAttribute{
								Description: "Permission to read devices.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"all": schema.BoolAttribute{
										Description: "Indicates if all read permissions for devices are enabled.",
										Optional:    true,
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
									"read": schema.SingleNestedAttribute{
										Description: "Permission to read sites and sensors.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"all": schema.BoolAttribute{
												Description: "Indicates if all read permissions for sites and sensors are enabled.",
												Optional:    true,
											},
										},
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

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating role", map[string]any{"name": plan.Name.ValueString()})

	if plan.Permissions == nil {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Permissions are required but not provided.",
		)
		return
	}

	if plan.Permissions.AdvancedPermissions == nil {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Advanced permissions are required but not provided.",
		)
		return
	}

	role := BuildRoleRequest(plan)

	tflog.Debug(ctx, "Creating role request", map[string]any{
		"name":        plan.Name.ValueString(),
		"permissions": plan.Permissions,
	})

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
	roleState := BuildRoleResourceModel(role, state)
	tflog.Debug(ctx, "Setting refreshed state for role", map[string]any{"role_id": state.ID.ValueString()})
	diags = resp.State.Set(ctx, roleState)
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
	role := BuildRoleRequest(plan)

	tflog.Debug(ctx, "Creating role request", map[string]any{
		"name":        plan.Name.ValueString(),
		"permissions": plan.Permissions,
	})

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
	updatedPlan := BuildRoleResourceModel(updatedRole, plan)

	// Save the updated state
	tflog.Info(ctx, "Setting updated state for role", map[string]any{"role_id": state.ID.ValueString()})
	diags = resp.State.Set(ctx, updatedPlan)
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

func BuildRoleRequest(role RoleResourceModel) armis.RoleSettings {
	return armis.RoleSettings{
		Name: role.Name.ValueString(),
		Permissions: armis.Permissions{
			AdvancedPermissions: armis.AdvancedPermissions{
				All: role.Permissions.AdvancedPermissions.All.ValueBool(),
				Behavioral: armis.Behavioral{
					All: role.Permissions.AdvancedPermissions.Behavioral.All.ValueBool(),
					ApplicationName: armis.Permission{
						All: role.Permissions.AdvancedPermissions.Behavioral.ApplicationName.ValueBool(),
					},
					HostName: armis.Permission{
						All: role.Permissions.AdvancedPermissions.Behavioral.HostName.ValueBool(),
					},
					ServiceName: armis.Permission{
						All: role.Permissions.AdvancedPermissions.Behavioral.ServiceName.ValueBool(),
					},
				},
				Device: armis.DeviceAdvanced{
					All: role.Permissions.AdvancedPermissions.Device.All.ValueBool(),
					DeviceNames: armis.Permission{
						All: role.Permissions.AdvancedPermissions.Device.DeviceNames.ValueBool(),
					},
					IPAddresses: armis.Permission{
						All: role.Permissions.AdvancedPermissions.Device.IPAddresses.ValueBool(),
					},
					MACAddresses: armis.Permission{
						All: role.Permissions.AdvancedPermissions.Device.MACAddresses.ValueBool(),
					},
					PhoneNumbers: armis.Permission{
						All: role.Permissions.AdvancedPermissions.Device.PhoneNumbers.ValueBool(),
					},
				},
			},
			Alert: armis.Alert{
				All: role.Permissions.Alert.All.ValueBool(),
				Manage: armis.Manage{
					All: role.Permissions.Alert.Manage.All.ValueBool(),
					Resolve: armis.Permission{
						All: role.Permissions.Alert.Manage.Resolve.ValueBool(),
					},
					Suppress: armis.Permission{
						All: role.Permissions.Alert.Manage.Suppress.ValueBool(),
					},
					WhitelistDevices: armis.Permission{
						All: role.Permissions.Alert.Manage.WhitelistDevices.ValueBool(),
					},
				},
				Read: armis.Permission{
					All: func() bool {
						if role.Permissions.Alert.Read != nil {
							return role.Permissions.Alert.Read.All.ValueBool()
						}
						return false
					}(),
				},
			},
			Device: armis.Device{
				All: role.Permissions.Device.All.ValueBool(),
				Manage: armis.ManageDevice{
					All: role.Permissions.Device.Manage.All.ValueBool(),
					Create: armis.Permission{
						All: role.Permissions.Device.Manage.Create.ValueBool(),
					},
					Delete: armis.Permission{
						All: role.Permissions.Device.Manage.Delete.ValueBool(),
					},
					Edit: armis.Permission{
						All: role.Permissions.Device.Manage.Edit.ValueBool(),
					},
					Enforce: armis.Enforce{
						All: role.Permissions.Device.Manage.Enforce.All.ValueBool(),
						Create: armis.Permission{
							All: role.Permissions.Device.Manage.Enforce.Create.ValueBool(),
						},
						Delete: armis.Permission{
							All: role.Permissions.Device.Manage.Enforce.Delete.ValueBool(),
						},
					},
					Merge: armis.Permission{
						All: role.Permissions.Device.Manage.Merge.ValueBool(),
					},
					RequestDeletedData: armis.Permission{
						All: role.Permissions.Device.Manage.RequestDeletedData.ValueBool(),
					},
					Tags: armis.Permission{
						All: role.Permissions.Device.Manage.Tags.ValueBool(),
					},
				},
				Read: armis.Permission{
					All: func() bool {
						if role.Permissions.Device.Read != nil {
							return role.Permissions.Device.Read.All.ValueBool()
						}
						return false
					}(),
				},
			},
			Policy: armis.Policy{
				All: role.Permissions.Policy.All.ValueBool(),
				Manage: armis.Permission{
					All: role.Permissions.Policy.Manage.ValueBool(),
				},
				Read: armis.Permission{
					All: role.Permissions.Policy.Read.ValueBool(),
				},
			},
			Report: armis.Report{
				All: role.Permissions.Report.All.ValueBool(),
				Export: armis.Permission{
					All: role.Permissions.Report.Export.ValueBool(),
				},
				Manage: armis.ManageReport{
					All: role.Permissions.Report.Manage.All.ValueBool(),
					Create: armis.Permission{
						All: role.Permissions.Report.Manage.Create.ValueBool(),
					},
					Delete: armis.Permission{
						All: role.Permissions.Report.Manage.Delete.ValueBool(),
					},
					Edit: armis.Permission{
						All: role.Permissions.Report.Manage.Edit.ValueBool(),
					},
				},
				Read: armis.Permission{
					All: role.Permissions.Report.Read.ValueBool(),
				},
			},
			RiskFactor: armis.RiskFactor{
				All: role.Permissions.RiskFactor.All.ValueBool(),
				Manage: armis.ManageRisk{
					All: role.Permissions.RiskFactor.Manage.All.ValueBool(),
					Customization: armis.Customization{
						All: role.Permissions.RiskFactor.Manage.Customization.All.ValueBool(),
						Create: armis.Permission{
							All: role.Permissions.RiskFactor.Manage.Customization.Create.ValueBool(),
						},
						Disable: armis.Permission{
							All: role.Permissions.RiskFactor.Manage.Customization.Disable.ValueBool(),
						},
						Edit: armis.Permission{
							All: role.Permissions.RiskFactor.Manage.Customization.Edit.ValueBool(),
						},
					},
					Status: armis.Status{
						All: role.Permissions.RiskFactor.Manage.Status.All.ValueBool(),
						Ignore: armis.Permission{
							All: role.Permissions.RiskFactor.Manage.Status.Ignore.ValueBool(),
						},
						Resolve: armis.Permission{
							All: role.Permissions.RiskFactor.Manage.Status.Resolve.ValueBool(),
						},
					},
				},
				Read: armis.Permission{
					All: role.Permissions.RiskFactor.Read.ValueBool(),
				},
			},
			Settings: armis.Settings{
				All: role.Permissions.Settings.All.ValueBool(),
				AuditLog: armis.Permission{
					All: role.Permissions.Settings.AuditLog.ValueBool(),
				},
				Boundary: armis.Boundary{
					All: role.Permissions.Settings.Boundary.All.ValueBool(),
					Manage: armis.ManageBoundary{
						All: role.Permissions.Settings.Boundary.Manage.All.ValueBool(),
						Create: armis.Permission{
							All: role.Permissions.Settings.Boundary.Manage.Create.ValueBool(),
						},
						Delete: armis.Permission{
							All: role.Permissions.Settings.Boundary.Manage.Delete.ValueBool(),
						},
						Edit: armis.Permission{
							All: role.Permissions.Settings.Boundary.Manage.Edit.ValueBool(),
						},
					},
					Read: armis.Permission{
						All: role.Permissions.Settings.Boundary.Read.ValueBool(),
					},
				},
				BusinessImpact: armis.ManageAndRead{
					All: role.Permissions.Settings.BusinessImpact.All.ValueBool(),
					Manage: armis.Permission{
						All: role.Permissions.Settings.BusinessImpact.Manage.ValueBool(),
					},
					Read: armis.Permission{
						All: role.Permissions.Settings.BusinessImpact.Read.ValueBool(),
					},
				},
				Collector: armis.ManageAndRead{
					All: role.Permissions.Settings.Collector.All.ValueBool(),
					Manage: armis.Permission{
						All: role.Permissions.Settings.Collector.Manage.ValueBool(),
					},
					Read: armis.Permission{
						All: role.Permissions.Settings.Collector.Read.ValueBool(),
					},
				},
				CustomProperties: armis.ManageAndRead{
					All: role.Permissions.Settings.CustomProperties.All.ValueBool(),
					Manage: armis.Permission{
						All: role.Permissions.Settings.CustomProperties.Manage.ValueBool(),
					},
					Read: armis.Permission{
						All: role.Permissions.Settings.CustomProperties.Read.ValueBool(),
					},
				},
				Integration: armis.ManageAndRead{
					All: role.Permissions.Settings.Integration.All.ValueBool(),
					Manage: armis.Permission{
						All: role.Permissions.Settings.Integration.Manage.ValueBool(),
					},
					Read: armis.Permission{
						All: role.Permissions.Settings.Integration.Read.ValueBool(),
					},
				},
				InternalIps: armis.ManageAndRead{
					All: role.Permissions.Settings.InternalIps.All.ValueBool(),
					Manage: armis.Permission{
						All: role.Permissions.Settings.InternalIps.Manage.ValueBool(),
					},
					Read: armis.Permission{
						All: role.Permissions.Settings.InternalIps.Read.ValueBool(),
					},
				},
				Notifications: armis.ManageAndRead{
					All: role.Permissions.Settings.Notifications.All.ValueBool(),
					Manage: armis.Permission{
						All: role.Permissions.Settings.Notifications.Manage.ValueBool(),
					},
					Read: armis.Permission{
						All: role.Permissions.Settings.Notifications.Read.ValueBool(),
					},
				},
				OIDC: armis.ManageAndRead{
					All: role.Permissions.Settings.OIDC.All.ValueBool(),
					Manage: armis.Permission{
						All: role.Permissions.Settings.OIDC.Manage.ValueBool(),
					},
					Read: armis.Permission{
						All: role.Permissions.Settings.OIDC.Read.ValueBool(),
					},
				},
				SAML: armis.ManageAndRead{
					All: role.Permissions.Settings.SAML.All.ValueBool(),
					Manage: armis.Permission{
						All: role.Permissions.Settings.SAML.Manage.ValueBool(),
					},
					Read: armis.Permission{
						All: role.Permissions.Settings.SAML.Read.ValueBool(),
					},
				},
				SecretKey: armis.Permission{
					All: role.Permissions.Settings.SecretKey.ValueBool(),
				},
				SecuritySettings: armis.Permission{
					All: role.Permissions.Settings.SecuritySettings.ValueBool(),
				},
				SitesAndSensors: armis.SitesAndSensors{
					All: role.Permissions.Settings.SitesAndSensors.All.ValueBool(),
					Manage: armis.ManageSensors{
						All: role.Permissions.Settings.SitesAndSensors.Manage.All.ValueBool(),
						Sensors: armis.Permission{
							All: role.Permissions.Settings.SitesAndSensors.Manage.Sensors.ValueBool(),
						},
						Sites: armis.Permission{
							All: role.Permissions.Settings.SitesAndSensors.Manage.Sites.ValueBool(),
						},
					},
					Read: armis.Permission{
						All: func() bool {
							if role.Permissions.Settings.SitesAndSensors.Read != nil {
								return role.Permissions.Settings.SitesAndSensors.Read.All.ValueBool()
							}
							return false
						}(),
					},
				},
				UsersAndRoles: armis.UsersAndRoles{
					All: role.Permissions.Settings.UsersAndRoles.All.ValueBool(),
					Manage: armis.ManageUsers{
						All: role.Permissions.Settings.UsersAndRoles.Manage.All.ValueBool(),
						Roles: armis.UserActions{
							All: role.Permissions.Settings.UsersAndRoles.Manage.Roles.All.ValueBool(),
							Create: armis.Permission{
								All: role.Permissions.Settings.UsersAndRoles.Manage.Roles.Create.ValueBool(),
							},
							Delete: armis.Permission{
								All: role.Permissions.Settings.UsersAndRoles.Manage.Roles.Delete.ValueBool(),
							},
							Edit: armis.Permission{
								All: role.Permissions.Settings.UsersAndRoles.Manage.Roles.Edit.ValueBool(),
							},
						},
						Users: armis.UserActions{
							All: role.Permissions.Settings.UsersAndRoles.Manage.Users.All.ValueBool(),
							Create: armis.Permission{
								All: role.Permissions.Settings.UsersAndRoles.Manage.Users.Create.ValueBool(),
							},
							Delete: armis.Permission{
								All: role.Permissions.Settings.UsersAndRoles.Manage.Users.Delete.ValueBool(),
							},
							Edit: armis.Permission{
								All: role.Permissions.Settings.UsersAndRoles.Manage.Users.Edit.ValueBool(),
							},
						},
					},
					Read: armis.Permission{
						All: role.Permissions.Settings.UsersAndRoles.Read.ValueBool(),
					},
				},
			},
			User: armis.User{
				All: role.Permissions.User.All.ValueBool(),
				Manage: armis.ManageUser{
					All: role.Permissions.User.Manage.All.ValueBool(),
					Upsert: armis.Permission{
						All: role.Permissions.User.Manage.Upsert.ValueBool(),
					},
				},
				Read: armis.Permission{
					All: role.Permissions.User.Read.ValueBool(),
				},
			},
			Vulnerability: armis.Vulnerability{
				All: role.Permissions.Vulnerability.All.ValueBool(),
				Manage: armis.ManageVuln{
					All: role.Permissions.Vulnerability.Manage.All.ValueBool(),
					Ignore: armis.Permission{
						All: role.Permissions.Vulnerability.Manage.Ignore.ValueBool(),
					},
					Resolve: armis.Permission{
						All: role.Permissions.Vulnerability.Manage.Resolve.ValueBool(),
					},
					Write: armis.Permission{
						All: role.Permissions.Vulnerability.Manage.Write.ValueBool(),
					},
				},
				Read: armis.Permission{
					All: role.Permissions.Vulnerability.Read.ValueBool(),
				},
			},
		},
	}
}

func BuildRoleResourceModel(role *armis.RoleSettings, model RoleResourceModel) RoleResourceModel {
	result := model
	result.Name = types.StringValue(role.Name)
	result.ID = types.StringValue(strconv.Itoa(role.ID))

	// Advanced Permissions
	result.Permissions.AdvancedPermissions.All = types.BoolValue(role.Permissions.AdvancedPermissions.All)

	// Advanced Permissions - Behavioral
	result.Permissions.AdvancedPermissions.Behavioral.All = types.BoolValue(role.Permissions.AdvancedPermissions.Behavioral.All)
	result.Permissions.AdvancedPermissions.Behavioral.ApplicationName = types.BoolValue(role.Permissions.AdvancedPermissions.Behavioral.ApplicationName.All)
	result.Permissions.AdvancedPermissions.Behavioral.HostName = types.BoolValue(role.Permissions.AdvancedPermissions.Behavioral.HostName.All)
	result.Permissions.AdvancedPermissions.Behavioral.ServiceName = types.BoolValue(role.Permissions.AdvancedPermissions.Behavioral.ServiceName.All)

	// Advanced Permissions - Device
	result.Permissions.AdvancedPermissions.Device.All = types.BoolValue(role.Permissions.AdvancedPermissions.Device.All)
	result.Permissions.AdvancedPermissions.Device.DeviceNames = types.BoolValue(role.Permissions.AdvancedPermissions.Device.DeviceNames.All)
	result.Permissions.AdvancedPermissions.Device.IPAddresses = types.BoolValue(role.Permissions.AdvancedPermissions.Device.IPAddresses.All)
	result.Permissions.AdvancedPermissions.Device.MACAddresses = types.BoolValue(role.Permissions.AdvancedPermissions.Device.MACAddresses.All)
	result.Permissions.AdvancedPermissions.Device.PhoneNumbers = types.BoolValue(role.Permissions.AdvancedPermissions.Device.PhoneNumbers.All)

	// Alert Permissions
	result.Permissions.Alert.All = types.BoolValue(role.Permissions.Alert.All)
	result.Permissions.Alert.Read.All = types.BoolValue(role.Permissions.Alert.Read.All)

	// Alert Manage Permissions
	result.Permissions.Alert.Manage.All = types.BoolValue(role.Permissions.Alert.Manage.All)
	result.Permissions.Alert.Manage.Resolve = types.BoolValue(role.Permissions.Alert.Manage.Resolve.All)
	result.Permissions.Alert.Manage.Suppress = types.BoolValue(role.Permissions.Alert.Manage.Suppress.All)
	result.Permissions.Alert.Manage.WhitelistDevices = types.BoolValue(role.Permissions.Alert.Manage.WhitelistDevices.All)

	// Device Permissions
	result.Permissions.Device.All = types.BoolValue(role.Permissions.Device.All)
	result.Permissions.Device.Read.All = types.BoolValue(role.Permissions.Device.Read.All)

	// Device Manage Permissions
	result.Permissions.Device.Manage.All = types.BoolValue(role.Permissions.Device.Manage.All)
	result.Permissions.Device.Manage.Create = types.BoolValue(role.Permissions.Device.Manage.Create.All)
	result.Permissions.Device.Manage.Delete = types.BoolValue(role.Permissions.Device.Manage.Delete.All)
	result.Permissions.Device.Manage.Edit = types.BoolValue(role.Permissions.Device.Manage.Edit.All)
	result.Permissions.Device.Manage.Merge = types.BoolValue(role.Permissions.Device.Manage.Merge.All)
	result.Permissions.Device.Manage.RequestDeletedData = types.BoolValue(role.Permissions.Device.Manage.RequestDeletedData.All)
	result.Permissions.Device.Manage.Tags = types.BoolValue(role.Permissions.Device.Manage.Tags.All)

	// Device Enforce Permissions
	result.Permissions.Device.Manage.Enforce.All = types.BoolValue(role.Permissions.Device.Manage.Enforce.All)
	result.Permissions.Device.Manage.Enforce.Create = types.BoolValue(role.Permissions.Device.Manage.Enforce.Create.All)
	result.Permissions.Device.Manage.Enforce.Delete = types.BoolValue(role.Permissions.Device.Manage.Enforce.Delete.All)

	// Policy Permissions
	result.Permissions.Policy.All = types.BoolValue(role.Permissions.Policy.All)
	result.Permissions.Policy.Manage = types.BoolValue(role.Permissions.Policy.Manage.All)
	result.Permissions.Policy.Read = types.BoolValue(role.Permissions.Policy.Read.All)

	// Report Permissions
	result.Permissions.Report.All = types.BoolValue(role.Permissions.Report.All)
	result.Permissions.Report.Export = types.BoolValue(role.Permissions.Report.Export.All)
	result.Permissions.Report.Read = types.BoolValue(role.Permissions.Report.Read.All)

	// Report Manage Permissions
	result.Permissions.Report.Manage.All = types.BoolValue(role.Permissions.Report.Manage.All)
	result.Permissions.Report.Manage.Create = types.BoolValue(role.Permissions.Report.Manage.Create.All)
	result.Permissions.Report.Manage.Delete = types.BoolValue(role.Permissions.Report.Manage.Delete.All)
	result.Permissions.Report.Manage.Edit = types.BoolValue(role.Permissions.Report.Manage.Edit.All)

	// Risk Factor Permissions
	result.Permissions.RiskFactor.All = types.BoolValue(role.Permissions.RiskFactor.All)
	result.Permissions.RiskFactor.Read = types.BoolValue(role.Permissions.RiskFactor.Read.All)

	// Risk Factor Manage Permissions
	result.Permissions.RiskFactor.Manage.All = types.BoolValue(role.Permissions.RiskFactor.Manage.All)

	// Risk Factor Customization Permissions
	result.Permissions.RiskFactor.Manage.Customization.All = types.BoolValue(role.Permissions.RiskFactor.Manage.Customization.All)
	result.Permissions.RiskFactor.Manage.Customization.Create = types.BoolValue(role.Permissions.RiskFactor.Manage.Customization.Create.All)
	result.Permissions.RiskFactor.Manage.Customization.Disable = types.BoolValue(role.Permissions.RiskFactor.Manage.Customization.Disable.All)
	result.Permissions.RiskFactor.Manage.Customization.Edit = types.BoolValue(role.Permissions.RiskFactor.Manage.Customization.Edit.All)

	// Risk Factor Status Permissions
	result.Permissions.RiskFactor.Manage.Status.All = types.BoolValue(role.Permissions.RiskFactor.Manage.Status.All)
	result.Permissions.RiskFactor.Manage.Status.Ignore = types.BoolValue(role.Permissions.RiskFactor.Manage.Status.Ignore.All)
	result.Permissions.RiskFactor.Manage.Status.Resolve = types.BoolValue(role.Permissions.RiskFactor.Manage.Status.Resolve.All)

	// Settings Permissions
	result.Permissions.Settings.All = types.BoolValue(role.Permissions.Settings.All)
	result.Permissions.Settings.AuditLog = types.BoolValue(role.Permissions.Settings.AuditLog.All)
	result.Permissions.Settings.SecretKey = types.BoolValue(role.Permissions.Settings.SecretKey.All)
	result.Permissions.Settings.SecuritySettings = types.BoolValue(role.Permissions.Settings.SecuritySettings.All)

	// Settings Boundary Permissions
	result.Permissions.Settings.Boundary.All = types.BoolValue(role.Permissions.Settings.Boundary.All)
	result.Permissions.Settings.Boundary.Read = types.BoolValue(role.Permissions.Settings.Boundary.Read.All)
	result.Permissions.Settings.Boundary.Manage.All = types.BoolValue(role.Permissions.Settings.Boundary.Manage.All)
	result.Permissions.Settings.Boundary.Manage.Create = types.BoolValue(role.Permissions.Settings.Boundary.Manage.Create.All)
	result.Permissions.Settings.Boundary.Manage.Delete = types.BoolValue(role.Permissions.Settings.Boundary.Manage.Delete.All)
	result.Permissions.Settings.Boundary.Manage.Edit = types.BoolValue(role.Permissions.Settings.Boundary.Manage.Edit.All)

	// Settings Business Impact Permissions
	result.Permissions.Settings.BusinessImpact.All = types.BoolValue(role.Permissions.Settings.BusinessImpact.All)
	result.Permissions.Settings.BusinessImpact.Manage = types.BoolValue(role.Permissions.Settings.BusinessImpact.Manage.All)
	result.Permissions.Settings.BusinessImpact.Read = types.BoolValue(role.Permissions.Settings.BusinessImpact.Read.All)

	// Settings Collector Permissions
	result.Permissions.Settings.Collector.All = types.BoolValue(role.Permissions.Settings.Collector.All)
	result.Permissions.Settings.Collector.Manage = types.BoolValue(role.Permissions.Settings.Collector.Manage.All)
	result.Permissions.Settings.Collector.Read = types.BoolValue(role.Permissions.Settings.Collector.Read.All)

	// Settings Custom Properties Permissions
	result.Permissions.Settings.CustomProperties.All = types.BoolValue(role.Permissions.Settings.CustomProperties.All)
	result.Permissions.Settings.CustomProperties.Manage = types.BoolValue(role.Permissions.Settings.CustomProperties.Manage.All)
	result.Permissions.Settings.CustomProperties.Read = types.BoolValue(role.Permissions.Settings.CustomProperties.Read.All)

	// Settings Integration Permissions
	result.Permissions.Settings.Integration.All = types.BoolValue(role.Permissions.Settings.Integration.All)
	result.Permissions.Settings.Integration.Manage = types.BoolValue(role.Permissions.Settings.Integration.Manage.All)
	result.Permissions.Settings.Integration.Read = types.BoolValue(role.Permissions.Settings.Integration.Read.All)

	// Settings Internal IPs Permissions
	result.Permissions.Settings.InternalIps.All = types.BoolValue(role.Permissions.Settings.InternalIps.All)
	result.Permissions.Settings.InternalIps.Manage = types.BoolValue(role.Permissions.Settings.InternalIps.Manage.All)
	result.Permissions.Settings.InternalIps.Read = types.BoolValue(role.Permissions.Settings.InternalIps.Read.All)

	// Settings Notifications Permissions
	result.Permissions.Settings.Notifications.All = types.BoolValue(role.Permissions.Settings.Notifications.All)
	result.Permissions.Settings.Notifications.Manage = types.BoolValue(role.Permissions.Settings.Notifications.Manage.All)
	result.Permissions.Settings.Notifications.Read = types.BoolValue(role.Permissions.Settings.Notifications.Read.All)

	// Settings OIDC Permissions
	result.Permissions.Settings.OIDC.All = types.BoolValue(role.Permissions.Settings.OIDC.All)
	result.Permissions.Settings.OIDC.Manage = types.BoolValue(role.Permissions.Settings.OIDC.Manage.All)
	result.Permissions.Settings.OIDC.Read = types.BoolValue(role.Permissions.Settings.OIDC.Read.All)

	// Settings SAML Permissions
	result.Permissions.Settings.SAML.All = types.BoolValue(role.Permissions.Settings.SAML.All)
	result.Permissions.Settings.SAML.Manage = types.BoolValue(role.Permissions.Settings.SAML.Manage.All)
	result.Permissions.Settings.SAML.Read = types.BoolValue(role.Permissions.Settings.SAML.Read.All)

	// Settings Sites and Sensors Permissions
	result.Permissions.Settings.SitesAndSensors.All = types.BoolValue(role.Permissions.Settings.SitesAndSensors.All)
	result.Permissions.Settings.SitesAndSensors.Read.All = types.BoolValue(role.Permissions.Settings.SitesAndSensors.Read.All)
	result.Permissions.Settings.SitesAndSensors.Manage.All = types.BoolValue(role.Permissions.Settings.SitesAndSensors.Manage.All)
	result.Permissions.Settings.SitesAndSensors.Manage.Sensors = types.BoolValue(role.Permissions.Settings.SitesAndSensors.Manage.Sensors.All)
	result.Permissions.Settings.SitesAndSensors.Manage.Sites = types.BoolValue(role.Permissions.Settings.SitesAndSensors.Manage.Sites.All)

	// Settings Users and Roles Permissions
	result.Permissions.Settings.UsersAndRoles.All = types.BoolValue(role.Permissions.Settings.UsersAndRoles.All)
	result.Permissions.Settings.UsersAndRoles.Read = types.BoolValue(role.Permissions.Settings.UsersAndRoles.Read.All)
	result.Permissions.Settings.UsersAndRoles.Manage.All = types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.All)

	// Settings Users and Roles - Roles Permissions
	result.Permissions.Settings.UsersAndRoles.Manage.Roles.All = types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Roles.All)
	result.Permissions.Settings.UsersAndRoles.Manage.Roles.Create = types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Roles.Create.All)
	result.Permissions.Settings.UsersAndRoles.Manage.Roles.Delete = types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Roles.Delete.All)
	result.Permissions.Settings.UsersAndRoles.Manage.Roles.Edit = types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Roles.Edit.All)

	// Settings Users and Roles - Users Permissions
	result.Permissions.Settings.UsersAndRoles.Manage.Users.All = types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Users.All)
	result.Permissions.Settings.UsersAndRoles.Manage.Users.Create = types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Users.Create.All)
	result.Permissions.Settings.UsersAndRoles.Manage.Users.Delete = types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Users.Delete.All)
	result.Permissions.Settings.UsersAndRoles.Manage.Users.Edit = types.BoolValue(role.Permissions.Settings.UsersAndRoles.Manage.Users.Edit.All)

	// User Permissions
	result.Permissions.User.All = types.BoolValue(role.Permissions.User.All)
	result.Permissions.User.Read = types.BoolValue(role.Permissions.User.Read.All)
	result.Permissions.User.Manage.All = types.BoolValue(role.Permissions.User.Manage.All)
	result.Permissions.User.Manage.Upsert = types.BoolValue(role.Permissions.User.Manage.Upsert.All)

	// Vulnerability Permissions
	result.Permissions.Vulnerability.All = types.BoolValue(role.Permissions.Vulnerability.All)
	result.Permissions.Vulnerability.Read = types.BoolValue(role.Permissions.Vulnerability.Read.All)
	result.Permissions.Vulnerability.Manage.All = types.BoolValue(role.Permissions.Vulnerability.Manage.All)
	result.Permissions.Vulnerability.Manage.Ignore = types.BoolValue(role.Permissions.Vulnerability.Manage.Ignore.All)
	result.Permissions.Vulnerability.Manage.Resolve = types.BoolValue(role.Permissions.Vulnerability.Manage.Resolve.All)
	result.Permissions.Vulnerability.Manage.Write = types.BoolValue(role.Permissions.Vulnerability.Manage.Write.All)

	return result
}
