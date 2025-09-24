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

func TestAcc_PolicyResource(t *testing.T) {
	resourceName := "armis_policy.test"

	rName := strings.ToLower(acctest.RandomWithPrefix("tfacc-policy"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyResourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is an example security policy with all options."),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "ACTIVITY"),
					resource.TestCheckResourceAttr(resourceName, "labels.0", "Security"),
					resource.TestCheckResourceAttr(resourceName, "mitre_attack_labels.0", "Enterprise.TA0009.T1056.001"),
					resource.TestCheckResourceAttr(resourceName, "mitre_attack_labels.1", "Enterprise.TA0009.T1056.004"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.type", "alert"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.params.severity", "high"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.params.title", "Test Security Alert"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.params.type", "Security - Threat"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.params.consolidation.amount", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.params.consolidation.unit", "Hours"),
					resource.TestCheckResourceAttr(resourceName, "rules.and.0", "protocol:BMS"),
				),
			},
			// ImportState testing
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
		},
	})
}

func testAccPolicyResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "armis_policy" "test" {
  name                = %q
  description         = "This is an example security policy with all options."
  enabled             = true
  rule_type           = "ACTIVITY"
  labels              = ["Security"]
  mitre_attack_labels = ["Enterprise.TA0009.T1056.001", "Enterprise.TA0009.T1056.004"]

  actions = [
    {
      type = "alert"
      params = {
        severity = "high"
        title    = "Test Security Alert"
        type     = "Security - Threat"
        consolidation = {
          amount = 2
          unit   = "Hours"
        }
      }
    }
  ]

  rules = {
    and = [
      "protocol:BMS",
    ]
  }
}
`, name)
}
