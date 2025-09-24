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

func TestAcc_CollectorResource(t *testing.T) {
	resourceName := "armis_collector.test"

	rName := strings.ToLower(acctest.RandomWithPrefix("tfacc-collector"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectorResourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "deployment_type", "OVA"),
				),
			},
		},
	})
}

func testAccCollectorResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "armis_collector" "test" {
  name            = %q
  deployment_type = "OVA"
}
`, name)
}
