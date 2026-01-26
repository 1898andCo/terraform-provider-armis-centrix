// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package verify_test

import (
	"context"
	"testing"

	"github.com/1898andCo/terraform-provider-armis-centrix/internal/verify"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestValidMitreAttackLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		// Valid formats - Enterprise matrix
		{"valid enterprise with subtechnique", "Enterprise.TA0009.T1056.001", false},
		{"valid enterprise without subtechnique", "Enterprise.TA0009.T1056", false},
		{"valid enterprise different tactic", "Enterprise.TA0001.T1234", false},

		// Valid formats - Mobile matrix
		{"valid mobile with subtechnique", "Mobile.TA0001.T1234.001", false},
		{"valid mobile without subtechnique", "Mobile.TA0001.T1234", false},

		// Valid formats - ICS matrix
		{"valid ics with subtechnique", "ICS.TA0010.T0800.001", false},
		{"valid ics without subtechnique", "ICS.TA0010.T0800", false},

		// Invalid - wrong matrix/domain
		{"lowercase enterprise fails", "enterprise.TA0009.T1056.001", true},
		{"lowercase mobile fails", "mobile.TA0001.T1234", true},
		{"lowercase ics fails", "ics.TA0010.T0800", true},
		{"invalid matrix fails", "Unknown.TA0009.T1056.001", true},
		{"mixed case matrix fails", "ENTERPRISE.TA0009.T1056.001", true},

		// Invalid - wrong tactic format
		{"lowercase tactic fails", "Enterprise.ta0009.T1056.001", true},
		{"tactic missing TA prefix fails", "Enterprise.0009.T1056.001", true},
		{"tactic wrong digit count fails", "Enterprise.TA009.T1056.001", true},
		{"tactic too many digits fails", "Enterprise.TA00091.T1056.001", true},

		// Invalid - wrong technique format
		{"lowercase technique fails", "Enterprise.TA0009.t1056.001", true},
		{"technique missing T prefix fails", "Enterprise.TA0009.1056.001", true},
		{"technique wrong digit count fails", "Enterprise.TA0009.T105.001", true},
		{"technique too many digits fails", "Enterprise.TA0009.T10561.001", true},

		// Invalid - wrong subtechnique format
		{"subtechnique wrong digit count fails", "Enterprise.TA0009.T1056.01", true},
		{"subtechnique too many digits fails", "Enterprise.TA0009.T1056.0011", true},
		{"subtechnique with letters fails", "Enterprise.TA0009.T1056.00a", true},

		// Invalid - missing components
		{"missing technique fails", "Enterprise.TA0009", true},
		{"missing tactic fails", "Enterprise.T1056.001", true},
		{"only matrix fails", "Enterprise", true},

		// Invalid - extra/wrong separators
		{"wrong separator dash fails", "Enterprise-TA0009-T1056-001", true},
		{"wrong separator underscore fails", "Enterprise_TA0009_T1056_001", true},
		{"extra dot at end fails", "Enterprise.TA0009.T1056.001.", true},
		{"leading dot fails", ".Enterprise.TA0009.T1056.001", true},

		// Invalid - edge cases
		{"empty string fails", "", true},
		{"whitespace fails", " ", true},
		{"spaces in value fails", "Enterprise .TA0009.T1056.001", true},
		{"trailing space fails", "Enterprise.TA0009.T1056.001 ", true},
		{"leading space fails", " Enterprise.TA0009.T1056.001", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := verify.ValidMitreAttackLabel()
			req := validator.StringRequest{
				ConfigValue: types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(context.Background(), req, resp)

			if tt.expectError && !resp.Diagnostics.HasError() {
				t.Errorf("expected error for value %q, but got none", tt.value)
			}
			if !tt.expectError && resp.Diagnostics.HasError() {
				t.Errorf("expected no error for value %q, but got: %s", tt.value, resp.Diagnostics.Errors())
			}
		})
	}
}
