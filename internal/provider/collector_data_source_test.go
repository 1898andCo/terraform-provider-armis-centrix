// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_CollectorDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectorDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.#"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.name"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.boot_time"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.city"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.cluster_id"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.collector_number"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.country"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.default_gateway"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.ip_address"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.last_seen"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.mac_address"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.product_serial"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.status"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.subnet"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.system_vendor"),
					resource.TestCheckResourceAttrSet("data.armis_collectors.test", "collectors.0.type"),
				),
			},
		},
	})
}

func testAccCollectorDataSourceConfig() string {
	return fmt.Sprintf(`
data "armis_collectors" "test" {}
`)
}
