// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

// Struct to match the entire API response.
type CollectorApiResponse struct {
	Data    Collectors `json:"data"`
	Success bool       `json:"success,omitempty"`
	Count   *int       `json:"count,omitempty"`
	Next    *string    `json:"next,omitempty"`
	Prev    *string    `json:"prev,omitempty"`
	Total   *int       `json:"total,omitempty"`
}

// Struct for the "data" key.
type Collectors struct {
	Collectors []CollectorSettings `json:"collectors,omitempty"`
}

// Struct for individual collector settings.
type CollectorSettings struct {
	BootTime           string `json:"bootTime,omitempty"`
	City               string `json:"city,omitempty"`
	ClusterID          int    `json:"clusterId,omitempty"`
	CollectorNumber    int    `json:"collectorNumber,omitempty"`
	Country            string `json:"country,omitempty"`
	DefaultGateway     string `json:"defaultGateway,omitempty"`
	HTTPSProxyRedacted string `json:"httpsProxyRedacted,omitempty"`
	IPAddress          string `json:"ipAddress,omitempty"`
	LastSeen           string `json:"lastSeen,omitempty"`
	MacAddress         string `json:"macAddress,omitempty"`
	Name               string `json:"name,omitempty"`
	Namespace          string `json:"namespace,omitempty"`
	ProductSerial      string `json:"productSerial,omitempty"`
	Status             string `json:"status,omitempty"`
	Subnet             string `json:"subnet,omitempty"`
	SystemVendor       string `json:"systemVendor,omitempty"`
	Type               string `json:"type,omitempty"`
}

type CreateCollectorSettings struct {
	DeploymentType string `json:"deploymentType"`
	Name           string `json:"name"`
}

type UpdateCollectorSettings struct {
	DeploymentType string `json:"deploymentType"`
	Name           string `json:"name"`
}

// Struct for creating collectors.
type CreateCollectorApiResponse struct {
	Data    NewCollectorSettings `json:"data"`
	Success bool                 `json:"success,omitempty"`
}

// Struct for getting a single collector.
type GetCollectorApiResponse struct {
	Data    CollectorSettings `json:"data"`
	Success bool              `json:"success,omitempty"`
}

// Struct for updating collector.
type UpdateCollectorApiResponse struct {
	Data    CollectorSettings `json:"data"`
	Success bool              `json:"success"`
}

// Structs for deleting collector.
type DeleteCollectorApiResponse struct {
	Success bool `json:"success"`
}

// Struct for creating a new collector response.
type NewCollectorSettings struct {
	CollectorID int    `json:"collectorId,omitempty"`
	LicenseKey  string `json:"licenseKey,omitempty"`
	Password    string `json:"password,omitempty"`
	User        string `json:"user,omitempty"`
}
