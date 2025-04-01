// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_CollectorResource(t *testing.T) {
	resourceName := "armis_collector.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectorResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Collector"),
					resource.TestCheckResourceAttr(resourceName, "deployment_type", "OVA"),
				),
			},
		},
	})
}

func testAccCollectorResourceConfig() string {
	return `
resource "armis_collector" "test" {
  name            = "Test Collector"
  deployment_type = "OVA"
}
`
}
