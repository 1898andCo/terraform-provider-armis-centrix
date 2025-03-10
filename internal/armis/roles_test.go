// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"
	"testing"

	log "github.com/charmbracelet/log"
)

func TestCreatingRole(t *testing.T) {
	// Initialize the client
	options := Client{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	log.Info("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	role := RoleSettings{
		Name: "Test Role",
		Permissions: Permissions{
			AdvancedPermissions: AdvancedPermissions{
				All: false,
				Behavioral: Behavioral{
					All: false,
					ApplicationName: Permission{
						All: false,
					},
					HostName: Permission{
						All: false,
					},
					ServiceName: Permission{
						All: false,
					},
				},
				Device: DeviceAdvanced{
					All: false,
					DeviceNames: Permission{
						All: false,
					},
					IPAddresses: Permission{
						All: true,
					},
					MACAddresses: Permission{
						All: false,
					},
					PhoneNumbers: Permission{
						All: false,
					},
				},
			},
			Alert: Alert{
				All: true,
				Manage: Manage{
					All: true,
					Resolve: Permission{
						All: true,
					},
					Suppress: Permission{
						All: true,
					},
					WhitelistDevices: Permission{
						All: true,
					},
				},
				Read: Permission{
					All: true,
				},
			},
			Device: Device{
				All: true,
				Manage: ManageDevice{
					All: true,
					Create: Permission{
						All: true,
					},
					Delete: Permission{
						All: true,
					},
					Edit: Permission{
						All: true,
					},
					Enforce: Enforce{
						All: true,
						Create: Permission{
							All: true,
						},
						Delete: Permission{
							All: true,
						},
					},
					Merge: Permission{
						All: true,
					},
					RequestDeletedData: Permission{
						All: true,
					},
					Tags: Permission{
						All: true,
					},
				},
				Read: Permission{
					All: true,
				},
			},
			Policy: Policy{
				All: true,
				Manage: Permission{
					All: true,
				},
				Read: Permission{
					All: true,
				},
			},
			Report: Report{
				All: true,
				Export: Permission{
					All: true,
				},
				Manage: ManageReport{
					All: true,
					Create: Permission{
						All: true,
					},
					Delete: Permission{
						All: true,
					},
					Edit: Permission{
						All: true,
					},
				},
				Read: Permission{
					All: true,
				},
			},
			RiskFactor: RiskFactor{
				All: true,
				Manage: ManageRisk{
					All: true,
					Customization: Customization{
						All: true,
						Create: Permission{
							All: true,
						},
						Disable: Permission{
							All: true,
						},
						Edit: Permission{
							All: true,
						},
					},
					Status: Status{
						All: true,
						Ignore: Permission{
							All: true,
						},
						Resolve: Permission{
							All: true,
						},
					},
				},
				Read: Permission{
					All: true,
				},
			},
			Settings: Settings{
				All: true,
				AuditLog: Permission{
					All: true,
				},
				Boundary: Boundary{
					All: true,
					Manage: ManageBoundary{
						All: true,
						Create: Permission{
							All: true,
						},
						Delete: Permission{
							All: true,
						},
						Edit: Permission{
							All: true,
						},
					},
					Read: Permission{
						All: true,
					},
				},
				BusinessImpact: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				Collector: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				CustomProperties: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				Integration: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				InternalIps: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				Notifications: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				OIDC: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				SAML: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				SecretKey: Permission{
					All: true,
				},
				SecuritySettings: Permission{
					All: true,
				},
				SitesAndSensors: SitesAndSensors{
					All: true,
					Manage: ManageSensors{
						All: true,
						Sensors: Permission{
							All: true,
						},
						Sites: Permission{
							All: true,
						},
					},
					Read: Permission{
						All: true,
					},
				},
				UsersAndRoles: UsersAndRoles{
					All: true,
					Manage: ManageUsers{
						All: true,
						Roles: UserActions{
							All: true,
							Create: Permission{
								All: true,
							},
							Delete: Permission{
								All: true,
							},
							Edit: Permission{
								All: true,
							},
						},
						Users: UserActions{
							All: true,
							Create: Permission{
								All: true,
							},
							Delete: Permission{
								All: true,
							},
							Edit: Permission{
								All: true,
							},
						},
					},
					Read: Permission{
						All: true,
					},
				},
			},
			User: User{
				All: true,
				Manage: ManageUser{
					All: true,
					Upsert: Permission{
						All: true,
					},
				},
				Read: Permission{
					All: true,
				},
			},
			Vulnerability: Vulnerability{
				All: true,
				Manage: ManageVuln{
					All: true,
					Ignore: Permission{
						All: true,
					},
					Resolve: Permission{
						All: true,
					},
					Write: Permission{
						All: true,
					},
				},
				Read: Permission{
					All: true,
				},
			},
		},
	}

	// Attempt to create the role
	response, err := client.CreateRole(role)
	if err != nil {
		t.Errorf("Error creating role: %s", err)
	}

	// Log the response
	if response != false {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Info("Error marshaling server response: %s\n", err)
		} else {
			var prettyResponse bytes.Buffer
			if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
				log.Info("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
			} else {
				log.Info("Server response (raw JSON): %s\n", responseJSON)
			}
		}
	} else {
		log.Info("No response received from server.")
	}
}

