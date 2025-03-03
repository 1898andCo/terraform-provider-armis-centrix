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
									"application_name": schema.BoolAttribute{
										Description: "Permission to access application names.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"host_name": schema.BoolAttribute{
										Description: "Permission to access host names.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"service_name": schema.BoolAttribute{
										Description: "Permission to access service names.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
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
									"device_names": schema.BoolAttribute{
										Description: "Permission to access device names.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"ip_addresses": schema.BoolAttribute{
										Description: "Permission to access device IP addresses.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"mac_addresses": schema.BoolAttribute{
										Description: "Permission to access device MAC addresses.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"phone_numbers": schema.BoolAttribute{
										Description: "Permission to access device phone numbers.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
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
									"resolve": schema.BoolAttribute{
										Description: "Permission to resolve alerts.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"suppress": schema.BoolAttribute{
										Description: "Permission to suppress alerts.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"whitelist_devices": schema.BoolAttribute{
										Description: "Permission to whitelist devices in alerts.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
								},
							},
							"read": schema.BoolAttribute{
								Description: "Permission to read alerts.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
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

// PermissionsModel maps the Permissions schema data.
type PermissionsModel struct {
	AdvancedPermissions *AdvancedPermissionsModel `tfsdk:"advanced_permissions"`
	Alert               *AlertModel               `tfsdk:"alert"`
	// Device              DeviceModel              `tfsdk:"device"`
	// Policy              PolicyModel              `tfsdk:"policy"`
	// Report              ReportModel              `tfsdk:"report"`
	// RiskFactor          RiskFactorModel          `tfsdk:"risk_factor"`
	// Settings            SettingsModel            `tfsdk:"settings"`
	// User                UserModel                `tfsdk:"user"`
	// Vulnerability       VulnerabilityModel       `tfsdk:"vulnerability"`
}

// AdvancedPermissionsModel maps the AdvancedPermissions schema data.
type AdvancedPermissionsModel struct {
	All        types.Bool           `tfsdk:"all"`
	Behavioral *BehavioralModel     `tfsdk:"behavioral"`
	Device     *DeviceAdvancedModel `tfsdk:"device"`
}

// BehavioralModel maps the Behavioral schema data.
type BehavioralModel struct {
	All             types.Bool `tfsdk:"all"`
	ApplicationName types.Bool `tfsdk:"application_name"`
	HostName        types.Bool `tfsdk:"host_name"`
	ServiceName     types.Bool `tfsdk:"service_name"`
}

// DeviceAdvancedModel maps the DeviceAdvanced schema data.
type DeviceAdvancedModel struct {
	All          types.Bool `tfsdk:"all"`
	DeviceNames  types.Bool `tfsdk:"device_names"`
	IPAddresses  types.Bool `tfsdk:"ip_addresses"`
	MACAddresses types.Bool `tfsdk:"mac_addresses"`
	PhoneNumbers types.Bool `tfsdk:"phone_numbers"`
}

// AlertModel maps the Alert schema data.
type AlertModel struct {
	All    types.Bool   `tfsdk:"all"`
	Manage *ManageModel `tfsdk:"manage"`
	Read   types.Bool   `tfsdk:"read"`
}

// ManageModel maps the Manage schema data.
type ManageModel struct {
	All              types.Bool `tfsdk:"all"`
	Resolve          types.Bool `tfsdk:"resolve"`
	Suppress         types.Bool `tfsdk:"suppress"`
	WhitelistDevices types.Bool `tfsdk:"whitelist_devices"`
}

// // DeviceModel maps the Device schema data.
// type DeviceModel struct {
// 	All    types.Bool        `tfsdk:"all"`
// 	Manage ManageDeviceModel `tfsdk:"manage"`
// 	Read   PermissionModel   `tfsdk:"read"`
// }
//
// // ManageDeviceModel maps the ManageDevice schema data.
// type ManageDeviceModel struct {
// 	All                types.Bool      `tfsdk:"all"`
// 	Create             PermissionModel `tfsdk:"create"`
// 	Delete             PermissionModel `tfsdk:"delete"`
// 	Edit               PermissionModel `tfsdk:"edit"`
// 	Enforce            EnforceModel    `tfsdk:"enforce"`
// 	Merge              PermissionModel `tfsdk:"merge"`
// 	RequestDeletedData PermissionModel `tfsdk:"request_deleted_data"`
// 	Tags               PermissionModel `tfsdk:"tags"`
// }
//
// // EnforceModel maps the Enforce schema data.
// type EnforceModel struct {
// 	All    types.Bool      `tfsdk:"all"`
// 	Create PermissionModel `tfsdk:"create"`
// 	Delete PermissionModel `tfsdk:"delete"`
// }
//
// // PolicyModel maps the Policy schema data.
// type PolicyModel struct {
// 	All    types.Bool      `tfsdk:"all"`
// 	Manage PermissionModel `tfsdk:"manage"`
// 	Read   PermissionModel `tfsdk:"read"`
// }
//
// // ReportModel maps the Report schema data.
// type ReportModel struct {
// 	All    types.Bool        `tfsdk:"all"`
// 	Export PermissionModel   `tfsdk:"export"`
// 	Manage ManageReportModel `tfsdk:"manage"`
// 	Read   PermissionModel   `tfsdk:"read"`
// }
//
// // ManageReportModel maps the ManageReport schema data.
// type ManageReportModel struct {
// 	All    types.Bool      `tfsdk:"all"`
// 	Create PermissionModel `tfsdk:"create"`
// 	Delete PermissionModel `tfsdk:"delete"`
// 	Edit   PermissionModel `tfsdk:"edit"`
// }
//
// // RiskFactorModel maps the RiskFactor schema data.
// type RiskFactorModel struct {
// 	All    types.Bool      `tfsdk:"all"`
// 	Manage ManageRiskModel `tfsdk:"manage"`
// 	Read   PermissionModel `tfsdk:"read"`
// }
//
// // ManageRiskModel maps the ManageRisk schema data.
// type ManageRiskModel struct {
// 	All           types.Bool         `tfsdk:"all"`
// 	Customization CustomizationModel `tfsdk:"customization"`
// 	Status        StatusModel        `tfsdk:"status"`
// }
//
// // CustomizationModel maps the Customization schema data.
// type CustomizationModel struct {
// 	All     types.Bool      `tfsdk:"all"`
// 	Create  PermissionModel `tfsdk:"create"`
// 	Disable PermissionModel `tfsdk:"disable"`
// 	Edit    PermissionModel `tfsdk:"edit"`
// }
//
// // StatusModel maps the Status schema data.
// type StatusModel struct {
// 	All     types.Bool      `tfsdk:"all"`
// 	Ignore  PermissionModel `tfsdk:"ignore"`
// 	Resolve PermissionModel `tfsdk:"resolve"`
// }
//
// // SettingsModel maps the Settings schema data.
// type SettingsModel struct {
// 	All              types.Bool           `tfsdk:"all"`
// 	AuditLog         PermissionModel      `tfsdk:"audit_log"`
// 	Boundary         BoundaryModel        `tfsdk:"boundary"`
// 	BusinessImpact   ManageAndReadModel   `tfsdk:"business_impact"`
// 	Collector        ManageAndReadModel   `tfsdk:"collector"`
// 	CustomProperties ManageAndReadModel   `tfsdk:"custom_properties"`
// 	Integration      ManageAndReadModel   `tfsdk:"integration"`
// 	InternalIps      ManageAndReadModel   `tfsdk:"internal_ips"`
// 	Notifications    ManageAndReadModel   `tfsdk:"notifications"`
// 	OIDC             ManageAndReadModel   `tfsdk:"oidc"`
// 	SAML             ManageAndReadModel   `tfsdk:"saml"`
// 	SecretKey        PermissionModel      `tfsdk:"secret_key"`
// 	SecuritySettings PermissionModel      `tfsdk:"security_settings"`
// 	SitesAndSensors  SitesAndSensorsModel `tfsdk:"sites_and_sensors"`
// 	UsersAndRoles    UsersAndRolesModel   `tfsdk:"users_and_roles"`
// }
//
// // BoundaryModel maps the Boundary schema data.
// type BoundaryModel struct {
// 	All    types.Bool          `tfsdk:"all"`
// 	Manage ManageBoundaryModel `tfsdk:"manage"`
// 	Read   PermissionModel     `tfsdk:"read"`
// }
//
// // ManageBoundaryModel maps the ManageBoundary schema data.
// type ManageBoundaryModel struct {
// 	All    types.Bool      `tfsdk:"all"`
// 	Create PermissionModel `tfsdk:"create"`
// 	Delete PermissionModel `tfsdk:"delete"`
// 	Edit   PermissionModel `tfsdk:"edit"`
// }
//
// // ManageAndReadModel maps the ManageAndRead schema data.
// type ManageAndReadModel struct {
// 	All    types.Bool      `tfsdk:"all"`
// 	Manage PermissionModel `tfsdk:"manage"`
// 	Read   PermissionModel `tfsdk:"read"`
// }
//
// // SitesAndSensorsModel maps the SitesAndSensors schema data.
// type SitesAndSensorsModel struct {
// 	All    types.Bool         `tfsdk:"all"`
// 	Manage ManageSensorsModel `tfsdk:"manage"`
// 	Read   PermissionModel    `tfsdk:"read"`
// }
//
// // ManageSensorsModel maps the ManageSensors schema data.
// type ManageSensorsModel struct {
// 	All     types.Bool      `tfsdk:"all"`
// 	Sensors PermissionModel `tfsdk:"sensors"`
// 	Sites   PermissionModel `tfsdk:"sites"`
// }
//
// // UsersAndRolesModel maps the UsersAndRoles schema data.
// type UsersAndRolesModel struct {
// 	All    types.Bool       `tfsdk:"all"`
// 	Manage ManageUsersModel `tfsdk:"manage"`
// 	Read   PermissionModel  `tfsdk:"read"`
// }
//
// // ManageUsersModel maps the ManageUsers schema data.
// type ManageUsersModel struct {
// 	All   types.Bool       `tfsdk:"all"`
// 	Roles UserActionsModel `tfsdk:"roles"`
// 	Users UserActionsModel `tfsdk:"users"`
// }
//
// // UserActionsModel maps the UserActions schema data.
// type UserActionsModel struct {
// 	All    types.Bool      `tfsdk:"all"`
// 	Create PermissionModel `tfsdk:"create"`
// 	Delete PermissionModel `tfsdk:"delete"`
// 	Edit   PermissionModel `tfsdk:"edit"`
// }
//
// // UserModel maps the User schema data.
// type UserModel struct {
// 	All    types.Bool      `tfsdk:"all"`
// 	Manage ManageUserModel `tfsdk:"manage"`
// 	Read   PermissionModel `tfsdk:"read"`
// }
//
// // ManageUserModel maps the ManageUser schema data.
// type ManageUserModel struct {
// 	All    types.Bool      `tfsdk:"all"`
// 	Upsert PermissionModel `tfsdk:"upsert"`
// }
//
// // VulnerabilityModel maps the Vulnerability schema data.
// type VulnerabilityModel struct {
// 	All    types.Bool      `tfsdk:"all"`
// 	Manage ManageVulnModel `tfsdk:"manage"`
// 	Read   PermissionModel `tfsdk:"read"`
// }
//
// // ManageVulnModel maps the ManageVuln schema data.
// type ManageVulnModel struct {
// 	All     types.Bool      `tfsdk:"all"`
// 	Ignore  PermissionModel `tfsdk:"ignore"`
// 	Resolve PermissionModel `tfsdk:"resolve"`
// 	Write   PermissionModel `tfsdk:"write"`
// }

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
	success, err := r.client.CreateRole(role)
	if err != nil || !success {
		resp.Diagnostics.AddError(
			"Error creating role",
			fmt.Sprintf("Failed to create role %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	// Fetch the created role to get its ID and other attributes
	createdRole, err := r.client.GetRoleByName(plan.Name.ValueString())
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
	role, err := r.client.GetRoleByID(state.ID.ValueString())
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
	_, err := r.client.UpdateRole(role, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Role",
			fmt.Sprintf("Failed to update role with ID %q: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Fetch the updated role details
	updatedRole, err := r.client.GetRoleByID(state.ID.ValueString())
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
	success, err := r.client.DeleteRole(state.ID.ValueString())
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
