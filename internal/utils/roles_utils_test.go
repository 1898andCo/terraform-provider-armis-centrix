// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/1898andCo/terraform-provider-armis-centrix/armis"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestShouldIncludeRole tests the ShouldIncludeRole filter function.
func TestShouldIncludeRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		model    RoleDataSourceSummaryModel
		prefix   types.String
		expected bool
	}{
		{
			name: "null prefix should include all",
			model: RoleDataSourceSummaryModel{
				Name: types.StringValue("AdminRole"),
			},
			prefix:   types.StringNull(),
			expected: true,
		},
		{
			name: "unknown prefix should include all",
			model: RoleDataSourceSummaryModel{
				Name: types.StringValue("AdminRole"),
			},
			prefix:   types.StringUnknown(),
			expected: true,
		},
		{
			name: "empty prefix should include all",
			model: RoleDataSourceSummaryModel{
				Name: types.StringValue("AdminRole"),
			},
			prefix:   types.StringValue(""),
			expected: true,
		},
		{
			name: "matching prefix should include",
			model: RoleDataSourceSummaryModel{
				Name: types.StringValue("AdminRole"),
			},
			prefix:   types.StringValue("Admin"),
			expected: true,
		},
		{
			name: "non-matching prefix should not include",
			model: RoleDataSourceSummaryModel{
				Name: types.StringValue("UserRole"),
			},
			prefix:   types.StringValue("Admin"),
			expected: false,
		},
		{
			name: "exact match should include",
			model: RoleDataSourceSummaryModel{
				Name: types.StringValue("ReadOnly"),
			},
			prefix:   types.StringValue("ReadOnly"),
			expected: true,
		},
		{
			name: "null model name should not include",
			model: RoleDataSourceSummaryModel{
				Name: types.StringNull(),
			},
			prefix:   types.StringValue("Admin"),
			expected: false,
		},
		{
			name: "unknown model name should not include",
			model: RoleDataSourceSummaryModel{
				Name: types.StringUnknown(),
			},
			prefix:   types.StringValue("Admin"),
			expected: false,
		},
		{
			name: "case sensitive matching",
			model: RoleDataSourceSummaryModel{
				Name: types.StringValue("adminRole"),
			},
			prefix:   types.StringValue("Admin"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ShouldIncludeRole(tt.model, tt.prefix)
			if result != tt.expected {
				t.Errorf("ShouldIncludeRole() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestShouldExcludeRole tests the ShouldExcludeRole filter function.
func TestShouldExcludeRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		model    RoleDataSourceSummaryModel
		prefix   types.String
		expected bool
	}{
		{
			name: "null prefix should not exclude",
			model: RoleDataSourceSummaryModel{
				Name: types.StringValue("AdminRole"),
			},
			prefix:   types.StringNull(),
			expected: false,
		},
		{
			name: "unknown prefix should not exclude",
			model: RoleDataSourceSummaryModel{
				Name: types.StringValue("AdminRole"),
			},
			prefix:   types.StringUnknown(),
			expected: false,
		},
		{
			name: "empty prefix should not exclude",
			model: RoleDataSourceSummaryModel{
				Name: types.StringValue("AdminRole"),
			},
			prefix:   types.StringValue(""),
			expected: false,
		},
		{
			name: "matching prefix should exclude",
			model: RoleDataSourceSummaryModel{
				Name: types.StringValue("TestRole"),
			},
			prefix:   types.StringValue("Test"),
			expected: true,
		},
		{
			name: "non-matching prefix should not exclude",
			model: RoleDataSourceSummaryModel{
				Name: types.StringValue("UserRole"),
			},
			prefix:   types.StringValue("Admin"),
			expected: false,
		},
		{
			name: "exact match should exclude",
			model: RoleDataSourceSummaryModel{
				Name: types.StringValue("System"),
			},
			prefix:   types.StringValue("System"),
			expected: true,
		},
		{
			name: "null model name should not exclude",
			model: RoleDataSourceSummaryModel{
				Name: types.StringNull(),
			},
			prefix:   types.StringValue("Test"),
			expected: false,
		},
		{
			name: "unknown model name should not exclude",
			model: RoleDataSourceSummaryModel{
				Name: types.StringUnknown(),
			},
			prefix:   types.StringValue("Test"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ShouldExcludeRole(tt.model, tt.prefix)
			if result != tt.expected {
				t.Errorf("ShouldExcludeRole() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestBuildRoleDataSourceSummaryModel tests the BuildRoleDataSourceSummaryModel function.
func TestBuildRoleDataSourceSummaryModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    armis.RoleSettings
		expected RoleDataSourceSummaryModel
	}{
		{
			name: "basic role conversion",
			input: armis.RoleSettings{
				ID:       123,
				Name:     "AdminRole",
				ViprRole: true,
			},
			expected: RoleDataSourceSummaryModel{
				ID:       types.StringValue("123"),
				Name:     types.StringValue("AdminRole"),
				ViprRole: types.BoolValue(true),
			},
		},
		{
			name: "role with false vipr",
			input: armis.RoleSettings{
				ID:       456,
				Name:     "UserRole",
				ViprRole: false,
			},
			expected: RoleDataSourceSummaryModel{
				ID:       types.StringValue("456"),
				Name:     types.StringValue("UserRole"),
				ViprRole: types.BoolValue(false),
			},
		},
		{
			name: "role with empty name",
			input: armis.RoleSettings{
				ID:       789,
				Name:     "",
				ViprRole: false,
			},
			expected: RoleDataSourceSummaryModel{
				ID:       types.StringValue("789"),
				Name:     types.StringValue(""),
				ViprRole: types.BoolValue(false),
			},
		},
		{
			name: "role with special characters in name",
			input: armis.RoleSettings{
				ID:       999,
				Name:     "Role-With_Special.Chars",
				ViprRole: true,
			},
			expected: RoleDataSourceSummaryModel{
				ID:       types.StringValue("999"),
				Name:     types.StringValue("Role-With_Special.Chars"),
				ViprRole: types.BoolValue(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := BuildRoleDataSourceSummaryModel(&tt.input)

			// Compare ID
			if result.ID.ValueString() != tt.expected.ID.ValueString() {
				t.Errorf("BuildRoleDataSourceSummaryModel().ID = %v, want %v",
					result.ID.ValueString(), tt.expected.ID.ValueString())
			}

			// Compare Name
			if result.Name.ValueString() != tt.expected.Name.ValueString() {
				t.Errorf("BuildRoleDataSourceSummaryModel().Name = %v, want %v",
					result.Name.ValueString(), tt.expected.Name.ValueString())
			}

			// Compare ViprRole
			if result.ViprRole.ValueBool() != tt.expected.ViprRole.ValueBool() {
				t.Errorf("BuildRoleDataSourceSummaryModel().ViprRole = %v, want %v",
					result.ViprRole.ValueBool(), tt.expected.ViprRole.ValueBool())
			}
		})
	}
}

// createMinimalRoleModel creates a RoleResourceModel with all nested pointers initialized.
func createMinimalRoleModel() RoleResourceModel {
	return RoleResourceModel{
		Name: types.StringValue(""),
		Permissions: &PermissionsModel{
			AdvancedPermissions: &AdvancedPermissionsModel{
				All:        types.BoolValue(false),
				Behavioral: &BehavioralModel{},
				Device:     &DeviceAdvancedModel{},
			},
			Alert: &AlertModel{
				All:    types.BoolValue(false),
				Manage: &ManageAlertsModel{},
				Read:   types.BoolValue(false),
			},
			Device: &DeviceModel{
				All: types.BoolValue(false),
				Manage: &ManageDeviceModel{
					All: types.BoolValue(false),
					Enforce: &EnforceModel{
						All: types.BoolValue(false),
					},
				},
				Read: types.BoolValue(false),
			},
			Policy: &PolicyModel{
				All: types.BoolValue(false),
			},
			Report: &ReportModel{
				All:    types.BoolValue(false),
				Manage: &ManageReportModel{},
				Read:   types.BoolValue(false),
			},
			RiskFactor: &RiskFactorModel{
				All: types.BoolValue(false),
				Manage: &ManageRiskModel{
					All:           types.BoolValue(false),
					Customization: &CustomizationModel{},
					Status:        &StatusModel{},
				},
				Read: types.BoolValue(false),
			},
			Settings: &SettingsModel{
				All: types.BoolValue(false),
				Boundary: &BoundaryModel{
					All:    types.BoolValue(false),
					Manage: &ManageBoundaryModel{},
					Read:   types.BoolValue(false),
				},
				BusinessImpact:   &ManageAndReadModel{},
				Collector:        &ManageAndReadModel{},
				CustomProperties: &ManageAndReadModel{},
				Integration:      &ManageAndReadModel{},
				InternalIps:      &ManageAndReadModel{},
				Notifications:    &ManageAndReadModel{},
				OIDC:             &ManageAndReadModel{},
				SAML:             &ManageAndReadModel{},
				SitesAndSensors: &SitesAndSensorsModel{
					All: types.BoolValue(false),
					Manage: &ManageSitesAndSensorsModel{
						All: types.BoolValue(false),
					},
					Read: types.BoolValue(false),
				},
				UsersAndRoles: &UsersAndRolesModel{
					All: types.BoolValue(false),
					Manage: &ManageUsersAndRolesModel{
						All: types.BoolValue(false),
						Roles: &ManageRolesModel{
							All: types.BoolValue(false),
						},
						Users: &ManageUsersModel{
							All: types.BoolValue(false),
						},
					},
					Read: types.BoolValue(false),
				},
			},
			User: &UserModel{
				All: types.BoolValue(false),
				Manage: &ManageUserModel{
					All: types.BoolValue(false),
				},
				Read: types.BoolValue(false),
			},
			Vulnerability: &VulnerabilityModel{
				All: types.BoolValue(false),
				Manage: &ManageVulnerabilityModel{
					All: types.BoolValue(false),
				},
				Read: types.BoolValue(false),
			},
		},
	}
}

// TestBuildRoleRequest_BasicPermissions tests basic permission conversion.
func TestBuildRoleRequest_BasicPermissions(t *testing.T) {
	t.Parallel()

	t.Run("role with advanced permissions", func(t *testing.T) {
		t.Parallel()

		model := createMinimalRoleModel()
		model.Name = types.StringValue("TestRole")
		model.Permissions.AdvancedPermissions.All = types.BoolValue(true)
		model.Permissions.AdvancedPermissions.Behavioral.ApplicationName = types.BoolValue(true)

		result := BuildRoleRequest(model)

		if result.Name != "TestRole" {
			t.Errorf("Expected Name 'TestRole', got '%s'", result.Name)
		}
		if !result.Permissions.AdvancedPermissions.All {
			t.Error("Expected AdvancedPermissions.All to be true")
		}
		if !result.Permissions.AdvancedPermissions.Behavioral.ApplicationName.All {
			t.Error("Expected ApplicationName.All to be true")
		}
	})

	t.Run("role with alert permissions", func(t *testing.T) {
		t.Parallel()

		model := createMinimalRoleModel()
		model.Name = types.StringValue("AlertRole")
		model.Permissions.Alert.All = types.BoolValue(true)
		model.Permissions.Alert.Manage.Resolve = types.BoolValue(true)
		model.Permissions.Alert.Read = types.BoolValue(true)

		result := BuildRoleRequest(model)

		if result.Name != "AlertRole" {
			t.Errorf("Expected Name 'AlertRole', got '%s'", result.Name)
		}
		if !result.Permissions.Alert.All {
			t.Error("Expected Alert.All to be true")
		}
		if !result.Permissions.Alert.Manage.Resolve.All {
			t.Error("Expected Alert.Manage.Resolve.All to be true")
		}
		if !result.Permissions.Alert.Read.All {
			t.Error("Expected Alert.Read.All to be true")
		}
	})
}

// TestBuildRoleResourceModel_BasicConversion tests reverse conversion from API to Terraform model.
func TestBuildRoleResourceModel_BasicConversion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    armis.RoleSettings
		validate func(t *testing.T, result RoleResourceModel)
	}{
		{
			name: "basic role with name and id",
			input: armis.RoleSettings{
				ID:   123,
				Name: "TestRole",
				Permissions: armis.Permissions{
					AdvancedPermissions: armis.AdvancedPermissions{
						All: true,
						Behavioral: armis.Behavioral{
							All: false,
							ApplicationName: armis.Permission{
								All: true,
							},
							HostName: armis.Permission{
								All: false,
							},
							ServiceName: armis.Permission{
								All: true,
							},
						},
					},
				},
			},
			validate: func(t *testing.T, result RoleResourceModel) {
				if result.Name.ValueString() != "TestRole" {
					t.Errorf("Expected Name 'TestRole', got '%s'", result.Name.ValueString())
				}
				if result.ID.ValueString() != "123" {
					t.Errorf("Expected ID '123', got '%s'", result.ID.ValueString())
				}
				if result.Permissions == nil {
					t.Fatal("Expected Permissions to be non-nil")
				}
				if result.Permissions.AdvancedPermissions == nil {
					t.Fatal("Expected AdvancedPermissions to be non-nil")
				}
				if !result.Permissions.AdvancedPermissions.All.ValueBool() {
					t.Error("Expected AdvancedPermissions.All to be true")
				}
				if result.Permissions.AdvancedPermissions.Behavioral == nil {
					t.Fatal("Expected Behavioral to be non-nil")
				}
				if result.Permissions.AdvancedPermissions.Behavioral.All.ValueBool() {
					t.Error("Expected Behavioral.All to be false")
				}
				if !result.Permissions.AdvancedPermissions.Behavioral.ApplicationName.ValueBool() {
					t.Error("Expected ApplicationName to be true")
				}
			},
		},
		{
			name: "role with alert permissions",
			input: armis.RoleSettings{
				ID:   456,
				Name: "AlertManager",
				Permissions: armis.Permissions{
					Alert: armis.Alert{
						All: true,
						Manage: armis.Manage{
							All: false,
							Resolve: armis.Permission{
								All: true,
							},
							WhitelistDevices: armis.Permission{
								All: false,
							},
						},
						Read: armis.Permission{
							All: true,
						},
					},
				},
			},
			validate: func(t *testing.T, result RoleResourceModel) {
				if result.Name.ValueString() != "AlertManager" {
					t.Errorf("Expected Name 'AlertManager', got '%s'", result.Name.ValueString())
				}
				if result.ID.ValueString() != "456" {
					t.Errorf("Expected ID '456', got '%s'", result.ID.ValueString())
				}
				if result.Permissions == nil || result.Permissions.Alert == nil {
					t.Fatal("Expected Alert permissions to be non-nil")
				}
				if !result.Permissions.Alert.All.ValueBool() {
					t.Error("Expected Alert.All to be true")
				}
				if result.Permissions.Alert.Manage == nil {
					t.Fatal("Expected Alert.Manage to be non-nil")
				}
				if !result.Permissions.Alert.Read.ValueBool() {
					t.Error("Expected Alert.Read to be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Create empty model for BuildRoleResourceModel to populate
			emptyModel := RoleResourceModel{}
			result := BuildRoleResourceModel(&tt.input, emptyModel)
			tt.validate(t, result)
		})
	}
}

// TestBuildRoleDataSourceModel tests conversion from API to data source model.
func TestBuildRoleDataSourceModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    armis.RoleSettings
		validate func(t *testing.T, result RoleDataSourceModel)
	}{
		{
			name: "comprehensive role with all permission types",
			input: armis.RoleSettings{
				ID:       789,
				Name:     "ComprehensiveRole",
				ViprRole: true,
				Permissions: armis.Permissions{
					AdvancedPermissions: armis.AdvancedPermissions{
						All: true,
						Behavioral: armis.Behavioral{
							All:             false,
							ApplicationName: armis.Permission{All: true},
							HostName:        armis.Permission{All: false},
							ServiceName:     armis.Permission{All: true},
						},
						Device: armis.DeviceAdvanced{
							All:          true,
							DeviceNames:  armis.Permission{All: true},
							IPAddresses:  armis.Permission{All: true},
							MACAddresses: armis.Permission{All: false},
							PhoneNumbers: armis.Permission{All: true},
						},
					},
					Alert: armis.Alert{
						All: true,
						Manage: armis.Manage{
							All:              false,
							Resolve:          armis.Permission{All: true},
							WhitelistDevices: armis.Permission{All: true},
						},
						Read: armis.Permission{All: true},
					},
					Policy: armis.Policy{
						All:    true,
						Manage: armis.Permission{All: true},
						Read:   armis.Permission{All: true},
					},
					Settings: armis.Settings{
						All:      true,
						AuditLog: armis.Permission{All: true},
						Boundary: armis.Boundary{
							All: true,
							Manage: armis.ManageBoundary{
								All:    true,
								Create: armis.Permission{All: true},
								Delete: armis.Permission{All: true},
								Edit:   armis.Permission{All: true},
							},
							Read: armis.Permission{All: true},
						},
						UsersAndRoles: armis.UsersAndRoles{
							All: true,
							Manage: armis.ManageUsers{
								All: true,
								Roles: armis.UserActions{
									All:    true,
									Create: armis.Permission{All: true},
									Delete: armis.Permission{All: false},
									Edit:   armis.Permission{All: true},
								},
								Users: armis.UserActions{
									All:    true,
									Create: armis.Permission{All: true},
									Delete: armis.Permission{All: true},
									Edit:   armis.Permission{All: true},
								},
							},
							Read: armis.Permission{All: true},
						},
					},
					User: armis.User{
						All: true,
						Manage: armis.ManageUser{
							All:    true,
							Upsert: armis.Permission{All: true},
						},
						Read: armis.Permission{All: true},
					},
					Vulnerability: armis.Vulnerability{
						All: true,
						Manage: armis.ManageVuln{
							All:     true,
							Ignore:  armis.Permission{All: false},
							Resolve: armis.Permission{All: true},
							Write:   armis.Permission{All: true},
						},
						Read: armis.Permission{All: true},
					},
				},
			},
			validate: func(t *testing.T, result RoleDataSourceModel) {
				// Validate basic fields
				if result.ID.ValueString() != "789" {
					t.Errorf("Expected ID '789', got '%s'", result.ID.ValueString())
				}
				if result.Name.ValueString() != "ComprehensiveRole" {
					t.Errorf("Expected Name 'ComprehensiveRole', got '%s'", result.Name.ValueString())
				}
				if !result.ViprRole.ValueBool() {
					t.Error("Expected ViprRole to be true")
				}

				// Validate permissions structure exists
				if result.Permissions == nil {
					t.Fatal("Expected Permissions to be non-nil")
				}

				// Validate AdvancedPermissions
				if !result.Permissions.AdvancedPermissions.All.ValueBool() {
					t.Error("Expected AdvancedPermissions.All to be true")
				}
				if result.Permissions.AdvancedPermissions.Behavioral.All.ValueBool() {
					t.Error("Expected Behavioral.All to be false")
				}
				if !result.Permissions.AdvancedPermissions.Behavioral.ApplicationName.ValueBool() {
					t.Error("Expected ApplicationName to be true")
				}

				// Validate Alert permissions
				if !result.Permissions.Alert.All.ValueBool() {
					t.Error("Expected Alert.All to be true")
				}
				if result.Permissions.Alert.Manage.All.ValueBool() {
					t.Error("Expected Alert.Manage.All to be false")
				}
				if !result.Permissions.Alert.Manage.Resolve.ValueBool() {
					t.Error("Expected Alert.Manage.Resolve to be true")
				}

				// Validate Policy permissions
				if !result.Permissions.Policy.All.ValueBool() {
					t.Error("Expected Policy.All to be true")
				}
				if !result.Permissions.Policy.Manage.ValueBool() {
					t.Error("Expected Policy.Manage to be true")
				}

				// Validate Settings permissions
				if !result.Permissions.Settings.All.ValueBool() {
					t.Error("Expected Settings.All to be true")
				}
				if !result.Permissions.Settings.Boundary.All.ValueBool() {
					t.Error("Expected Settings.Boundary.All to be true")
				}

				// Validate UsersAndRoles
				if result.Permissions.Settings.UsersAndRoles.Manage.Roles.Delete.ValueBool() {
					t.Error("Expected Roles.Delete to be false")
				}
				if !result.Permissions.Settings.UsersAndRoles.Manage.Roles.Create.ValueBool() {
					t.Error("Expected Roles.Create to be true")
				}

				// Validate User permissions
				if !result.Permissions.User.All.ValueBool() {
					t.Error("Expected User.All to be true")
				}
				if !result.Permissions.User.Manage.Upsert.ValueBool() {
					t.Error("Expected User.Manage.Upsert to be true")
				}

				// Validate Vulnerability permissions
				if !result.Permissions.Vulnerability.All.ValueBool() {
					t.Error("Expected Vulnerability.All to be true")
				}
				if result.Permissions.Vulnerability.Manage.Ignore.ValueBool() {
					t.Error("Expected Vulnerability.Manage.Ignore to be false")
				}
				if !result.Permissions.Vulnerability.Manage.Write.ValueBool() {
					t.Error("Expected Vulnerability.Manage.Write to be true")
				}
			},
		},
		{
			name: "minimal role with no permissions",
			input: armis.RoleSettings{
				ID:       1,
				Name:     "ReadOnlyRole",
				ViprRole: false,
				Permissions: armis.Permissions{
					AdvancedPermissions: armis.AdvancedPermissions{
						All:        false,
						Behavioral: armis.Behavioral{All: false},
						Device:     armis.DeviceAdvanced{All: false},
					},
					Alert: armis.Alert{
						All:    false,
						Manage: armis.Manage{All: false},
						Read:   armis.Permission{All: true}, // Only read permission
					},
				},
			},
			validate: func(t *testing.T, result RoleDataSourceModel) {
				if result.ID.ValueString() != "1" {
					t.Errorf("Expected ID '1', got '%s'", result.ID.ValueString())
				}
				if result.Name.ValueString() != "ReadOnlyRole" {
					t.Errorf("Expected Name 'ReadOnlyRole', got '%s'", result.Name.ValueString())
				}
				if result.ViprRole.ValueBool() {
					t.Error("Expected ViprRole to be false")
				}
				if result.Permissions.AdvancedPermissions.All.ValueBool() {
					t.Error("Expected AdvancedPermissions.All to be false")
				}
				if !result.Permissions.Alert.Read.ValueBool() {
					t.Error("Expected Alert.Read to be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := BuildRoleDataSourceModel(&tt.input)
			tt.validate(t, result)
		})
	}
}

// TestBuildRoleRequest_RoundTrip tests that converting to API model and back preserves data.
func TestBuildRoleRequest_RoundTrip(t *testing.T) {
	t.Parallel()

	// Create a comprehensive role model using the helper
	original := createMinimalRoleModel()
	original.Name = types.StringValue("FullRole")
	original.Permissions.AdvancedPermissions.All = types.BoolValue(false)
	original.Permissions.AdvancedPermissions.Behavioral.All = types.BoolValue(true)
	original.Permissions.AdvancedPermissions.Behavioral.ApplicationName = types.BoolValue(true)
	original.Permissions.AdvancedPermissions.Behavioral.HostName = types.BoolValue(false)
	original.Permissions.AdvancedPermissions.Behavioral.ServiceName = types.BoolValue(true)
	original.Permissions.AdvancedPermissions.Device.All = types.BoolValue(true)
	original.Permissions.AdvancedPermissions.Device.DeviceNames = types.BoolValue(true)
	original.Permissions.AdvancedPermissions.Device.IPAddresses = types.BoolValue(true)
	original.Permissions.AdvancedPermissions.Device.MACAddresses = types.BoolValue(false)
	original.Permissions.AdvancedPermissions.Device.PhoneNumbers = types.BoolValue(true)
	original.Permissions.Alert.All = types.BoolValue(true)
	original.Permissions.Alert.Manage.All = types.BoolValue(false)
	original.Permissions.Alert.Manage.Resolve = types.BoolValue(true)
	original.Permissions.Alert.Manage.WhitelistDevices = types.BoolValue(true)
	original.Permissions.Alert.Read = types.BoolValue(true)

	// Convert to API model
	apiModel := BuildRoleRequest(original)

	// Set an ID (simulating API response)
	apiModel.ID = 999

	// Convert back to Terraform model
	emptyModel := RoleResourceModel{}
	result := BuildRoleResourceModel(&apiModel, emptyModel)

	// Validate key fields
	if result.Name.ValueString() != original.Name.ValueString() {
		t.Errorf("Round trip failed: Name changed from '%s' to '%s'",
			original.Name.ValueString(), result.Name.ValueString())
	}

	// Validate Advanced Permissions
	if result.Permissions.AdvancedPermissions.All.ValueBool() != original.Permissions.AdvancedPermissions.All.ValueBool() {
		t.Error("Round trip failed: AdvancedPermissions.All changed")
	}

	if result.Permissions.AdvancedPermissions.Behavioral.ApplicationName.ValueBool() !=
		original.Permissions.AdvancedPermissions.Behavioral.ApplicationName.ValueBool() {
		t.Error("Round trip failed: Behavioral.ApplicationName changed")
	}

	// Validate Alert Permissions
	if result.Permissions.Alert.All.ValueBool() != original.Permissions.Alert.All.ValueBool() {
		t.Error("Round trip failed: Alert.All changed")
	}

	if result.Permissions.Alert.Manage.Resolve.ValueBool() != original.Permissions.Alert.Manage.Resolve.ValueBool() {
		t.Error("Round trip failed: Alert.Manage.Resolve changed")
	}

	// Validate ID was set
	if result.ID.ValueString() != "999" {
		t.Errorf("Expected ID '999', got '%s'", result.ID.ValueString())
	}
}