func TestUpdatingRole(t *testing.T) {
	// Initialize the client
	options := Client{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	log.Info("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	role := RoleSettings{
		Name: "Test Role",
		Permissions: Permissions{
			AdvancedPermissions: AdvancedPermissions{
				All: false,
				Behavioral: Behavioral{
					All: false,
					ApplicationName: Permission{
						All: false,
					},
					HostName: Permission{
						All: true,
					},
					ServiceName: Permission{
						All: true,
					},
				},
				Device: DeviceAdvanced{
					All: true,
					DeviceNames: Permission{
						All: true,
					},
					IPAddresses: Permission{
						All: true,
					},
					MACAddresses: Permission{
						All: true,
					},
					PhoneNumbers: Permission{
						All: true,
					},
				},
			},
			Alert: Alert{
				All: false,
				Manage: Manage{
					All: false,
					Resolve: Permission{
						All: false,
					},
					Suppress: Permission{
						All: true,
					},
					WhitelistDevices: Permission{
						All: true,
					},
				},
				Read: Permission{
					All: true,
				},
			},
			Device: Device{
				All: false,
				Manage: ManageDevice{
					All: false,
					Create: Permission{
						All: true,
					},
					Delete: Permission{
						All: true,
					},
					Edit: Permission{
						All: true,
					},
					Enforce: Enforce{
						All: false,
						Create: Permission{
							All: true,
						},
						Delete: Permission{
							All: true,
						},
					},
					Merge: Permission{
						All: true,
					},
					RequestDeletedData: Permission{
						All: true,
					},
					Tags: Permission{
						All: true,
					},
				},
				Read: Permission{
					All: true,
				},
			},
			Policy: Policy{
				All: true,
				Manage: Permission{
					All: true,
				},
				Read: Permission{
					All: true,
				},
			},
			Report: Report{
				All: true,
				Export: Permission{
					All: true,
				},
				Manage: ManageReport{
					All: true,
					Create: Permission{
						All: true,
					},
					Delete: Permission{
						All: true,
					},
					Edit: Permission{
						All: true,
					},
				},
				Read: Permission{
					All: true,
				},
			},
			RiskFactor: RiskFactor{
				All: true,
				Manage: ManageRisk{
					All: true,
					Customization: Customization{
						All: true,
						Create: Permission{
							All: true,
						},
						Disable: Permission{
							All: true,
						},
						Edit: Permission{
							All: true,
						},
					},
					Status: Status{
						All: true,
						Ignore: Permission{
							All: true,
						},
						Resolve: Permission{
							All: true,
						},
					},
				},
				Read: Permission{
					All: true,
				},
			},
			Settings: Settings{
				All: true,
				AuditLog: Permission{
					All: true,
				},
				Boundary: Boundary{
					All: true,
					Manage: ManageBoundary{
						All: true,
						Create: Permission{
							All: true,
						},
						Delete: Permission{
							All: true,
						},
						Edit: Permission{
							All: true,
						},
					},
					Read: Permission{
						All: true,
					},
				},
				BusinessImpact: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				Collector: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				CustomProperties: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				Integration: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				InternalIps: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				Notifications: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				OIDC: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				SAML: ManageAndRead{
					All: true,
					Manage: Permission{
						All: true,
					},
					Read: Permission{
						All: true,
					},
				},
				SecretKey: Permission{
					All: true,
				},
				SecuritySettings: Permission{
					All: true,
				},
				SitesAndSensors: SitesAndSensors{
					All: true,
					Manage: ManageSensors{
						All: true,
						Sensors: Permission{
							All: true,
						},
						Sites: Permission{
							All: true,
						},
					},
					Read: Permission{
						All: true,
					},
				},
				UsersAndRoles: UsersAndRoles{
					All: true,
					Manage: ManageUsers{
						All: true,
						Roles: UserActions{
							All: true,
							Create: Permission{
								All: true,
							},
							Delete: Permission{
								All: true,
							},
							Edit: Permission{
								All: true,
							},
						},
						Users: UserActions{
							All: true,
							Create: Permission{
								All: true,
							},
							Delete: Permission{
								All: true,
							},
							Edit: Permission{
								All: true,
							},
						},
					},
					Read: Permission{
						All: true,
					},
				},
			},
			User: User{
				All: true,
				Manage: ManageUser{
					All: true,
					Upsert: Permission{
						All: true,
					},
				},
				Read: Permission{
					All: true,
				},
			},
			Vulnerability: Vulnerability{
				All: true,
				Manage: ManageVuln{
					All: true,
					Ignore: Permission{
						All: true,
					},
					Resolve: Permission{
						All: true,
					},
					Write: Permission{
						All: true,
					},
				},
				Read: Permission{
					All: true,
				},
			},
		},
	}

	// Attempt to get a role by name
	rolesName, err := client.GetRoleByName("Test Role")
	if err != nil {
		t.Errorf("Error getting role: %s", err)
	}

	roleId := strconv.Itoa(rolesName.ID)

	// Attempt to update the role
	response, err := client.UpdateRole(role, roleId)
	if err != nil {
		t.Errorf("Error updating role: %s", err)
	}

	// Log the response
	if response != nil {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Info("Error marshaling server response: %s\n", err)
		} else {
			var prettyResponse bytes.Buffer
			if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
				log.Info("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
			} else {
				log.Info("Server response (raw JSON): %s\n", responseJSON)
			}
		}
	} else {
		log.Info("No response received from server.")
	}
}

