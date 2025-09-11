// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_UserResource(t *testing.T) {
	// Generate unique values to avoid conflicts in CI (email/username must be unique).
	suffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-acc-user-%s", suffix)
	email := fmt.Sprintf("tfacc%s@example.com", suffix) // starts with alpha, satisfies regex validator
	username := email

	resourceName := "armis_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfig(name, username, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "phone", "8675309"),
					resource.TestCheckResourceAttr(resourceName, "location", "Houston"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "email", email),
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

func testAccUserResourceConfig(name, username, email string) string {
	return fmt.Sprintf(`
resource "armis_user" "test" {
  name = %q

  phone    = "8675309"
  location = "Houston"
  username = %q
  email    = %q

  role_assignments = {
    name  = ["Read Only"]
    sites = ["Lab"]
  }
}
`, name, username, email)
}
