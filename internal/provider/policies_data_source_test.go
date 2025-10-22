// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_PoliciesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPoliciesDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.armis_policies.test", "policies.#"),
					resource.TestCheckResourceAttrSet("data.armis_policies.test", "policies.0.id"),
					resource.TestCheckResourceAttrSet("data.armis_policies.test", "policies.0.name"),
				),
			},
		},
	})
}

func testAccPoliciesDataSourceConfig() string {
	return `
	data "armis_policies" "test" {}
	`
}
