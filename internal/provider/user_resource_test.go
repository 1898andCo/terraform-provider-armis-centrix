// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_UserResource(t *testing.T) {
	resourceName := "armis_user.test"

	rName := strings.ToLower(acctest.RandomWithPrefix("tfacc-user"))
	// Generate a random ID to avoid duplicate usernames and emails
	randomID := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfig(rName, randomID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "phone", "8675309"),
					resource.TestCheckResourceAttr(resourceName, "location", "Houston"),
					resource.TestCheckResourceAttr(resourceName, "username", fmt.Sprintf("test.user-%d@test.com", randomID)),
					resource.TestCheckResourceAttr(resourceName, "email", fmt.Sprintf("test.user-%d@test.com", randomID)),
					resource.TestCheckResourceAttr(resourceName, "role_assignments.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "role_assignments.0.name.0", "Read Only"),
					resource.TestCheckResourceAttr(resourceName, "role_assignments.0.sites.0", "Lab"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
		},
	})
}

func testAccUserResourceConfig(name string, randomID int) string {
	return fmt.Sprintf(`
resource "armis_user" "test" {
  name = %q

  phone    = "8675309"
  location = "Houston"
  username = "test.user-%d@test.com"
  email    = "test.user-%d@test.com"

  role_assignments = [{
    name  = ["Read Only"]
    sites = ["Lab"]
  }]
}
`, name, randomID, randomID)
}