func TestGettingRoles(t *testing.T) {
	// Initialize the client
	options := Client{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	log.Info("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	// Attempt to get all roles
	response, err := client.GetRoles()
	if err != nil {
		t.Errorf("Error getting roles: %s", err)
	}

	// Log the response
	if response != nil {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Info("Error marshaling server response: %s\n", err)
		}

		// Attempt to pretty-print the JSON
		var prettyResponse bytes.Buffer
		if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
			log.Info("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
		} else {
			log.Info("Failed to pretty-print JSON.")
		}
	} else {
		log.Info("No response received from server.")
	}
}

func TestGettingRoleByName(t *testing.T) {
	// Initialize the client
	options := Client{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	log.Info("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	// Attempt to get a role
	response, err := client.GetRoleByName("Test Role")
	if err != nil {
		t.Errorf("Error getting role: %s", err)
	}

	// Log the response
	if response != nil {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Info("Error marshaling server response: %s\n", err)
		}

		// Attempt to pretty-print the JSON
		var prettyResponse bytes.Buffer
		if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
			log.Info("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
		} else {
			log.Info("Failed to pretty-print JSON.")
		}
	} else {
		log.Info("No response received from server.")
	}
}

func TestGettingRoleByID(t *testing.T) {
	// Initialize the client
	options := Client{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	log.Info("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	// Attempt to get a role
	response, err := client.GetRoleByID("10")
	if err != nil {
		t.Errorf("Error getting role: %s", err)
	}

	// Log the response
	if response != nil {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Info("Error marshaling server response: %s\n", err)
		}

		// Attempt to pretty-print the JSON
		var prettyResponse bytes.Buffer
		if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
			log.Info("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
		} else {
			log.Info("Failed to pretty-print JSON.")
		}
	} else {
		log.Info("No response received from server.")
	}
}

func TestDeletingRole(t *testing.T) {
	// Initialize the client
	options := Client{
		ApiUrl: os.Getenv("ARMIS_API_URL"),
		ApiKey: os.Getenv("ARMIS_API_KEY"),
	}
	log.Info("Initializing client with API URL: %s\n", options.ApiUrl)

	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	// Attempt to get a role by name
	role, err := client.GetRoleByName("Test Role")
	if err != nil {
		t.Errorf("Error getting role: %s", err)
	}

	roleId := strconv.Itoa(role.ID)

	// Attempt to delete a role
	response, err := client.DeleteRole(roleId)
	if err != nil {
		t.Errorf("Error deleting role: %s", err)
	}

	// Log the response
	if response != false {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Info("Error marshaling server response: %s\n", err)
		}

		// Attempt to pretty-print the JSON
		var prettyResponse bytes.Buffer
		if err := json.Indent(&prettyResponse, responseJSON, "", "  "); err == nil {
			log.Info("\n=== Parsed Response Body ===\n%s\n", prettyResponse.String())
		} else {
			log.Info("Failed to pretty-print JSON.")
		}
	} else {
		log.Info("No response received from server.")
	}
}
