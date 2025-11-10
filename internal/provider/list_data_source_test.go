// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_ListsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccListsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.armis_lists.test", "lists.#"),
					resource.TestCheckResourceAttrSet("data.armis_lists.test", "lists.0.id"),
					resource.TestCheckResourceAttrSet("data.armis_lists.test", "lists.0.name"),
					resource.TestCheckResourceAttrSet("data.armis_lists.test", "lists.0.created_by"),
					resource.TestCheckResourceAttrSet("data.armis_lists.test", "lists.0.list_type"),
					resource.TestCheckResourceAttrSet("data.armis_lists.test", "lists.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.armis_lists.test", "lists.0.last_update_time"),
				),
			},
			{
				Config: testAccListsDataSourceFilteredConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.armis_lists.filtered", "name", "IP Address"),
				),
			},
		},
	})
}

func testAccListsDataSourceConfig() string {
	return `
data "armis_lists" "test" {}
`
}

func testAccListsDataSourceFilteredConfig() string {
	return `
data "armis_lists" "filtered" {
  name = "IP Address"
}
`
}
