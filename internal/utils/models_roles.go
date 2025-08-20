// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RoleResourceModel maps the RoleSettings schema data.
type RoleResourceModel struct {
	Name        types.String      `tfsdk:"name"`
	Permissions *PermissionsModel `tfsdk:"permissions"`
	ID          types.String      `tfsdk:"id"`
}

// RoleDataSourceModel defines the structure for the role data source model.
type RoleDataSourceModel struct {
	Name        types.String      `tfsdk:"name"`
	Permissions *PermissionsModel `tfsdk:"permissions"`
	ID          types.String      `tfsdk:"role_id"`
	ViprRole    types.Bool        `tfsdk:"vipr_role"`
}

// AllPermissionsModel defines the structure for all permissions.
type AllPermissionsModel struct {
	All types.Bool `tfsdk:"all"`
}

// PermissionsModel defines the structure for permissions.
type PermissionsModel struct {
	AdvancedPermissions *AdvancedPermissionsModel `tfsdk:"advanced_permissions"`
	Alert               *AlertModel               `tfsdk:"alert"`
	Device              *DeviceModel              `tfsdk:"device"`
	Policy              *PolicyModel              `tfsdk:"policy"`
	Report              *ReportModel              `tfsdk:"report"`
	RiskFactor          *RiskFactorModel          `tfsdk:"risk_factor"`
	Settings            *SettingsModel            `tfsdk:"settings"`
	User                *UserModel                `tfsdk:"user"`
	Vulnerability       *VulnerabilityModel       `tfsdk:"vulnerability"`
}

// AdvancedPermissionsModel defines the structure for advanced permissions.
type AdvancedPermissionsModel struct {
	All        types.Bool           `tfsdk:"all"`
	Behavioral *BehavioralModel     `tfsdk:"behavioral"`
	Device     *DeviceAdvancedModel `tfsdk:"device"`
}

// BehavioralModel defines the structure for behavioral permissions.
type BehavioralModel struct {
	All             types.Bool `tfsdk:"all"`
	ApplicationName types.Bool `tfsdk:"application_name"`
	HostName        types.Bool `tfsdk:"host_name"`
	ServiceName     types.Bool `tfsdk:"service_name"`
}

// DeviceAdvancedModel defines the structure for device advanced permissions.
type DeviceAdvancedModel struct {
	All          types.Bool `tfsdk:"all"`
	DeviceNames  types.Bool `tfsdk:"device_names"`
	IPAddresses  types.Bool `tfsdk:"ip_addresses"`
	MACAddresses types.Bool `tfsdk:"mac_addresses"`
	PhoneNumbers types.Bool `tfsdk:"phone_numbers"`
}

// AlertModel defines the structure for alert permissions.
type AlertModel struct {
	All    types.Bool           `tfsdk:"all"`
	Manage *ManageAlertsModel   `tfsdk:"manage"`
	Read   *AllPermissionsModel `tfsdk:"read"`
}

// ManageAlertsModel defines the structure for manage permissions.
type ManageAlertsModel struct {
	All              types.Bool `tfsdk:"all"`
	Resolve          types.Bool `tfsdk:"resolve"`
	Suppress         types.Bool `tfsdk:"suppress"`
	WhitelistDevices types.Bool `tfsdk:"whitelist_devices"`
}

// DeviceModel defines the structure for device permissions.
type DeviceModel struct {
	All    types.Bool           `tfsdk:"all"`
	Manage *ManageDeviceModel   `tfsdk:"manage"`
	Read   *AllPermissionsModel `tfsdk:"read"`
}

// ManageDeviceModel defines the structure for device management permissions.
type ManageDeviceModel struct {
	All                types.Bool    `tfsdk:"all"`
	Create             types.Bool    `tfsdk:"create"`
	Delete             types.Bool    `tfsdk:"delete"`
	Edit               types.Bool    `tfsdk:"edit"`
	Enforce            *EnforceModel `tfsdk:"enforce"`
	Merge              types.Bool    `tfsdk:"merge"`
	RequestDeletedData types.Bool    `tfsdk:"request_deleted_data"`
	Tags               types.Bool    `tfsdk:"tags"`
}

// EnforceModel defines the structure for enforce permissions.
type EnforceModel struct {
	All    types.Bool `tfsdk:"all"`
	Create types.Bool `tfsdk:"create"`
	Delete types.Bool `tfsdk:"delete"`
}

// PolicyModel defines the structure for policy permissions.
type PolicyModel struct {
	All    types.Bool `tfsdk:"all"`
	Manage types.Bool `tfsdk:"manage"`
	Read   types.Bool `tfsdk:"read"`
}

// ReportModel defines the structure for report permissions.
type ReportModel struct {
	All    types.Bool         `tfsdk:"all"`
	Export types.Bool         `tfsdk:"export"`
	Manage *ManageReportModel `tfsdk:"manage"`
	Read   types.Bool         `tfsdk:"read"`
}

// ManageReportModel defines the structure for report management permissions.
type ManageReportModel struct {
	All    types.Bool `tfsdk:"all"`
	Create types.Bool `tfsdk:"create"`
	Delete types.Bool `tfsdk:"delete"`
	Edit   types.Bool `tfsdk:"edit"`
}

// RiskFactorModel defines the structure for risk factor permissions.
type RiskFactorModel struct {
	All    types.Bool       `tfsdk:"all"`
	Manage *ManageRiskModel `tfsdk:"manage"`
	Read   types.Bool       `tfsdk:"read"`
}

// ManageRiskModel defines the structure for risk management permissions.
type ManageRiskModel struct {
	All           types.Bool          `tfsdk:"all"`
	Customization *CustomizationModel `tfsdk:"customization"`
	Status        *StatusModel        `tfsdk:"status"`
}

