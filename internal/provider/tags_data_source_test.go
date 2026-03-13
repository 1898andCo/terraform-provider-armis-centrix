// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	errUnexpectedPrefix = errors.New("unexpected tag prefix")
	errMissingPrefix    = errors.New("missing expected tag prefix")
)

// Note: match_prefix and exclude_prefix tests assume the Armis test environment
// contains at least one tag starting with "OT" and at least one tag that does not.

// TestAcc_TagsDataSource tests fetching all tags without any filters.
func TestAcc_TagsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.armis_tags.test", "tags.#"),
				),
			},
		},
	})
}

// TestAcc_TagsDataSource_MatchPrefix tests filtering tags by match_prefix.
func TestAcc_TagsDataSource_MatchPrefix(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagsDataSourceMatchPrefixConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.armis_tags.filtered", "tags.#"),
					resource.TestCheckResourceAttr("data.armis_tags.filtered", "match_prefix", "OT"),
					resource.TestCheckResourceAttrWith("data.armis_tags.filtered", "tags.0", func(value string) error {
						if !strings.HasPrefix(value, "OT") {
							return fmt.Errorf("%w: expected \"OT\", got %q", errMissingPrefix, value)
						}
						return nil
					}),
				),
			},
		},
	})
}

// TestAcc_TagsDataSource_ExcludePrefix tests filtering tags by exclude_prefix.
func TestAcc_TagsDataSource_ExcludePrefix(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagsDataSourceExcludePrefixConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.armis_tags.excluded", "tags.#"),
					resource.TestCheckResourceAttr("data.armis_tags.excluded", "exclude_prefix", "OT"),
					resource.TestCheckResourceAttrWith("data.armis_tags.excluded", "tags.0", func(value string) error {
						if strings.HasPrefix(value, "OT") {
							return fmt.Errorf("%w: got %q", errUnexpectedPrefix, value)
						}
						return nil
					}),
				),
			},
		},
	})
}

// TestAcc_TagsDataSource_CombinedFilters tests using both match_prefix and exclude_prefix together.
func TestAcc_TagsDataSource_CombinedFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagsDataSourceCombinedConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.armis_tags.combined", "tags.#"),
					resource.TestCheckResourceAttr("data.armis_tags.combined", "match_prefix", "OT"),
					resource.TestCheckResourceAttr("data.armis_tags.combined", "exclude_prefix", "OT-"),
					resource.TestCheckResourceAttrWith("data.armis_tags.combined", "tags.0", func(value string) error {
						if !strings.HasPrefix(value, "OT") {
							return fmt.Errorf("%w: expected \"OT\", got %q", errMissingPrefix, value)
						}
						if strings.HasPrefix(value, "OT-") {
							return fmt.Errorf("%w: got %q", errUnexpectedPrefix, value)
						}
						return nil
					}),
				),
			},
		},
	})
}

func testAccTagsDataSourceConfig() string {
	return `
data "armis_tags" "test" {}
`
}

func testAccTagsDataSourceMatchPrefixConfig() string {
	return `
data "armis_tags" "filtered" {
  match_prefix = "OT"
}
`
}

func testAccTagsDataSourceExcludePrefixConfig() string {
	return `
data "armis_tags" "excluded" {
  exclude_prefix = "OT"
}
`
}

func testAccTagsDataSourceCombinedConfig() string {
	return `
data "armis_tags" "combined" {
  match_prefix   = "OT"
  exclude_prefix = "OT-"
}
`
}
