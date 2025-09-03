// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

// Package utils provides utility functions and resource and data source models
// for managing resources in the Armis Centrix Terraform provider.
// This includes building request models, converting API responses to Terraform models,
// and defining the structure of Armis Centrix resources and data sources.
package utils

import (
	"fmt"
	"strconv"

	"github.com/1898andCo/terraform-provider-armis-centrix/armis"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
					All: role.Permissions.Alert.Read.ValueBool(),
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
					All: role.Permissions.Device.Read.ValueBool(),
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
						All: role.Permissions.Settings.SitesAndSensors.Read.ValueBool(),
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
	result.Permissions.Alert.Read = types.BoolValue(role.Permissions.Alert.Read.All)

	// Alert Manage Permissions
	result.Permissions.Alert.Manage.All = types.BoolValue(role.Permissions.Alert.Manage.All)
	result.Permissions.Alert.Manage.Resolve = types.BoolValue(role.Permissions.Alert.Manage.Resolve.All)
	result.Permissions.Alert.Manage.Suppress = types.BoolValue(role.Permissions.Alert.Manage.Suppress.All)
	result.Permissions.Alert.Manage.WhitelistDevices = types.BoolValue(role.Permissions.Alert.Manage.WhitelistDevices.All)

	// Device Permissions
	result.Permissions.Device.All = types.BoolValue(role.Permissions.Device.All)
	result.Permissions.Device.Read = types.BoolValue(role.Permissions.Device.Read.All)

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
	result.Permissions.Settings.SitesAndSensors.Read = types.BoolValue(role.Permissions.Settings.SitesAndSensors.Read.All)
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
