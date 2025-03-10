// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_UserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.armis_user.test", "users.#"),
					resource.TestCheckResourceAttrSet("data.armis_user.test", "users.0.id"),
					resource.TestCheckResourceAttrSet("data.armis_user.test", "users.0.name"),
					resource.TestCheckResourceAttrSet("data.armis_user.test", "users.0.email"),
					resource.TestCheckResourceAttrSet("data.armis_user.test", "users.0.is_active"),
					resource.TestCheckResourceAttrSet("data.armis_user.test", "users.0.role_assignment.#"),
					resource.TestCheckResourceAttrSet("data.armis_user.test", "users.0.role_assignment.0.name.#"),
					resource.TestCheckResourceAttrSet("data.armis_user.test", "users.0.two_factor_authentication"),
					resource.TestCheckResourceAttrSet("data.armis_user.test", "users.0.username"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig() string {
	return `
data "armis_user" "test" {}
`
}
