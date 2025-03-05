// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_SitesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSitesDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.armis_sites.test", "sites.#"),
					resource.TestCheckResourceAttrSet("data.armis_sites.test", "sites.0.id"),
					resource.TestCheckResourceAttrSet("data.armis_sites.test", "sites.0.name"),
					resource.TestCheckResourceAttrSet("data.armis_sites.test", "sites.0.latitude"),
					resource.TestCheckResourceAttrSet("data.armis_sites.test", "sites.0.longitude"),
					resource.TestCheckResourceAttrSet("data.armis_sites.test", "sites.0.location"),
					resource.TestCheckResourceAttrSet("data.armis_sites.test", "sites.0.user"),
				),
			},
		},
	})
}

func testAccSitesDataSourceConfig() string {
	return fmt.Sprintf(`
data "armis_sites" "test" {}
`)
}
