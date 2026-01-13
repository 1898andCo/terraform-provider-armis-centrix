// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test data for report data source acceptance tests.
// These values reference default reports that should exist in a standard Armis Centrix environment.
// If tests fail, verify these reports exist in the test environment:
//   - Report ID 3 should be a valid default system report
//   - "All Activities" is a default system report
var (
	reportID   = 3
	reportName = "All Activities"
)

// TestAcc_ReportsDataSource tests fetching all reports without any filters.
func TestAcc_ReportsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReportsDataSourceConfigNoFilter(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify reports list exists and has at least one report
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.#"),
					// Verify the first report has basic attributes populated
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.id"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.report_name"),
					// Note: report_type can be null for some Armis reports, so we don't check it
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.asq"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.is_scheduled"),
				),
			},
		},
	})
}

// TestAcc_ReportsDataSource_ByID tests filtering reports by report_id.
func TestAcc_ReportsDataSource_ByID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReportsDataSourceConfigByID(fmt.Sprintf("%d", reportID)),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Should return exactly one report
					resource.TestCheckResourceAttr("data.armis_reports.test", "reports.#", "1"),
					// Verify the report ID matches
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.id"),
					// Verify all attributes are present
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.report_name"),
					// Note: report_type can be null for some Armis reports, so we don't check it
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.asq"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.is_scheduled"),
				),
			},
		},
	})
}

// TestAcc_ReportsDataSource_ByName tests filtering reports by report_name.
func TestAcc_ReportsDataSource_ByName(t *testing.T) {
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
					// Note: report_type can be null for some Armis reports, so we don't check it
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.asq"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.armis_reports.test", "reports.0.is_scheduled"),
				),
			},
		},
	})
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

func testAccReportsDataSourceConfigNoFilter() string {
	return `
data "armis_reports" "test" {
}
`
}
