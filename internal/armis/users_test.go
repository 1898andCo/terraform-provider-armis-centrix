// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"testing"
)

func TestGettingUsers(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	res, err := client.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("get users: %v", err)
	}
	prettyPrint(res)
}

func TestCreatingUser(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	payload := UserSettings{
		Name:     "Test User",
		Username: "testuser",
		Email:    "test.user@1898andco.io",
		RoleAssignment: []RoleAssignment{{
			Name:  []string{"Read Only"},
			Sites: []string{"Lab"},
		}},
	}

	res, err := client.CreateUser(context.Background(), payload)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	prettyPrint(res)

	if res.Name != payload.Name {
		t.Errorf("expected name %q, got %q", payload.Name, res.Name)
	}
	if res.Email != payload.Email {
		t.Errorf("expected email %q, got %q", payload.Email, res.Email)
	}
}

func TestGettingUser(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	email := "test.user@1898andco.io"
	res, err := client.GetUser(context.Background(), email)
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	prettyPrint(res)
}

func TestUpdatingUser(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	updated := UserSettings{
		Name:     "Test User",
		Username: "testupdateduser",
		Email:    "test.user@1898andco.io",
		RoleAssignment: []RoleAssignment{{
			Name:  []string{"Admin"},
			Sites: []string{"Lab"},
		}},
	}

	res, err := client.UpdateUser(context.Background(), updated, updated.Email)
	if err != nil {
		t.Fatalf("update user: %v", err)
	}
	prettyPrint(res)

	if res.Username != updated.Username {
		t.Errorf("expected username %q, got %q", updated.Username, res.Username)
	}
	if res.RoleAssignment[0].Name[0] != "Admin" {
		t.Errorf("expected role 'Admin', got %q", res.RoleAssignment[0].Name[0])
	}
}

func TestDeletingUser(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	email := "test.user@1898andco.io"
	ok, err := client.DeleteUser(context.Background(), email)
	if err != nil {
		t.Fatalf("delete user: %v", err)
	}
	if !ok {
		t.Fatalf("user %s not deleted", email)
	}
}
