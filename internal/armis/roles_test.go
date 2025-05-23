// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"strconv"
	"testing"
)

func TestCreatingRole(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	role := RoleSettings{
		Name: "Test Role",
		Permissions: Permissions{
			AdvancedPermissions: AdvancedPermissions{All: true},
			Alert:               Alert{All: true},
			Device:              Device{All: true},
			Policy:              Policy{All: true},
			Report:              Report{All: true},
			RiskFactor:          RiskFactor{All: true},
			Settings:            Settings{All: true},
			User:                User{All: true},
			Vulnerability:       Vulnerability{All: true},
		},
	}

	res, err := client.CreateRole(context.Background(), role)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	prettyPrint(res)
}

func TestUpdatingRole(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	// Lookup the role ID by name first.
	roleMeta, err := client.GetRoleByName(context.Background(), "Test Role")
	if err != nil {
		t.Fatalf("lookup role: %v", err)
	}
	id := strconv.Itoa(roleMeta.ID)

	updated := RoleSettings{
		Name: "Test Role Updated",
		Permissions: Permissions{
			AdvancedPermissions: AdvancedPermissions{All: false},
			Alert:               Alert{All: true},
		},
	}

	res, err := client.UpdateRole(context.Background(), updated, id)
	if err != nil {
		t.Fatalf("update role: %v", err)
	}
	prettyPrint(res)
}

func TestGettingRoles(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	res, err := client.GetRoles(context.Background())
	if err != nil {
		t.Fatalf("get roles: %v", err)
	}
	prettyPrint(res)
}

func TestGettingRoleByName(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	res, err := client.GetRoleByName(context.Background(), "Test Role")
	if err != nil {
		t.Fatalf("get role by name: %v", err)
	}
	prettyPrint(res)
}

func TestGettingRoleByID(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	const id = "10" // adjust as needed
	res, err := client.GetRoleByID(context.Background(), id)
	if err != nil {
		t.Fatalf("get role by id: %v", err)
	}
	prettyPrint(res)
}

func TestDeletingRole(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	roleMeta, err := client.GetRoleByName(context.Background(), "Test Role")
	if err != nil {
		t.Fatalf("lookup role: %v", err)
	}
	id := strconv.Itoa(roleMeta.ID)

	ok, err := client.DeleteRole(context.Background(), id)
	if err != nil {
		t.Fatalf("delete role: %v", err)
	}
	if !ok {
		t.Fatalf("role %s not deleted", id)
	}
}
