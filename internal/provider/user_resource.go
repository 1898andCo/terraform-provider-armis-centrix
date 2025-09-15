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
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
			"Unexpected Data Source Configure Type",
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
			"role_assignments": schema.SingleNestedAttribute{
				Required:    true,
				Description: "A list of role assignments for the user.",
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
	}
}

// userResourceModel maps the resource schema data.
type userResourceModel struct {
	ID              types.String   `tfsdk:"id"`
	Name            types.String   `tfsdk:"name"`
	Phone           types.String   `tfsdk:"phone"`
	Email           types.String   `tfsdk:"email"`
	Location        types.String   `tfsdk:"location"`
	Title           types.String   `tfsdk:"title"`
	Username        types.String   `tfsdk:"username"`
	RoleAssignments roleAssignment `tfsdk:"role_assignments"` // value, not pointer
}

type roleAssignment struct {
	Name  types.List `tfsdk:"name"`  // list(string)
	Sites types.List `tfsdk:"sites"` // list(string)
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	tflog.Info(ctx, "Creating user")

	// Parse the plan from Terraform
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build API payload
	user := armis.UserSettings{
		Name:     plan.Name.ValueString(),
		Phone:    plan.Phone.ValueString(),
		Email:    plan.Email.ValueString(),
		Location: plan.Location.ValueString(),
		Title:    plan.Title.ValueString(),
		Username: plan.Username.ValueString(),
	}

	// role assignments (one nested object in API)
	if ra := buildAPIRoleAssignmentFromModel(plan.RoleAssignments); ra != nil {
		user.RoleAssignment = []armis.RoleAssignment{*ra}
	}

	// Create the user via the client
	newUser, err := r.client.CreateUser(ctx, user)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	// Map the response to Terraform state
	plan.ID = types.StringValue(strconv.Itoa(newUser.ID))
	plan.Name = types.StringValue(newUser.Name)
	plan.Phone = types.StringValue(newUser.Phone)
	plan.Email = types.StringValue(newUser.Email)
	plan.Location = types.StringValue(newUser.Location)
	plan.Title = types.StringValue(newUser.Title)
	plan.Username = types.StringValue(newUser.Username)
	plan.RoleAssignments = mapAPIRoleAssignmentsToModel(newUser.RoleAssignment)

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

		resp.Diagnostics.AddError(
			"Error Reading Armis User",
			fmt.Sprintf("Failed to fetch user: %v", err),
		)
		return
	}

	if user == nil {
		resp.State.RemoveResource(ctx)
		tflog.Warn(ctx, "User not found, removing from state", map[string]any{
			"user_id": state.ID.ValueString(),
		})
		return
	}

	// Overwrite users with refreshed state
	state.Name = types.StringValue(user.Name)
	state.Phone = types.StringValue(user.Phone)
	state.Email = types.StringValue(user.Email)
	state.Location = types.StringValue(user.Location)
	state.Title = types.StringValue(user.Title)
	state.Username = types.StringValue(user.Username)
	state.RoleAssignments = mapAPIRoleAssignmentsToModel(user.RoleAssignment)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

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

	// Build API payload
	user := armis.UserSettings{
		Name:     plan.Name.ValueString(),
		Phone:    plan.Phone.ValueString(),
		Email:    plan.Email.ValueString(),
		Location: plan.Location.ValueString(),
		Title:    plan.Title.ValueString(),
		Username: plan.Username.ValueString(),
	}
	if ra := buildAPIRoleAssignmentFromModel(plan.RoleAssignments); ra != nil {
		user.RoleAssignment = []armis.RoleAssignment{*ra}
	}

	// Update existing user and then fetch the updated user from the API.
	_, err := r.client.UpdateUser(ctx, user, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Armis user",
			"Could not update user, unexpected error: "+err.Error(),
		)
		return
	}

	updatedUser, err := r.client.GetUser(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Armis User",
			"Could not read Armis user ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map the response to Terraform state
	plan.ID = types.StringValue(strconv.Itoa(updatedUser.ID))
	plan.Name = types.StringValue(updatedUser.Name)
	plan.Phone = types.StringValue(updatedUser.Phone)
	plan.Email = types.StringValue(updatedUser.Email)
	plan.Location = types.StringValue(updatedUser.Location)
	plan.Title = types.StringValue(updatedUser.Title)
	plan.Username = types.StringValue(updatedUser.Username)
	plan.RoleAssignments = mapAPIRoleAssignmentsToModel(updatedUser.RoleAssignment)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

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
	if err != nil || !success {
		resp.Diagnostics.AddError(
			"Error Deleting Armis user",
			"Could not delete user, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *userResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


func buildAPIRoleAssignmentFromModel(m roleAssignment) *armis.RoleAssignment {
	names := listToStringSlice(m.Name)
	sites := listToStringSlice(m.Sites)

	if len(names) == 0 && len(sites) == 0 {
		return nil
	}

	return &armis.RoleAssignment{
		Name:  names,
		Sites: sites,
	}
}

func mapAPIRoleAssignmentsToModel(api []armis.RoleAssignment) roleAssignment {
	if len(api) == 0 {
		return roleAssignment{
			Name:  types.ListNull(types.StringType),
			Sites: types.ListNull(types.StringType),
		}
	}
	first := api[0]
	return roleAssignment{
		Name:  stringSliceToList(first.Name),
		Sites: stringSliceToList(first.Sites),
	}
}

func listToStringSlice(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	out := make([]string, 0, len(list.Elements()))
	for _, v := range list.Elements() {
		if s, ok := v.(types.String); ok && !s.IsNull() && !s.IsUnknown() {
			out = append(out, s.ValueString())
		}
	}
	return out
}

func stringSliceToList(slice []string) types.List {
	if slice == nil {
		return types.ListNull(types.StringType)
	}
	vals := make([]attr.Value, len(slice))
	for i, s := range slice {
		vals[i] = types.StringValue(s)
	}
	lv, _ := types.ListValue(types.StringType, vals)
	return lv
}
