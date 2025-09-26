// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

type BaseResponse struct {
	Success bool `json:"success"`
}

type CreateRoleAPIResponse struct {
	BaseResponse
}

type DeleteRoleAPIResponse struct {
	BaseResponse
}

type RolesAPIResponse struct {
	Data []RoleSettings `json:"data,omitempty"`
	BaseResponse
}

type RoleSettings struct {
	Name        string      `json:"name"`
	Permissions Permissions `json:"permissions"`
	ID          int         `json:"roleId,omitempty"`
	ViprRole    bool        `json:"viprRole,omitempty"`
}

type Permissions struct {
	AdvancedPermissions AdvancedPermissions `json:"advancedPermissions,omitempty"`
	Alert               Alert               `json:"alert,omitempty"`
	Device              Device              `json:"device,omitempty"`
	Policy              Policy              `json:"policy,omitempty"`
	Report              Report              `json:"report,omitempty"`
	RiskFactor          RiskFactor          `json:"risk_factor,omitempty"`
	Settings            Settings            `json:"settings,omitempty"`
	User                User                `json:"user,omitempty"`
	Vulnerability       Vulnerability       `json:"vulnerability,omitempty"`
}

type AdvancedPermissions struct {
	All        bool           `json:"all,omitempty"`
	Behavioral Behavioral     `json:"behavioral,omitempty"`
	Device     DeviceAdvanced `json:"device,omitempty"`
}

type Behavioral struct {
	All             bool       `json:"all,omitempty"`
	ApplicationName Permission `json:"applicationName,omitempty"`
	HostName        Permission `json:"hostName,omitempty"`
	ServiceName     Permission `json:"serviceName,omitempty"`
}

type DeviceAdvanced struct {
	All          bool       `json:"all,omitempty"`
	DeviceNames  Permission `json:"deviceNames,omitempty"`
	IPAddresses  Permission `json:"ipAddresses,omitempty"`
	MACAddresses Permission `json:"macAddresses,omitempty"`
	PhoneNumbers Permission `json:"phoneNumbers,omitempty"`
}

type Permission struct {
	All bool `json:"all,omitempty"`
}

type Alert struct {
	All    bool       `json:"all,omitempty"`
	Manage Manage     `json:"manage,omitempty"`
	Read   Permission `json:"read,omitempty"`
}

type Manage struct {
	All              bool       `json:"all,omitempty"`
	Resolve          Permission `json:"resolve,omitempty"`
	WhitelistDevices Permission `json:"whitelistDevices,omitempty"`
}

type Device struct {
	All    bool         `json:"all,omitempty"`
	Manage ManageDevice `json:"manage,omitempty"`
	Read   Permission   `json:"read,omitempty"`
}

type ManageDevice struct {
	All                bool       `json:"all,omitempty"`
	Create             Permission `json:"create,omitempty"`
	Delete             Permission `json:"delete,omitempty"`
	Edit               Permission `json:"edit,omitempty"`
	Enforce            Enforce    `json:"enforce,omitempty"`
	Merge              Permission `json:"merge,omitempty"`
	RequestDeletedData Permission `json:"request_deleted_data,omitempty"`
	Tags               Permission `json:"tags,omitempty"`
}

type Enforce struct {
	All    bool       `json:"all,omitempty"`
	Create Permission `json:"create,omitempty"`
	Delete Permission `json:"delete,omitempty"`
}

type Policy struct {
	All    bool       `json:"all,omitempty"`
	Manage Permission `json:"manage,omitempty"`
	Read   Permission `json:"read,omitempty"`
}

type Report struct {
	All    bool         `json:"all,omitempty"`
	Export Permission   `json:"export,omitempty"`
	Manage ManageReport `json:"manage,omitempty"`
	Read   Permission   `json:"read,omitempty"`
}

type ManageReport struct {
	All    bool       `json:"all,omitempty"`
	Create Permission `json:"create,omitempty"`
	Delete Permission `json:"delete,omitempty"`
	Edit   Permission `json:"edit,omitempty"`
}

type RiskFactor struct {
	All    bool       `json:"all,omitempty"`
	Manage ManageRisk `json:"manage,omitempty"`
	Read   Permission `json:"read,omitempty"`
}

