// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/1898andCo/armis-sdk-go/armis"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// UserDataSource is a helper function to simplify the provider implementation.
func UserDataSource() datasource.DataSource {
	return &usersDataSource{}
}

// usersDataSource is the data source implementation.
type usersDataSource struct {
	client *armis.Client
}

// Metadata returns the data source type name.
func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the users data source.
func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Armis user information. If an email is provided, the data source will retrieve information about the specified user",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Description: "An optional email address used to filter the retrieved user information. If specified, only the user matching this email will be returned.",
				Optional:    true,
			},
			"users": schema.ListNestedAttribute{
				Description: "A computed list of users. Each object in the list contains detailed information about a user.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "A unique identifier for the user.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The full name of the user.",
							Computed:    true,
						},
						"email": schema.StringAttribute{
							Description: "The email address of the user.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Indicates whether the user account is active.",
							Computed:    true,
						},
						"last_login_time": schema.StringAttribute{
							Description: "The timestamp of the user's last login, formatted as a string.",
							Computed:    true,
						},
						"location": schema.StringAttribute{
							Description: "The physical location or address of the user.",
							Computed:    true,
						},
						"phone": schema.StringAttribute{
							Description: "The user's phone number.",
							Computed:    true,
						},
						"pov_eula_signing_date": schema.StringAttribute{
							Description: "The date when the user signed the Proof of Value (PoV) End User License Agreement (EULA).",
							Computed:    true,
						},
						"prod_eula_signing_date": schema.StringAttribute{
							Description: "The date when the user signed the Production End User License Agreement (EULA).",
							Computed:    true,
						},
						"report_permissions": schema.StringAttribute{
							Description: "The level of permissions the user has for generating or accessing reports.",
							Computed:    true,
						},
						"role": schema.StringAttribute{
							Description: "The user's role within the system.",
							Computed:    true,
						},
						"role_assignment": schema.ListNestedAttribute{
							Description: "A list of role assignments for the user. Each object contains details about a specific role.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.ListAttribute{
										Description: "The names of the roles assigned to the user.",
										ElementType: types.StringType,
										Computed:    true,
									},
								},
							},
						},
						"title": schema.StringAttribute{
							Description: "The job title of the user.",
							Computed:    true,
						},
						"two_factor_authentication": schema.BoolAttribute{
							Description: "Indicates whether the user has two-factor authentication enabled.",
							Computed:    true,
						},
						"username": schema.StringAttribute{
							Description: "The username associated with the user's account.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// usersDataSourceModel maps the data source schema data.
type usersDataSourceModel struct {
	Email types.String `tfsdk:"email"`
	Users []userModel  `tfsdk:"users"`
}

// userModel maps the user schema data.
type userModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Email                   types.String `tfsdk:"email"`
	IsActive                types.Bool   `tfsdk:"is_active"`
	LastLoginTime           types.String `tfsdk:"last_login_time"`
	Location                types.String `tfsdk:"location"`
	Phone                   types.String `tfsdk:"phone"`
	PovEULASigningDate      types.String `tfsdk:"pov_eula_signing_date"`
	ProdEULASigningDate     types.String `tfsdk:"prod_eula_signing_date"`
	ReportPermissions       types.String `tfsdk:"report_permissions"`
	Role                    types.String `tfsdk:"role"`
	RoleAssignment          []roleModel  `tfsdk:"role_assignment"`
	Title                   types.String `tfsdk:"title"`
	TwoFactorAuthentication types.Bool   `tfsdk:"two_factor_authentication"`
	Username                types.String `tfsdk:"username"`
}

// roleModel maps the role assignment schema data.
type roleModel struct {
	Name []types.String `tfsdk:"name"`
}

// Read refreshes the Terraform state with the latest data.
func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config usersDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var users []userModel

	if !config.Email.IsNull() {
		// Fetch a specific user by email
		user, err := d.client.GetUser(ctx, config.Email.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Armis User",
				err.Error(),
			)
			return
		}

		// Map response body to model
		userState := userModel{
			ID:                      types.StringValue(fmt.Sprintf("%d", user.ID)),
			Name:                    types.StringValue(user.Name),
			Email:                   types.StringValue(user.Email),
			IsActive:                types.BoolValue(user.IsActive),
			Location:                types.StringValue(user.Location),
			Phone:                   types.StringValue(user.Phone),
			Title:                   types.StringValue(user.Title),
			Username:                types.StringValue(user.Username),
			Role:                    types.StringValue(user.Role),
			ReportPermissions:       types.StringValue(user.ReportPermissions),
			TwoFactorAuthentication: types.BoolValue(user.TwoFactorAuthentication),
			LastLoginTime:           types.StringValue(user.LastLoginTime),
			PovEULASigningDate:      types.StringValue(user.PovEULASigningDate),
			ProdEULASigningDate:     types.StringValue(user.ProdEULASigningDate),
		}

		// Map role assignments
		for _, role := range user.RoleAssignment {
			var roleNames []types.String
			for _, name := range role.Name {
				roleNames = append(roleNames, types.StringValue(name))
			}
			userState.RoleAssignment = append(userState.RoleAssignment, roleModel{
				Name: roleNames,
			})
		}

		users = append(users, userState)
	} else {
		// Fetch all users
		allUsers, err := d.client.GetUsers(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Armis Users",
				err.Error(),
			)
			return
		}

		// Map response body to model
		for _, user := range allUsers {
			userState := userModel{
				ID:                      types.StringValue(fmt.Sprintf("%d", user.ID)),
				Name:                    types.StringValue(user.Name),
				Email:                   types.StringValue(user.Email),
				IsActive:                types.BoolValue(user.IsActive),
				Location:                types.StringValue(user.Location),
				Phone:                   types.StringValue(user.Phone),
				Title:                   types.StringValue(user.Title),
				Username:                types.StringValue(user.Username),
				Role:                    types.StringValue(user.Role),
				ReportPermissions:       types.StringValue(user.ReportPermissions),
				TwoFactorAuthentication: types.BoolValue(user.TwoFactorAuthentication),
				LastLoginTime:           types.StringValue(user.LastLoginTime),
				PovEULASigningDate:      types.StringValue(user.PovEULASigningDate),
				ProdEULASigningDate:     types.StringValue(user.ProdEULASigningDate),
			}

			// Map role assignments
			for _, role := range user.RoleAssignment {
				var roleNames []types.String
				for _, name := range role.Name {
					roleNames = append(roleNames, types.StringValue(name))
				}
				userState.RoleAssignment = append(userState.RoleAssignment, roleModel{
					Name: roleNames,
				})
			}

			users = append(users, userState)
		}
	}

	// Save data into Terraform state
	config.Users = users
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