// CustomizationModel defines the structure for customization permissions.
type CustomizationModel struct {
	All     types.Bool `tfsdk:"all"`
	Create  types.Bool `tfsdk:"create"`
	Disable types.Bool `tfsdk:"disable"`
	Edit    types.Bool `tfsdk:"edit"`
}

// StatusModel defines the structure for status permissions.
type StatusModel struct {
	All     types.Bool `tfsdk:"all"`
	Ignore  types.Bool `tfsdk:"ignore"`
	Resolve types.Bool `tfsdk:"resolve"`
}

// SettingsModel defines the structure for settings permissions.
type SettingsModel struct {
	All              types.Bool            `tfsdk:"all"`
	AuditLog         types.Bool            `tfsdk:"audit_log"`
	Boundary         *BoundaryModel        `tfsdk:"boundary"`
	BusinessImpact   *ManageAndReadModel   `tfsdk:"business_impact"`
	Collector        *ManageAndReadModel   `tfsdk:"collector"`
	CustomProperties *ManageAndReadModel   `tfsdk:"custom_properties"`
	Integration      *ManageAndReadModel   `tfsdk:"integration"`
	InternalIps      *ManageAndReadModel   `tfsdk:"internal_ips"`
	Notifications    *ManageAndReadModel   `tfsdk:"notifications"`
	OIDC             *ManageAndReadModel   `tfsdk:"oidc"`
	SAML             *ManageAndReadModel   `tfsdk:"saml"`
	SecretKey        types.Bool            `tfsdk:"secret_key"`
	SecuritySettings types.Bool            `tfsdk:"security_settings"`
	SitesAndSensors  *SitesAndSensorsModel `tfsdk:"sites_and_sensors"`
	UsersAndRoles    *UsersAndRolesModel   `tfsdk:"users_and_roles"`
}

// BoundaryModel defines the structure for boundary permissions.
type BoundaryModel struct {
	All    types.Bool           `tfsdk:"all"`
	Manage *ManageBoundaryModel `tfsdk:"manage"`
	Read   types.Bool           `tfsdk:"read"`
}

// ManageBoundaryModel defines the structure for managing boundary permissions.
type ManageBoundaryModel struct {
	All    types.Bool `tfsdk:"all"`
	Create types.Bool `tfsdk:"create"`
	Delete types.Bool `tfsdk:"delete"`
	Edit   types.Bool `tfsdk:"edit"`
}

// ManageAndReadModel defines the structure for manage and read permissions.
type ManageAndReadModel struct {
	All    types.Bool `tfsdk:"all"`
	Manage types.Bool `tfsdk:"manage"`
	Read   types.Bool `tfsdk:"read"`
}

// SitesAndSensorsModel defines the structure for sites and sensors permissions.
type SitesAndSensorsModel struct {
	All    types.Bool                  `tfsdk:"all"`
	Manage *ManageSitesAndSensorsModel `tfsdk:"manage"`
	Read   *AllPermissionsModel        `tfsdk:"read"`
}

// ManageSitesAndSensorsModel defines the structure for managing sites and sensors permissions.
type ManageSitesAndSensorsModel struct {
	All     types.Bool `tfsdk:"all"`
	Sensors types.Bool `tfsdk:"sensors"`
	Sites   types.Bool `tfsdk:"sites"`
}

// UsersAndRolesModel defines the structure for users and roles permissions.
type UsersAndRolesModel struct {
	All    types.Bool                `tfsdk:"all"`
	Manage *ManageUsersAndRolesModel `tfsdk:"manage"`
	Read   types.Bool                `tfsdk:"read"`
}

// ManageUsersAndRolesModel defines the structure for managing users and roles permissions.
type ManageUsersAndRolesModel struct {
	All   types.Bool        `tfsdk:"all"`
	Roles *ManageRolesModel `tfsdk:"roles"`
	Users *ManageUsersModel `tfsdk:"users"`
}

// ManageRolesModel defines the structure for managing roles permissions.
type ManageRolesModel struct {
	All    types.Bool `tfsdk:"all"`
	Create types.Bool `tfsdk:"create"`
	Delete types.Bool `tfsdk:"delete"`
	Edit   types.Bool `tfsdk:"edit"`
}

// ManageUsersModel defines the structure for managing users permissions.
type ManageUsersModel struct {
	All    types.Bool `tfsdk:"all"`
	Create types.Bool `tfsdk:"create"`
	Delete types.Bool `tfsdk:"delete"`
	Edit   types.Bool `tfsdk:"edit"`
}

// UserModel defines the structure for user permissions.
type UserModel struct {
	All    types.Bool       `tfsdk:"all"`
	Manage *ManageUserModel `tfsdk:"manage"`
	Read   types.Bool       `tfsdk:"read"`
}

// ManageUserModel defines the structure for managing user permissions.
type ManageUserModel struct {
	All    types.Bool `tfsdk:"all"`
	Upsert types.Bool `tfsdk:"upsert"`
}

// VulnerabilityModel defines the structure for vulnerability permissions.
type VulnerabilityModel struct {
	All    types.Bool                `tfsdk:"all"`
	Manage *ManageVulnerabilityModel `tfsdk:"manage"`
	Read   types.Bool                `tfsdk:"read"`
}

// ManageVulnerabilityModel defines the structure for managing vulnerabilities permissions.
type ManageVulnerabilityModel struct {
	All     types.Bool `tfsdk:"all"`
	Ignore  types.Bool `tfsdk:"ignore"`
	Resolve types.Bool `tfsdk:"resolve"`
	Write   types.Bool `tfsdk:"write"`
}