type ManageRisk struct {
	All           bool          `json:"all,omitempty"`
	Customization Customization `json:"customization,omitempty"`
	Status        Status        `json:"status,omitempty"`
}

type Customization struct {
	All     bool       `json:"all,omitempty"`
	Create  Permission `json:"create,omitempty"`
	Disable Permission `json:"disable,omitempty"`
	Edit    Permission `json:"edit,omitempty"`
}

type Status struct {
	All     bool       `json:"all,omitempty"`
	Ignore  Permission `json:"ignore,omitempty"`
	Resolve Permission `json:"resolve,omitempty"`
}

type Settings struct {
	All              bool            `json:"all,omitempty"`
	AuditLog         Permission      `json:"auditLog,omitempty"`
	Boundary         Boundary        `json:"boundary,omitempty"`
	BusinessImpact   ManageAndRead   `json:"businessImpact,omitempty"`
	Collector        ManageAndRead   `json:"collector,omitempty"`
	CustomProperties ManageAndRead   `json:"customProperties,omitempty"`
	Integration      ManageAndRead   `json:"integration,omitempty"`
	InternalIps      ManageAndRead   `json:"internalIps,omitempty"`
	Notifications    ManageAndRead   `json:"notifications,omitempty"`
	OIDC             ManageAndRead   `json:"oidc,omitempty"`
	SAML             ManageAndRead   `json:"saml,omitempty"`
	SecretKey        Permission      `json:"secretKey,omitempty"`
	SecuritySettings Permission      `json:"securitySettings,omitempty"`
	SitesAndSensors  SitesAndSensors `json:"sitesAndSensors,omitempty"`
	UsersAndRoles    UsersAndRoles   `json:"usersAndRoles,omitempty"`
}

type Boundary struct {
	All    bool           `json:"all,omitempty"`
	Manage ManageBoundary `json:"manage,omitempty"`
	Read   Permission     `json:"read,omitempty"`
}

type ManageBoundary struct {
	All    bool       `json:"all,omitempty"`
	Create Permission `json:"create,omitempty"`
	Delete Permission `json:"delete,omitempty"`
	Edit   Permission `json:"edit,omitempty"`
}

type ManageAndRead struct {
	All    bool       `json:"all,omitempty"`
	Manage Permission `json:"manage,omitempty"`
	Read   Permission `json:"read,omitempty"`
}

type SitesAndSensors struct {
	All    bool          `json:"all,omitempty"`
	Manage ManageSensors `json:"manage,omitempty"`
	Read   Permission    `json:"read,omitempty"`
}

type ManageSensors struct {
	All     bool       `json:"all,omitempty"`
	Sensors Permission `json:"sensors,omitempty"`
	Sites   Permission `json:"sites,omitempty"`
}

type UsersAndRoles struct {
	All    bool        `json:"all,omitempty"`
	Manage ManageUsers `json:"manage,omitempty"`
	Read   Permission  `json:"read,omitempty"`
}

type ManageUsers struct {
	All   bool        `json:"all,omitempty"`
	Roles UserActions `json:"roles,omitempty"`
	Users UserActions `json:"users,omitempty"`
}

type UserActions struct {
	All    bool       `json:"all,omitempty"`
	Create Permission `json:"create,omitempty"`
	Delete Permission `json:"delete,omitempty"`
	Edit   Permission `json:"edit,omitempty"`
}

type User struct {
	All    bool       `json:"all,omitempty"`
	Manage ManageUser `json:"manage,omitempty"`
	Read   Permission `json:"read,omitempty"`
}

type ManageUser struct {
	All    bool       `json:"all,omitempty"`
	Upsert Permission `json:"upsert,omitempty"`
}

type Vulnerability struct {
	All    bool       `json:"all,omitempty"`
	Manage ManageVuln `json:"manage,omitempty"`
	Read   Permission `json:"read,omitempty"`
}

type ManageVuln struct {
	All     bool       `json:"all,omitempty"`
	Ignore  Permission `json:"ignore,omitempty"`
	Resolve Permission `json:"resolve,omitempty"`
	Write   Permission `json:"write,omitempty"`
}
