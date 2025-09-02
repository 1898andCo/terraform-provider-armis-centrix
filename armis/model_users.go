// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

// Struct to match the entire API response.
type UserAPIResponse struct {
	Data    Users `json:"data"`
	Success bool  `json:"success,omitempty"`
}

// Struct for the "data" key.
type Users struct {
	Users []UserSettings `json:"users,omitempty"`
}

// Struct for individual user settings.
type UserSettings struct {
	Email                   string           `json:"email,omitempty"`
	ID                      int              `json:"id,omitempty"`
	IsActive                bool             `json:"isActive,omitempty"`
	LastLoginTime           string           `json:"lastLoginTime,omitempty"`
	Location                string           `json:"location,omitempty"`
	Name                    string           `json:"name,omitempty"`
	Phone                   string           `json:"phone,omitempty"`
	PovEULASigningDate      string           `json:"povEulaSigningDate,omitempty"`
	ProdEULASigningDate     string           `json:"prodEulaSigningDate,omitempty"`
	ReportPermissions       string           `json:"reportPermissions,omitempty"`
	Role                    string           `json:"role,omitempty"`
	RoleAssignment          []RoleAssignment `json:"roleAssignment,omitempty"`
	Title                   string           `json:"title,omitempty"`
	TwoFactorAuthentication bool             `json:"twoFactorAuthentication,omitempty"`
	Username                string           `json:"username,omitempty"`
}

type RoleAssignment struct {
	Name  []string `json:"name,omitempty"`
	Sites []string `json:"sites,omitempty"`
}

// Structs for getting a singular user.
type GetUserAPIResponse struct {
	Data    UserSettings `json:"data"`
	Success bool         `json:"success"`
}

// Structs for creating users.
type CreateUserAPIResponse struct {
	Data    UserSettings `json:"data"` // Directly map to a single user object
	Success bool         `json:"success"`
}

// Structs for updating users.
type UpdateUserAPIResponse struct {
	Data    UserSettings `json:"data"` // Directly map to a single user object
	Success bool         `json:"success"`
}

// Structs for deleting users.
type DeleteUserAPIResponse struct {
	Success bool `json:"success"`
}
