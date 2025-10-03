// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	armis "github.com/1898andCo/terraform-provider-armis-centrix/armis"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

type userResource struct {
	client *armis.Client
}

func UserResource() resource.Resource {
	return &userResource{}
}

// Configure adds the provider configured client to the resource.
func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*armis.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *armis.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the user resource.
func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides an Armis user

The resource provisions a user with the ability to define location, email, roles, and role assignments.
`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:      true,
				Description:   "The full name of the user.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured(), stringplanmodifier.UseStateForUnknown()},
			},
			"phone": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The phone number of the user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "The email address of the user.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9._-]*[a-zA-Z0-9])?@[a-zA-Z0-9]([a-zA-Z0-9.-]*[a-zA-Z0-9])?\.[a-zA-Z]{2,}$`),
						"must be a valid email address",
					),
				},
			},
			"location": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The physical location or address of the user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The job title or designation of the user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "The unique username of the user.",
			},
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "A unique identifier for the user resource.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"role_assignments": schema.ListNestedAttribute{
				Required:    true,
				Description: "Role assignments for the user.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "The names of the roles assigned to the user.",
						},
						"sites": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "A list of site identifiers associated with the role.",
						},
					},
				},
			},
		},
	}
}

// userResourceModel maps the resource schema data.
type userResourceModel struct {
	ID              types.String      `tfsdk:"id"`
	Name            types.String      `tfsdk:"name"`
	Phone           types.String      `tfsdk:"phone"`
	Email           types.String      `tfsdk:"email"`
	Location        types.String      `tfsdk:"location"`
	Title           types.String      `tfsdk:"title"`
	Username        types.String      `tfsdk:"username"`
	RoleAssignments []RoleAssignments `tfsdk:"role_assignments"`
}

type RoleAssignments struct {
	Name  []types.String `tfsdk:"name"`
	Sites []types.String `tfsdk:"sites"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel

	// Parse the plan from Terraform
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := buildArmisUser(plan)
	tflog.Info(ctx, "Creating user in Armis", map[string]any{"user": user})

	newUser, err := r.client.CreateUser(ctx, user)
	if err != nil {
		appendAPIError(&resp.Diagnostics, fmt.Sprintf("Error creating user %q", plan.Username.ValueString()), err)
		return
	}

	tflog.Info(ctx, "API returned user data after creation", map[string]any{
		"user_id":                   newUser.ID,
		"role_assignments_returned": fmt.Sprintf("%+v", newUser.RoleAssignment),
		"role_assignments_count":    len(newUser.RoleAssignment),
	})

	if len(newUser.RoleAssignment) > 0 {
		for i, ra := range newUser.RoleAssignment {
			tflog.Info(ctx, "Role assignment details", map[string]any{
				"index": i,
				"names": ra.Name,
				"sites": ra.Sites,
			})
		}
	}

	plan.ID = types.StringValue(strconv.Itoa(newUser.ID))
	plan.Name = types.StringValue(newUser.Name)
	plan.Phone = types.StringValue(newUser.Phone)
	plan.Email = types.StringValue(newUser.Email)
	plan.Location = types.StringValue(newUser.Location)
	plan.Title = types.StringValue(newUser.Title)
	plan.Username = types.StringValue(newUser.Username)

	// Keep the planned role assignments - don't overwrite with API response as sites aren't returned
	// The API may return them in different order or incomplete
	// plan.RoleAssignments stays as originally planned

	// Save the state
	tflog.Info(ctx, "Setting state for user")
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read user resource information.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state userResourceModel
	tflog.Info(ctx, "Retrieving current state")
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed user value from Armis
	user, err := r.client.GetUser(ctx, state.ID.ValueString())
	if err != nil {
		// Handle 404 Not Found by removing resource from state
		if strings.Contains(err.Error(), "Status Code: 404") {
			tflog.Warn(ctx, "User not found, removing from state", map[string]any{
				"user_id": state.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}

		appendAPIError(&resp.Diagnostics, fmt.Sprintf("Error reading user %s", state.ID.ValueString()), err)
		return
	}

	if user == nil {
		resp.State.RemoveResource(ctx)
		tflog.Warn(ctx, "User not found, removing from state", map[string]any{
			"user_id": state.ID.ValueString(),
		})
		return
	}

	state.Name = types.StringValue(user.Name)
	state.Phone = types.StringValue(user.Phone)
	state.Email = types.StringValue(user.Email)
	state.Location = types.StringValue(user.Location)
	state.Title = types.StringValue(user.Title)
	state.Username = types.StringValue(user.Username)

	// state.RoleAssignments - DON'T UPDATE, keep existing state
	// Only update the basic fields - leave role_assignments alone
	// The API returns them in inconsistent order, so we trust our state

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the current state to get the user ID
	var state userResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that the user ID is available
	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Error Updating User",
			"The user ID is missing from the state. This is required to update the user.",
		)
		return
	}

	user := buildArmisUser(plan)

	// Update user
	_, err := r.client.UpdateUser(ctx, user, state.ID.ValueString())
	if err != nil {
		appendAPIError(&resp.Diagnostics, fmt.Sprintf("Error updating user %s", state.ID.ValueString()), err)
		return
	}

	// For update, only fetch and update the basic fields
	// Keep the role assignments from the plan since we just sent them
	updatedUser, err := r.client.GetUser(ctx, state.ID.ValueString())
	if err != nil {
		appendAPIError(&resp.Diagnostics, fmt.Sprintf("Error reading user %s", state.ID.ValueString()), err)
		return
	}

	// Map only the basic fields from API response
	plan.ID = types.StringValue(strconv.Itoa(updatedUser.ID))
	plan.Name = types.StringValue(updatedUser.Name)
	plan.Phone = types.StringValue(updatedUser.Phone)
	plan.Email = types.StringValue(updatedUser.Email)
	plan.Location = types.StringValue(updatedUser.Location)
	plan.Title = types.StringValue(updatedUser.Title)
	plan.Username = types.StringValue(updatedUser.Username)

	// Keep role assignments from plan - don't overwrite with API response as it does not return sites
	// plan.RoleAssignments stays as planned

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing user
	success, err := r.client.DeleteUser(ctx, state.ID.ValueString())
	if err != nil {
		appendAPIError(&resp.Diagnostics, fmt.Sprintf("Error deleting user %s", state.ID.ValueString()), err)
		return
	}
	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Armis User",
			"Could not delete user: operation returned unsuccessful status",
		)
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func buildArmisUser(plan userResourceModel) armis.UserSettings {
	var roleAssignments []armis.RoleAssignment
	for _, roleAssignment := range plan.RoleAssignments {
		// Simple conversion from []types.String to []string
		var names []string
		for _, n := range roleAssignment.Name {
			names = append(names, n.ValueString())
		}

		var sites []string
		for _, s := range roleAssignment.Sites {
			sites = append(sites, s.ValueString())
		}

		roleAssignments = append(roleAssignments, armis.RoleAssignment{
			Name:  names,
			Sites: sites,
		})
	}

	return armis.UserSettings{
		Name:           plan.Name.ValueString(),
		Phone:          plan.Phone.ValueString(),
		Email:          plan.Email.ValueString(),
		Location:       plan.Location.ValueString(),
		Title:          plan.Title.ValueString(),
		Username:       plan.Username.ValueString(),
		RoleAssignment: roleAssignments,
	}
}
