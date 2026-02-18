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

// TestAcc_ReportResource_basic tests creating a basic report with required attributes only.
func TestAcc_ReportResource_basic(t *testing.T) {
	resourceName := "armis_report.test"
	rName := strings.ToLower(acctest.RandomWithPrefix("tfacc-report"))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReportResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "report_name", rName),
					resource.TestCheckResourceAttr(resourceName, "asq", `in:devices timeFrame:"1 Day"`),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "is_scheduled"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// email_subject is write-only and not returned by the API on read
				ImportStateVerifyIgnore: []string{"email_subject"},
			},
		},
	})
}

// TestAcc_ReportResource_withSchedule tests creating a report with schedule configuration.
func TestAcc_ReportResource_withSchedule(t *testing.T) {
	resourceName := "armis_report.test"
	rName := strings.ToLower(acctest.RandomWithPrefix("tfacc-report-sched"))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReportResourceConfig_withSchedule(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "report_name", rName),
					resource.TestCheckResourceAttr(resourceName, "asq", `in:devices timeFrame:"7 Days"`),
					resource.TestCheckResourceAttr(resourceName, "email_subject", "Weekly Device Report"),
					resource.TestCheckResourceAttr(resourceName, "is_scheduled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
				),
			},
		},
	})
}

// TestAcc_ReportResource_withExportConfiguration tests creating a report with export configuration.
func TestAcc_ReportResource_withExportConfiguration(t *testing.T) {
	resourceName := "armis_report.test"
	rName := strings.ToLower(acctest.RandomWithPrefix("tfacc-report-exp"))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReportResourceConfig_withExportConfiguration(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "report_name", rName),
					resource.TestCheckResourceAttr(resourceName, "asq", `in:devices timeFrame:"1 Day"`),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
				),
			},
		},
	})
}

func testAccReportResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "armis_report" "test" {
  report_name = %q
  asq         = "in:devices timeFrame:\"1 Day\""
}
`, name)
}

func testAccReportResourceConfig_withSchedule(name string) string {
	return fmt.Sprintf(`
resource "armis_report" "test" {
  report_name   = %q
  asq           = "in:devices timeFrame:\"7 Days\""
  email_subject = "Weekly Device Report"

  schedule = {
    email              = ["test@example.com"]
    repeat_amount      = "1"
    repeat_unit        = "Weeks"
    report_file_format = "csv"
    time_of_day        = "09:00"
    timezone           = "UTC"
    weekdays           = ["Monday"]
  }
}
`, name)
}

func testAccReportResourceConfig_withExportConfiguration(name string) string {
	return fmt.Sprintf(`
resource "armis_report" "test" {
  report_name = %q
  asq         = "in:devices timeFrame:\"1 Day\""

  export_configuration = {
    columns = {
      devices = ["name", "type", "ipAddress"]
    }
  }
}
`, name)
}
