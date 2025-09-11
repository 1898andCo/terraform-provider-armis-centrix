// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_UserResource(t *testing.T) {
	resourceName := "armis_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "phone", "8675309"),
					resource.TestCheckResourceAttr(resourceName, "location", "Houston"),
					resource.TestCheckResourceAttr(resourceName, "username", "test.user@test.com"),
					resource.TestCheckResourceAttr(resourceName, "email", "test.user@test.com"),
					resource.TestCheckResourceAttr(resourceName, "role_assignments.name.0", "Read Only"),
					resource.TestCheckResourceAttr(resourceName, "role_assignments.sites.#", "1"),
				),
			},
			{
               ResourceName: resourceName,
               ImportState:  true,
            },
		},
	})
}

func testAccUserResourceConfig() string {
	return `
resource "armis_user" "test" {
  name = "test"

  phone    = "8675309"
  location = "Houston"
  username = "test.user@test.com"
  email    = "test.user@test.com"

  role_assignments = {
    name  = ["Read Only"]
    sites = ["Lab"]
  }
}
`
}