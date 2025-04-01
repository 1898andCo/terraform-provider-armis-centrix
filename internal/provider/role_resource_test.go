// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_RoleResource(t *testing.T) {
	resourceName := "armis_role.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test"),

					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.behavioral.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.behavioral.application_name", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.behavioral.host_name", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.behavioral.service_name", "false"),

					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.device.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.device.device_names", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.device.ip_addresses", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.device.mac_addresses", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.device.phone_numbers", "false"),

					resource.TestCheckResourceAttr(resourceName, "permissions.alert.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.alert.manage.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.alert.manage.resolve", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.alert.manage.suppress", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.alert.manage.whitelist_devices", "false"),
				),
			},
		},
	})
}

func testAccRoleResourceConfig() string {
	return `
resource "armis_role" "test" {
  name = "Test"

  permissions = {
    advanced_permissions = {
      all = false
      behavioral = {
        all              = false
        application_name = false
        host_name        = true
        service_name     = false
      }
      device = {
        all           = false
        device_names  = true
        ip_addresses  = true
        mac_addresses = true
        phone_numbers = false
      }
    }
    alert = {
      all = false
      manage = {
        all               = false
        resolve           = true
        suppress          = true
        whitelist_devices = false
      }
    }
  }
}
`
}
