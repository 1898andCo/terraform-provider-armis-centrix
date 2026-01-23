// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

// Package verify contains the logic for validating schema definitions.
package verify

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ValidMitreAttackLabel validates MITRE ATT&CK label format.
func ValidMitreAttackLabel() validator.String {
	return stringvalidator.RegexMatches(
		regexp.MustCompile(`^(Enterprise|Mobile|ICS)\.TA\d{4}\.T\d{4}(\.\d{3})?$`),
		"must be a valid MITRE ATT&CK label (e.g., Enterprise.TA0009.T1056.001)",
	)
}
