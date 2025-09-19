// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_BoundaryDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBoundaryDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.armis_boundary.test", "boundaries.#"),
					resource.TestCheckResourceAttrSet("data.armis_boundary.test", "boundaries.0.id"),
					resource.TestCheckResourceAttrSet("data.armis_boundary.test", "boundaries.0.name"),
					resource.TestCheckResourceAttrSet("data.armis_boundary.test", "boundaries.0.affected_sites"),
					resource.TestCheckResourceAttrSet("data.armis_boundary.test", "boundaries.0.rule_aql.and.#"),
					resource.TestCheckResourceAttrSet("data.armis_boundary.test", "boundaries.0.rule_aql.or.#"),
				),
			},
		},
	})
}

func testAccBoundaryDataSourceConfig() string {
	return `
data "armis_boundary" "test" {}
`
}
