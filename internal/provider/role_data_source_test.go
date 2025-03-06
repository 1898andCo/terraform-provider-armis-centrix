// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_RoleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.armis_role.test", "name"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "role_id"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "vipr_role"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.advanced_permissions.all"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.advanced_permissions.behavioral.all"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.advanced_permissions.behavioral.application_name"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.advanced_permissions.behavioral.host_name"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.advanced_permissions.behavioral.service_name"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.advanced_permissions.device.all"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.advanced_permissions.device.device_names"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.advanced_permissions.device.ip_addresses"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.advanced_permissions.device.mac_addresses"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.advanced_permissions.device.phone_numbers"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.alert.all"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.alert.manage.all"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.alert.manage.resolve"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.alert.manage.suppress"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.alert.manage.whitelist_devices"),
					resource.TestCheckResourceAttrSet("data.armis_role.test", "permissions.alert.read"),
				),
			},
		},
	})
}

func testAccRoleDataSourceConfig() string {
	return fmt.Sprintf(`
		data "armis_role" "test" {
			name = "Stakeholder"
		}
	`)
}
