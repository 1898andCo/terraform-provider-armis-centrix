// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAcc_ReportsDataSource tests fetching all reports with comprehensive attribute validation.
func TestAcc_ReportsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReportsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify reports list exists and has at least one report
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.#"),
					// Verify report attributes are populated
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.id"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.report_name"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.report_type"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.asq"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.is_scheduled"),
				),
			},
		},
	})
}

// TestAcc_ReportsDataSource_ByID tests filtering reports by report_id.
// This test requires TEST_REPORT_ID environment variable to be set.
func TestAcc_ReportsDataSource_ByID(t *testing.T) {
	reportID := os.Getenv("TEST_REPORT_ID")
	if reportID == "" {
		t.Skip("TEST_REPORT_ID environment variable not set, skipping report_id filter test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReportsDataSourceConfigByID(reportID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Should return exactly one report
					resource.TestCheckResourceAttr("data.armis_reports.test", "reports.#", "1"),
					// Verify the report ID matches
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.id"),
					// Verify all attributes are present
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.report_name"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.report_type"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.asq"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.is_scheduled"),
				),
			},
		},
	})
}

// TestAcc_ReportsDataSource_ByName tests filtering reports by report_name.
// This test requires TEST_REPORT_NAME environment variable to be set.
func TestAcc_ReportsDataSource_ByName(t *testing.T) {
	reportName := os.Getenv("TEST_REPORT_NAME")
	if reportName == "" {
		t.Skip("TEST_REPORT_NAME environment variable not set, skipping report_name filter test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReportsDataSourceConfigByName(reportName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Should return at least one report with matching name
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.#"),
					// Verify the first report has the expected name
					resource.TestCheckResourceAttr("data.armis_reports.test", "reports.0.report_name", reportName),
					// Verify all attributes are present
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.id"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.report_type"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.asq"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.is_scheduled"),
				),
			},
		},
	})
}

func testAccReportsDataSourceConfig() string {
	return `data "armis_reports" "test" {}`
}

func testAccReportsDataSourceConfigByID(reportID string) string {
	return `
data "armis_reports" "test" {
  report_id = "` + reportID + `"
}
`
}

func testAccReportsDataSourceConfigByName(reportName string) string {
	return `
data "armis_reports" "test" {
  report_name = "` + reportName + `"
}
`
}
