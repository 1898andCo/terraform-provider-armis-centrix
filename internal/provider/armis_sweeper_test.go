// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/1898andCo/terraform-provider-armis-centrix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestMain is responsible for parsing the special test flags and invoking the sweepers
//
//	See: https://developer.hashicorp.com/terraform/plugin/testing/acceptance-tests/sweepers
func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func init() {
	resource.AddTestSweepers("collectors", sweep.SweepArmisCollectors("collectors"))
	resource.AddTestSweepers("users", sweep.SweepArmisUsers("users"))
	resource.AddTestSweepers("roles", sweep.SweepArmisRoles("roles"))
	resource.AddTestSweepers("policies", sweep.SweepArmisPolicies("policies"))
	resource.AddTestSweepers("reports", sweep.SweepArmisReports("reports"))
}
