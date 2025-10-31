// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/1898andCo/terraform-provider-armis-centrix/armis"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestShouldIncludePolicy tests the ShouldIncludePolicy filter function.
func TestShouldIncludePolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		model    PolicyDataSourcePolicyModel
		prefix   types.String
		expected bool
	}{
		{
			name: "null prefix should include all",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringValue("AdminPolicy"),
			},
			prefix:   types.StringNull(),
			expected: true,
		},
		{
			name: "unknown prefix should include all",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringValue("AdminPolicy"),
			},
			prefix:   types.StringUnknown(),
			expected: true,
		},
		{
			name: "empty prefix should include all",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringValue("AdminPolicy"),
			},
			prefix:   types.StringValue(""),
			expected: true,
		},
		{
			name: "matching prefix should include",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringValue("AdminPolicy"),
			},
			prefix:   types.StringValue("Admin"),
			expected: true,
		},
		{
			name: "non-matching prefix should not include",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringValue("UserPolicy"),
			},
			prefix:   types.StringValue("Admin"),
			expected: false,
		},
		{
			name: "exact match should include",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringValue("ReadOnly"),
			},
			prefix:   types.StringValue("ReadOnly"),
			expected: true,
		},
		{
			name: "null model name should not include",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringNull(),
			},
			prefix:   types.StringValue("Admin"),
			expected: false,
		},
		{
			name: "unknown model name should not include",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringUnknown(),
			},
			prefix:   types.StringValue("Admin"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ShouldIncludePolicy(tt.model, tt.prefix)
			if result != tt.expected {
				t.Errorf("ShouldIncludePolicy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestShouldExcludePolicy tests the ShouldExcludePolicy filter function.
func TestShouldExcludePolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		model    PolicyDataSourcePolicyModel
		prefix   types.String
		expected bool
	}{
		{
			name: "null prefix should not exclude",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringValue("AdminPolicy"),
			},
			prefix:   types.StringNull(),
			expected: false,
		},
		{
			name: "unknown prefix should not exclude",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringValue("AdminPolicy"),
			},
			prefix:   types.StringUnknown(),
			expected: false,
		},
		{
			name: "empty prefix should not exclude",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringValue("AdminPolicy"),
			},
			prefix:   types.StringValue(""),
			expected: false,
		},
		{
			name: "matching prefix should exclude",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringValue("TestPolicy"),
			},
			prefix:   types.StringValue("Test"),
			expected: true,
		},
		{
			name: "non-matching prefix should not exclude",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringValue("UserPolicy"),
			},
			prefix:   types.StringValue("Admin"),
			expected: false,
		},
		{
			name: "exact match should exclude",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringValue("System"),
			},
			prefix:   types.StringValue("System"),
			expected: true,
		},
		{
			name: "null model name should not exclude",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringNull(),
			},
			prefix:   types.StringValue("Test"),
			expected: false,
		},
		{
			name: "unknown model name should not exclude",
			model: PolicyDataSourcePolicyModel{
				Name: types.StringUnknown(),
			},
			prefix:   types.StringValue("Test"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ShouldExcludePolicy(tt.model, tt.prefix)
			if result != tt.expected {
				t.Errorf("ShouldExcludePolicy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestConvertListToStringSlice tests the ConvertListToStringSlice function.
func TestConvertListToStringSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    types.List
		expected []string
	}{
		{
			name:     "null list should return nil",
			input:    types.ListNull(types.StringType),
			expected: nil,
		},
		{
			name:     "unknown list should return nil",
			input:    types.ListUnknown(types.StringType),
			expected: nil,
		},
		{
			name: "valid list with strings",
			input: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("item1"),
				types.StringValue("item2"),
				types.StringValue("item3"),
			}),
			expected: []string{"item1", "item2", "item3"},
		},
		{
			name:     "empty list should return empty slice",
			input:    types.ListValueMust(types.StringType, []attr.Value{}),
			expected: []string{},
		},
		{
			name: "list with null element should skip it",
			input: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("item1"),
				types.StringNull(),
				types.StringValue("item2"),
			}),
			expected: []string{"item1", "item2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ConvertListToStringSlice(tt.input)

			// Compare length
			if (result == nil) != (tt.expected == nil) {
				t.Errorf("ConvertListToStringSlice() = %v, want %v", result, tt.expected)
				return
			}

			if result != nil && tt.expected != nil {
				if len(result) != len(tt.expected) {
					t.Errorf("ConvertListToStringSlice() length = %d, want %d", len(result), len(tt.expected))
					return
				}

				// Compare elements
				for i := range result {
					if result[i] != tt.expected[i] {
						t.Errorf("ConvertListToStringSlice()[%d] = %s, want %s", i, result[i], tt.expected[i])
					}
				}
			}
		})
	}
}

// TestConvertStringSliceToList tests the ConvertStringSliceToList function.
func TestConvertStringSliceToList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []string
		validate func(t *testing.T, result types.List)
	}{
		{
			name:  "nil slice should return null list",
			input: nil,
			validate: func(t *testing.T, result types.List) {
				if !result.IsNull() {
					t.Error("Expected null list")
				}
			},
		},
		{
			name:  "empty slice should return empty list",
			input: []string{},
			validate: func(t *testing.T, result types.List) {
				if result.IsNull() || result.IsUnknown() {
					t.Error("Expected non-null, non-unknown list")
				}
				elements := result.Elements()
				if len(elements) != 0 {
					t.Errorf("Expected 0 elements, got %d", len(elements))
				}
			},
		},
		{
			name:  "valid slice with strings",
			input: []string{"item1", "item2", "item3"},
			validate: func(t *testing.T, result types.List) {
				if result.IsNull() || result.IsUnknown() {
					t.Error("Expected non-null, non-unknown list")
				}
				elements := result.Elements()
				if len(elements) != 3 {
					t.Errorf("Expected 3 elements, got %d", len(elements))
					return
				}

				expected := []string{"item1", "item2", "item3"}
				for i, elem := range elements {
					if strVal, ok := elem.(types.String); ok {
						if strVal.ValueString() != expected[i] {
							t.Errorf("Element[%d] = %s, want %s", i, strVal.ValueString(), expected[i])
						}
					} else {
						t.Errorf("Element[%d] is not a string", i)
					}
				}
			},
		},
		{
			name:  "single element slice",
			input: []string{"single"},
			validate: func(t *testing.T, result types.List) {
				elements := result.Elements()
				if len(elements) != 1 {
					t.Errorf("Expected 1 element, got %d", len(elements))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ConvertStringSliceToList(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestConvertListToStringSlice_RoundTrip tests bidirectional conversion.
func TestConvertListToStringSlice_RoundTrip(t *testing.T) {
	t.Parallel()

	original := []string{"alpha", "beta", "gamma", "delta"}

	// Convert to list
	list := ConvertStringSliceToList(original)

	// Convert back to slice
	result := ConvertListToStringSlice(list)

	// Verify round trip
	if len(result) != len(original) {
		t.Errorf("Round trip failed: length changed from %d to %d", len(original), len(result))
		return
	}

	for i := range original {
		if result[i] != original[i] {
			t.Errorf("Round trip failed: element[%d] changed from %s to %s", i, original[i], result[i])
		}
	}
}

// TestConvertStringsToTypeStrings tests the convertStringsToTypeStrings function.
func TestConvertStringsToTypeStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []string
		validate func(t *testing.T, result []types.String)
	}{
		{
			name:  "nil slice should return nil",
			input: nil,
			validate: func(t *testing.T, result []types.String) {
				if result != nil {
					t.Error("Expected nil result")
				}
			},
		},
		{
			name:  "empty slice should return empty slice",
			input: []string{},
			validate: func(t *testing.T, result []types.String) {
				if len(result) != 0 {
					t.Errorf("Expected length 0, got %d", len(result))
				}
			},
		},
		{
			name:  "valid strings converted",
			input: []string{"label1", "label2", "label3"},
			validate: func(t *testing.T, result []types.String) {
				if len(result) != 3 {
					t.Errorf("Expected length 3, got %d", len(result))
				}
				expected := []string{"label1", "label2", "label3"}
				for i, str := range result {
					if str.ValueString() != expected[i] {
						t.Errorf("Expected result[%d] = %s, got %s", i, expected[i], str.ValueString())
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := convertStringsToTypeStrings(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestConvertInterfacesToTypeStrings tests the convertInterfacesToTypeStrings function.
func TestConvertInterfacesToTypeStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []any
		validate func(t *testing.T, result []types.String)
	}{
		{
			name:  "nil slice should return nil",
			input: nil,
			validate: func(t *testing.T, result []types.String) {
				if result != nil {
					t.Error("Expected nil result")
				}
			},
		},
		{
			name:  "empty slice should return empty slice",
			input: []any{},
			validate: func(t *testing.T, result []types.String) {
				if len(result) != 0 {
					t.Errorf("Expected length 0, got %d", len(result))
				}
			},
		},
		{
			name:  "valid interface strings converted",
			input: []any{"rule1", "rule2", "rule3"},
			validate: func(t *testing.T, result []types.String) {
				if len(result) != 3 {
					t.Errorf("Expected length 3, got %d", len(result))
				}
				expected := []string{"rule1", "rule2", "rule3"}
				for i, str := range result {
					if str.ValueString() != expected[i] {
						t.Errorf("Expected result[%d] = %s, got %s", i, expected[i], str.ValueString())
					}
				}
			},
		},
		{
			name:  "non-string interfaces skipped",
			input: []any{"valid", 123, "another", nil, "third"},
			validate: func(t *testing.T, result []types.String) {
				if len(result) != 3 {
					t.Errorf("Expected length 3 (non-strings skipped), got %d", len(result))
				}
				expected := []string{"valid", "another", "third"}
				for i, str := range result {
					if str.ValueString() != expected[i] {
						t.Errorf("Expected result[%d] = %s, got %s", i, expected[i], str.ValueString())
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := convertInterfacesToTypeStrings(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestConvertSliceToList tests the ConvertSliceToList function.
func TestConvertSliceToList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []any
		validate func(t *testing.T, result types.List)
	}{
		{
			name:  "nil slice should return null list",
			input: nil,
			validate: func(t *testing.T, result types.List) {
				if !result.IsNull() {
					t.Error("Expected null list")
				}
			},
		},
		{
			name:  "empty slice should return empty list",
			input: []any{},
			validate: func(t *testing.T, result types.List) {
				if result.IsNull() || result.IsUnknown() {
					t.Error("Expected non-null, non-unknown list")
				}
				if len(result.Elements()) != 0 {
					t.Errorf("Expected 0 elements, got %d", len(result.Elements()))
				}
			},
		},
		{
			name:  "valid string interfaces converted",
			input: []any{"item1", "item2", "item3"},
			validate: func(t *testing.T, result types.List) {
				elements := result.Elements()
				if len(elements) != 3 {
					t.Errorf("Expected 3 elements, got %d", len(elements))
				}
			},
		},
		{
			name:  "non-string interfaces skipped",
			input: []any{"valid", 123, "another"},
			validate: func(t *testing.T, result types.List) {
				elements := result.Elements()
				if len(elements) != 2 {
					t.Errorf("Expected 2 elements (non-strings skipped), got %d", len(elements))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ConvertSliceToList(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestConvertMitreLabelsToDataSource tests the convertMitreLabelsToDataSource function.
func TestConvertMitreLabelsToDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []armis.MitreAttackLabel
		validate func(t *testing.T, result []PolicyDataSourceMitreLabelModel)
	}{
		{
			name:  "empty labels should return empty slice",
			input: []armis.MitreAttackLabel{},
			validate: func(t *testing.T, result []PolicyDataSourceMitreLabelModel) {
				if len(result) != 0 {
					t.Errorf("Expected length 0, got %d", len(result))
				}
			},
		},
		{
			name: "valid labels converted correctly",
			input: []armis.MitreAttackLabel{
				{
					Matrix:       "Enterprise",
					SubTechnique: "T1566.001",
					Tactic:       "Initial Access",
					Technique:    "Phishing",
				},
				{
					Matrix:       "Mobile",
					SubTechnique: "T1234.567",
					Tactic:       "Persistence",
					Technique:    "Malware",
				},
			},
			validate: func(t *testing.T, result []PolicyDataSourceMitreLabelModel) {
				if len(result) != 2 {
					t.Errorf("Expected length 2, got %d", len(result))
				}
				if result[0].Matrix.ValueString() != "Enterprise" {
					t.Errorf("Expected Matrix 'Enterprise', got '%s'", result[0].Matrix.ValueString())
				}
				if result[0].Technique.ValueString() != "Phishing" {
					t.Errorf("Expected Technique 'Phishing', got '%s'", result[0].Technique.ValueString())
				}
				if result[1].Tactic.ValueString() != "Persistence" {
					t.Errorf("Expected Tactic 'Persistence', got '%s'", result[1].Tactic.ValueString())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := convertMitreLabelsToDataSource(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestConvertActionToDataSource tests the convertActionToDataSource function.
func TestConvertActionToDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    armis.Action
		validate func(t *testing.T, result PolicyDataSourceActionModel)
	}{
		{
			name: "action with all fields",
			input: armis.Action{
				Type: "alert",
				Params: armis.Params{
					Severity: "high",
					Title:    "Security Alert",
					Type:     "notification",
					Endpoint: "https://example.com/webhook",
					Tags:     []string{"critical", "security"},
					Consolidation: armis.Consolidation{
						Amount: 10,
						Unit:   "minutes",
					},
				},
			},
			validate: func(t *testing.T, result PolicyDataSourceActionModel) {
				if result.Type.ValueString() != "alert" {
					t.Errorf("Expected Type 'alert', got '%s'", result.Type.ValueString())
				}
				if result.Params.Severity.ValueString() != "high" {
					t.Errorf("Expected Severity 'high', got '%s'", result.Params.Severity.ValueString())
				}
				if result.Params.Title.ValueString() != "Security Alert" {
					t.Errorf("Expected Title 'Security Alert', got '%s'", result.Params.Title.ValueString())
				}
				if result.Params.Consolidation.Amount.ValueInt64() != 10 {
					t.Errorf("Expected Amount 10, got %d", result.Params.Consolidation.Amount.ValueInt64())
				}
				if result.Params.Consolidation.Unit.ValueString() != "minutes" {
					t.Errorf("Expected Unit 'minutes', got '%s'", result.Params.Consolidation.Unit.ValueString())
				}
				if len(result.Params.Tags) != 2 {
					t.Errorf("Expected 2 tags, got %d", len(result.Params.Tags))
				}
			},
		},
		{
			name: "action with minimal fields",
			input: armis.Action{
				Type: "log",
				Params: armis.Params{
					Severity: "low",
				},
			},
			validate: func(t *testing.T, result PolicyDataSourceActionModel) {
				if result.Type.ValueString() != "log" {
					t.Errorf("Expected Type 'log', got '%s'", result.Type.ValueString())
				}
				if result.Params.Severity.ValueString() != "low" {
					t.Errorf("Expected Severity 'low', got '%s'", result.Params.Severity.ValueString())
				}
				if !result.Params.Consolidation.Amount.IsNull() {
					t.Error("Expected null consolidation amount")
				}
				if !result.Params.Consolidation.Unit.IsNull() {
					t.Error("Expected null consolidation unit")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := convertActionToDataSource(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestBuildPolicyDataSourceModelFromGet tests the BuildPolicyDataSourceModelFromGet function.
func TestBuildPolicyDataSourceModelFromGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		policy   armis.GetPolicySettings
		id       string
		validate func(t *testing.T, result PolicyDataSourcePolicyModel)
	}{
		{
			name: "comprehensive policy with all fields",
			policy: armis.GetPolicySettings{
				Name:        "Security Policy",
				Description: "Monitors security threats",
				IsEnabled:   true,
				RuleType:    "ACTIVITY",
				Labels:      []string{"production", "security"},
				MitreAttackLabels: []armis.MitreAttackLabel{
					{
						Matrix:       "Enterprise",
						SubTechnique: "T1566.001",
						Tactic:       "Initial Access",
						Technique:    "Phishing",
					},
				},
				Actions: []armis.Action{
					{
						Type: "alert",
						Params: armis.Params{
							Severity: "high",
							Title:    "Security Alert",
						},
					},
				},
				Rules: armis.Rules{
					And: []any{"rule1", "rule2"},
					Or:  []any{"rule3"},
				},
			},
			id: "policy-123",
			validate: func(t *testing.T, result PolicyDataSourcePolicyModel) {
				if result.ID.ValueString() != "policy-123" {
					t.Errorf("Expected ID 'policy-123', got '%s'", result.ID.ValueString())
				}
				if result.Name.ValueString() != "Security Policy" {
					t.Errorf("Expected Name 'Security Policy', got '%s'", result.Name.ValueString())
				}
				if result.Description.ValueString() != "Monitors security threats" {
					t.Errorf("Expected Description, got '%s'", result.Description.ValueString())
				}
				if !result.IsEnabled.ValueBool() {
					t.Error("Expected IsEnabled to be true")
				}
				if result.RuleType.ValueString() != "ACTIVITY" {
					t.Errorf("Expected RuleType 'ACTIVITY', got '%s'", result.RuleType.ValueString())
				}
				if len(result.Labels) != 2 {
					t.Errorf("Expected 2 labels, got %d", len(result.Labels))
				}
				if len(result.MitreAttackLabels) != 1 {
					t.Errorf("Expected 1 MITRE label, got %d", len(result.MitreAttackLabels))
				}
				if len(result.Actions) != 1 {
					t.Errorf("Expected 1 action, got %d", len(result.Actions))
				}
				if len(result.Rules.And) != 2 {
					t.Errorf("Expected 2 And rules, got %d", len(result.Rules.And))
				}
			},
		},
		{
			name: "policy with empty id should have null ID",
			policy: armis.GetPolicySettings{
				Name:      "Test Policy",
				IsEnabled: false,
			},
			id: "",
			validate: func(t *testing.T, result PolicyDataSourcePolicyModel) {
				if !result.ID.IsNull() {
					t.Error("Expected null ID when id string is empty")
				}
				if result.Name.ValueString() != "Test Policy" {
					t.Errorf("Expected Name 'Test Policy', got '%s'", result.Name.ValueString())
				}
				if result.IsEnabled.ValueBool() {
					t.Error("Expected IsEnabled to be false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := BuildPolicyDataSourceModelFromGet(tt.policy, tt.id)
			tt.validate(t, result)
		})
	}
}

// TestBuildPolicyDataSourceModelFromSingle tests the BuildPolicyDataSourceModelFromSingle function.
func TestBuildPolicyDataSourceModelFromSingle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		policy   armis.SinglePolicy
		validate func(t *testing.T, result PolicyDataSourcePolicyModel)
	}{
		{
			name: "policy with Actions array",
			policy: armis.SinglePolicy{
				ID:          "single-123",
				Name:        "Single Policy",
				Description: "Test policy",
				IsEnabled:   true,
				RuleType:    "DEVICE",
				Actions: []armis.Action{
					{
						Type: "alert",
						Params: armis.Params{
							Severity: "medium",
						},
					},
				},
				Rules: armis.Rules{
					Or: []any{"deviceType:Laptop"},
				},
			},
			validate: func(t *testing.T, result PolicyDataSourcePolicyModel) {
				if result.ID.ValueString() != "single-123" {
					t.Errorf("Expected ID 'single-123', got '%s'", result.ID.ValueString())
				}
				if result.Name.ValueString() != "Single Policy" {
					t.Errorf("Expected Name 'Single Policy', got '%s'", result.Name.ValueString())
				}
				if len(result.Actions) != 1 {
					t.Errorf("Expected 1 action, got %d", len(result.Actions))
				}
			},
		},
		{
			name: "policy with Action field (fallback)",
			policy: armis.SinglePolicy{
				ID:        "single-456",
				Name:      "Fallback Policy",
				IsEnabled: false,
				Action: armis.Action{
					Type: "log",
					Params: armis.Params{
						Type: "system",
					},
				},
				Rules: armis.Rules{},
			},
			validate: func(t *testing.T, result PolicyDataSourcePolicyModel) {
				if result.ID.ValueString() != "single-456" {
					t.Errorf("Expected ID 'single-456', got '%s'", result.ID.ValueString())
				}
				// When Actions is empty but Action has data, it should be included
				if len(result.Actions) != 1 {
					t.Errorf("Expected 1 action (from Action field), got %d", len(result.Actions))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := BuildPolicyDataSourceModelFromSingle(tt.policy)
			tt.validate(t, result)
		})
	}
}

// TestStringValue tests the stringValue helper function.
func TestStringValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          types.String
		expectedValue  string
		expectedExists bool
	}{
		{
			name:           "null string returns false",
			input:          types.StringNull(),
			expectedValue:  "",
			expectedExists: false,
		},
		{
			name:           "unknown string returns false",
			input:          types.StringUnknown(),
			expectedValue:  "",
			expectedExists: false,
		},
		{
			name:           "valid string returns value and true",
			input:          types.StringValue("test"),
			expectedValue:  "test",
			expectedExists: true,
		},
		{
			name:           "empty string value returns empty and true",
			input:          types.StringValue(""),
			expectedValue:  "",
			expectedExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			value, exists := stringValue(tt.input)
			if value != tt.expectedValue {
				t.Errorf("Expected value '%s', got '%s'", tt.expectedValue, value)
			}
			if exists != tt.expectedExists {
				t.Errorf("Expected exists %v, got %v", tt.expectedExists, exists)
			}
		})
	}
}

// TestIntValue tests the intValue helper function.
func TestIntValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          types.Int64
		expectedValue  int
		expectedExists bool
	}{
		{
			name:           "null int returns false",
			input:          types.Int64Null(),
			expectedValue:  0,
			expectedExists: false,
		},
		{
			name:           "unknown int returns false",
			input:          types.Int64Unknown(),
			expectedValue:  0,
			expectedExists: false,
		},
		{
			name:           "valid int returns value and true",
			input:          types.Int64Value(42),
			expectedValue:  42,
			expectedExists: true,
		},
		{
			name:           "zero value returns zero and true",
			input:          types.Int64Value(0),
			expectedValue:  0,
			expectedExists: true,
		},
		{
			name:           "negative value returns negative and true",
			input:          types.Int64Value(-10),
			expectedValue:  -10,
			expectedExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			value, exists := intValue(tt.input)
			if value != tt.expectedValue {
				t.Errorf("Expected value %d, got %d", tt.expectedValue, value)
			}
			if exists != tt.expectedExists {
				t.Errorf("Expected exists %v, got %v", tt.expectedExists, exists)
			}
		})
	}
}

// TestConsolidationFromObject tests the consolidationFromObject function.
func TestConsolidationFromObject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        types.Object
		validate     func(t *testing.T, result armis.Consolidation, hasValue bool, diags diag.Diagnostics)
	}{
		{
			name:  "null object returns false",
			input: types.ObjectNull(map[string]attr.Type{
				"amount": types.Int64Type,
				"unit":   types.StringType,
			}),
			validate: func(t *testing.T, result armis.Consolidation, hasValue bool, diags diag.Diagnostics) {
				if hasValue {
					t.Error("Expected hasValue to be false for null object")
				}
				if diags.HasError() {
					t.Error("Expected no diagnostics errors")
				}
			},
		},
		{
			name:  "unknown object returns false",
			input: types.ObjectUnknown(map[string]attr.Type{
				"amount": types.Int64Type,
				"unit":   types.StringType,
			}),
			validate: func(t *testing.T, result armis.Consolidation, hasValue bool, diags diag.Diagnostics) {
				if hasValue {
					t.Error("Expected hasValue to be false for unknown object")
				}
				if diags.HasError() {
					t.Error("Expected no diagnostics errors")
				}
			},
		},
		{
			name: "valid consolidation with both fields",
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"amount": types.Int64Type,
					"unit":   types.StringType,
				},
				map[string]attr.Value{
					"amount": types.Int64Value(10),
					"unit":   types.StringValue("minutes"),
				},
			),
			validate: func(t *testing.T, result armis.Consolidation, hasValue bool, diags diag.Diagnostics) {
				if !hasValue {
					t.Error("Expected hasValue to be true")
				}
				if result.Amount != 10 {
					t.Errorf("Expected Amount 10, got %d", result.Amount)
				}
				if result.Unit != "minutes" {
					t.Errorf("Expected Unit 'minutes', got '%s'", result.Unit)
				}
				if diags.HasError() {
					t.Error("Expected no diagnostics errors")
				}
			},
		},
		{
			name: "consolidation with only amount",
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"amount": types.Int64Type,
					"unit":   types.StringType,
				},
				map[string]attr.Value{
					"amount": types.Int64Value(5),
					"unit":   types.StringNull(),
				},
			),
			validate: func(t *testing.T, result armis.Consolidation, hasValue bool, diags diag.Diagnostics) {
				if !hasValue {
					t.Error("Expected hasValue to be true")
				}
				if result.Amount != 5 {
					t.Errorf("Expected Amount 5, got %d", result.Amount)
				}
				if result.Unit != "" {
					t.Errorf("Expected empty Unit, got '%s'", result.Unit)
				}
			},
		},
		{
			name: "consolidation with only unit",
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"amount": types.Int64Type,
					"unit":   types.StringType,
				},
				map[string]attr.Value{
					"amount": types.Int64Null(),
					"unit":   types.StringValue("hours"),
				},
			),
			validate: func(t *testing.T, result armis.Consolidation, hasValue bool, diags diag.Diagnostics) {
				if !hasValue {
					t.Error("Expected hasValue to be true")
				}
				if result.Amount != 0 {
					t.Errorf("Expected Amount 0, got %d", result.Amount)
				}
				if result.Unit != "hours" {
					t.Errorf("Expected Unit 'hours', got '%s'", result.Unit)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, hasValue, diags := consolidationFromObject(tt.input)
			tt.validate(t, result, hasValue, diags)
		})
	}
}

// TestParamsFromObject tests the paramsFromObject function.
func TestParamsFromObject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    types.Object
		validate func(t *testing.T, result armis.Params, hasParams bool, diags diag.Diagnostics)
	}{
		{
			name: "null object returns false",
			input: types.ObjectNull(map[string]attr.Type{
				"severity":      types.StringType,
				"title":         types.StringType,
				"type":          types.StringType,
				"endpoint":      types.StringType,
				"tags":          types.ListType{ElemType: types.StringType},
				"consolidation": types.ObjectType{},
			}),
			validate: func(t *testing.T, result armis.Params, hasParams bool, diags diag.Diagnostics) {
				if hasParams {
					t.Error("Expected hasParams to be false for null object")
				}
			},
		},
		{
			name: "params with all fields",
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"severity": types.StringType,
					"title":    types.StringType,
					"type":     types.StringType,
					"endpoint": types.StringType,
					"tags":     types.ListType{ElemType: types.StringType},
					"consolidation": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"amount": types.Int64Type,
							"unit":   types.StringType,
						},
					},
				},
				map[string]attr.Value{
					"severity": types.StringValue("high"),
					"title":    types.StringValue("Security Alert"),
					"type":     types.StringValue("notification"),
					"endpoint": types.StringValue("https://example.com/webhook"),
					"tags": types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("critical"),
						types.StringValue("security"),
					}),
					"consolidation": types.ObjectValueMust(
						map[string]attr.Type{
							"amount": types.Int64Type,
							"unit":   types.StringType,
						},
						map[string]attr.Value{
							"amount": types.Int64Value(10),
							"unit":   types.StringValue("minutes"),
						},
					),
				},
			),
			validate: func(t *testing.T, result armis.Params, hasParams bool, diags diag.Diagnostics) {
				if !hasParams {
					t.Error("Expected hasParams to be true")
				}
				if result.Severity != "high" {
					t.Errorf("Expected Severity 'high', got '%s'", result.Severity)
				}
				if result.Title != "Security Alert" {
					t.Errorf("Expected Title 'Security Alert', got '%s'", result.Title)
				}
				if result.Type != "notification" {
					t.Errorf("Expected Type 'notification', got '%s'", result.Type)
				}
				if result.Endpoint != "https://example.com/webhook" {
					t.Errorf("Expected Endpoint, got '%s'", result.Endpoint)
				}
				if len(result.Tags) != 2 {
					t.Errorf("Expected 2 tags, got %d", len(result.Tags))
				}
				if result.Consolidation.Amount != 10 {
					t.Errorf("Expected Consolidation.Amount 10, got %d", result.Consolidation.Amount)
				}
				if result.Consolidation.Unit != "minutes" {
					t.Errorf("Expected Consolidation.Unit 'minutes', got '%s'", result.Consolidation.Unit)
				}
			},
		},
		{
			name: "params with minimal fields",
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"severity": types.StringType,
					"title":    types.StringType,
					"type":     types.StringType,
					"endpoint": types.StringType,
					"tags":     types.ListType{ElemType: types.StringType},
					"consolidation": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"amount": types.Int64Type,
							"unit":   types.StringType,
						},
					},
				},
				map[string]attr.Value{
					"severity": types.StringValue("low"),
					"title":    types.StringNull(),
					"type":     types.StringNull(),
					"endpoint": types.StringNull(),
					"tags":     types.ListNull(types.StringType),
					"consolidation": types.ObjectNull(map[string]attr.Type{
						"amount": types.Int64Type,
						"unit":   types.StringType,
					}),
				},
			),
			validate: func(t *testing.T, result armis.Params, hasParams bool, diags diag.Diagnostics) {
				if !hasParams {
					t.Error("Expected hasParams to be true (has severity)")
				}
				if result.Severity != "low" {
					t.Errorf("Expected Severity 'low', got '%s'", result.Severity)
				}
				if result.Title != "" {
					t.Errorf("Expected empty Title, got '%s'", result.Title)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, hasParams, diags := paramsFromObject(tt.input)
			tt.validate(t, result, hasParams, diags)
		})
	}
}

// TestConvertActionModel tests the convertActionModel function.
func TestConvertActionModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    ActionModel
		validate func(t *testing.T, result armis.Action, diags diag.Diagnostics)
	}{
		{
			name: "action with type and params",
			input: ActionModel{
				Type: types.StringValue("alert"),
				Params: types.ObjectValueMust(
					map[string]attr.Type{
						"severity": types.StringType,
						"title":    types.StringType,
						"type":     types.StringType,
						"endpoint": types.StringType,
						"tags":     types.ListType{ElemType: types.StringType},
						"consolidation": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"amount": types.Int64Type,
								"unit":   types.StringType,
							},
						},
					},
					map[string]attr.Value{
						"severity": types.StringValue("high"),
						"title":    types.StringValue("Test Alert"),
						"type":     types.StringNull(),
						"endpoint": types.StringNull(),
						"tags":     types.ListNull(types.StringType),
						"consolidation": types.ObjectNull(map[string]attr.Type{
							"amount": types.Int64Type,
							"unit":   types.StringType,
						}),
					},
				),
			},
			validate: func(t *testing.T, result armis.Action, diags diag.Diagnostics) {
				if result.Type != "alert" {
					t.Errorf("Expected Type 'alert', got '%s'", result.Type)
				}
				if result.Params.Severity != "high" {
					t.Errorf("Expected Severity 'high', got '%s'", result.Params.Severity)
				}
				if result.Params.Title != "Test Alert" {
					t.Errorf("Expected Title 'Test Alert', got '%s'", result.Params.Title)
				}
				if diags.HasError() {
					t.Error("Expected no diagnostics errors")
				}
			},
		},
		{
			name: "action with null params",
			input: ActionModel{
				Type: types.StringValue("log"),
				Params: types.ObjectNull(map[string]attr.Type{
					"severity":      types.StringType,
					"title":         types.StringType,
					"type":          types.StringType,
					"endpoint":      types.StringType,
					"tags":          types.ListType{ElemType: types.StringType},
					"consolidation": types.ObjectType{},
				}),
			},
			validate: func(t *testing.T, result armis.Action, diags diag.Diagnostics) {
				if result.Type != "log" {
					t.Errorf("Expected Type 'log', got '%s'", result.Type)
				}
				// Params should be empty when input params is null
				if result.Params.Severity != "" {
					t.Error("Expected empty Params when input is null")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, diags := convertActionModel(tt.input)
			tt.validate(t, result, diags)
		})
	}
}

// TestConvertListToActions tests the ConvertListToActions function.
func TestConvertListToActions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    types.List
		validate func(t *testing.T, result []armis.Action, diags diag.Diagnostics)
	}{
		{
			name: "null list returns nil",
			input: types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"type":   types.StringType,
					"params": types.ObjectType{},
				},
			}),
			validate: func(t *testing.T, result []armis.Action, diags diag.Diagnostics) {
				if result != nil {
					t.Error("Expected nil result for null list")
				}
				if diags.HasError() {
					t.Error("Expected no diagnostics errors")
				}
			},
		},
		{
			name: "unknown list returns nil",
			input: types.ListUnknown(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"type":   types.StringType,
					"params": types.ObjectType{},
				},
			}),
			validate: func(t *testing.T, result []armis.Action, diags diag.Diagnostics) {
				if result != nil {
					t.Error("Expected nil result for unknown list")
				}
				if diags.HasError() {
					t.Error("Expected no diagnostics errors")
				}
			},
		},
		{
			name:  "empty list returns empty slice",
			input: types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"type": types.StringType,
					"params": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"severity": types.StringType,
							"title":    types.StringType,
							"type":     types.StringType,
							"endpoint": types.StringType,
							"tags":     types.ListType{ElemType: types.StringType},
							"consolidation": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"amount": types.Int64Type,
									"unit":   types.StringType,
								},
							},
						},
					},
				},
			}, []attr.Value{}),
			validate: func(t *testing.T, result []armis.Action, diags diag.Diagnostics) {
				if len(result) != 0 {
					t.Errorf("Expected 0 actions, got %d", len(result))
				}
			},
		},
		{
			name: "list with single action",
			input: types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"type": types.StringType,
					"params": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"severity": types.StringType,
							"title":    types.StringType,
							"type":     types.StringType,
							"endpoint": types.StringType,
							"tags":     types.ListType{ElemType: types.StringType},
							"consolidation": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"amount": types.Int64Type,
									"unit":   types.StringType,
								},
							},
						},
					},
				},
			}, []attr.Value{
				types.ObjectValueMust(
					map[string]attr.Type{
						"type": types.StringType,
						"params": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"severity": types.StringType,
								"title":    types.StringType,
								"type":     types.StringType,
								"endpoint": types.StringType,
								"tags":     types.ListType{ElemType: types.StringType},
								"consolidation": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"amount": types.Int64Type,
										"unit":   types.StringType,
									},
								},
							},
						},
					},
					map[string]attr.Value{
						"type": types.StringValue("alert"),
						"params": types.ObjectValueMust(
							map[string]attr.Type{
								"severity": types.StringType,
								"title":    types.StringType,
								"type":     types.StringType,
								"endpoint": types.StringType,
								"tags":     types.ListType{ElemType: types.StringType},
								"consolidation": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"amount": types.Int64Type,
										"unit":   types.StringType,
									},
								},
							},
							map[string]attr.Value{
								"severity": types.StringValue("high"),
								"title":    types.StringValue("Test Alert"),
								"type":     types.StringNull(),
								"endpoint": types.StringNull(),
								"tags":     types.ListNull(types.StringType),
								"consolidation": types.ObjectNull(map[string]attr.Type{
									"amount": types.Int64Type,
									"unit":   types.StringType,
								}),
							},
						),
					},
				),
			}),
			validate: func(t *testing.T, result []armis.Action, diags diag.Diagnostics) {
				if len(result) != 1 {
					t.Fatalf("Expected 1 action, got %d", len(result))
				}
				if result[0].Type != "alert" {
					t.Errorf("Expected Type 'alert', got '%s'", result[0].Type)
				}
				if result[0].Params.Severity != "high" {
					t.Errorf("Expected Severity 'high', got '%s'", result[0].Params.Severity)
				}
				if result[0].Params.Title != "Test Alert" {
					t.Errorf("Expected Title 'Test Alert', got '%s'", result[0].Params.Title)
				}
				if diags.HasError() {
					t.Error("Expected no diagnostics errors")
				}
			},
		},
		{
			name: "list with multiple actions",
			input: types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"type": types.StringType,
					"params": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"severity": types.StringType,
							"title":    types.StringType,
							"type":     types.StringType,
							"endpoint": types.StringType,
							"tags":     types.ListType{ElemType: types.StringType},
							"consolidation": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"amount": types.Int64Type,
									"unit":   types.StringType,
								},
							},
						},
					},
				},
			}, []attr.Value{
				types.ObjectValueMust(
					map[string]attr.Type{
						"type": types.StringType,
						"params": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"severity": types.StringType,
								"title":    types.StringType,
								"type":     types.StringType,
								"endpoint": types.StringType,
								"tags":     types.ListType{ElemType: types.StringType},
								"consolidation": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"amount": types.Int64Type,
										"unit":   types.StringType,
									},
								},
							},
						},
					},
					map[string]attr.Value{
						"type": types.StringValue("alert"),
						"params": types.ObjectValueMust(
							map[string]attr.Type{
								"severity": types.StringType,
								"title":    types.StringType,
								"type":     types.StringType,
								"endpoint": types.StringType,
								"tags":     types.ListType{ElemType: types.StringType},
								"consolidation": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"amount": types.Int64Type,
										"unit":   types.StringType,
									},
								},
							},
							map[string]attr.Value{
								"severity": types.StringValue("critical"),
								"title":    types.StringValue("Security Alert"),
								"type":     types.StringNull(),
								"endpoint": types.StringNull(),
								"tags":     types.ListNull(types.StringType),
								"consolidation": types.ObjectNull(map[string]attr.Type{
									"amount": types.Int64Type,
									"unit":   types.StringType,
								}),
							},
						),
					},
				),
				types.ObjectValueMust(
					map[string]attr.Type{
						"type": types.StringType,
						"params": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"severity": types.StringType,
								"title":    types.StringType,
								"type":     types.StringType,
								"endpoint": types.StringType,
								"tags":     types.ListType{ElemType: types.StringType},
								"consolidation": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"amount": types.Int64Type,
										"unit":   types.StringType,
									},
								},
							},
						},
					},
					map[string]attr.Value{
						"type": types.StringValue("webhook"),
						"params": types.ObjectValueMust(
							map[string]attr.Type{
								"severity": types.StringType,
								"title":    types.StringType,
								"type":     types.StringType,
								"endpoint": types.StringType,
								"tags":     types.ListType{ElemType: types.StringType},
								"consolidation": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"amount": types.Int64Type,
										"unit":   types.StringType,
									},
								},
							},
							map[string]attr.Value{
								"severity": types.StringNull(),
								"title":    types.StringNull(),
								"type":     types.StringNull(),
								"endpoint": types.StringValue("https://example.com/webhook"),
								"tags":     types.ListNull(types.StringType),
								"consolidation": types.ObjectNull(map[string]attr.Type{
									"amount": types.Int64Type,
									"unit":   types.StringType,
								}),
							},
						),
					},
				),
			}),
			validate: func(t *testing.T, result []armis.Action, diags diag.Diagnostics) {
				if len(result) != 2 {
					t.Fatalf("Expected 2 actions, got %d", len(result))
				}
				// Validate first action
				if result[0].Type != "alert" {
					t.Errorf("Expected first action Type 'alert', got '%s'", result[0].Type)
				}
				if result[0].Params.Severity != "critical" {
					t.Errorf("Expected first action Severity 'critical', got '%s'", result[0].Params.Severity)
				}
				// Validate second action
				if result[1].Type != "webhook" {
					t.Errorf("Expected second action Type 'webhook', got '%s'", result[1].Type)
				}
				if result[1].Params.Endpoint != "https://example.com/webhook" {
					t.Errorf("Expected second action Endpoint 'https://example.com/webhook', got '%s'", result[1].Params.Endpoint)
				}
				if diags.HasError() {
					t.Error("Expected no diagnostics errors")
				}
			},
		},
		{
			name: "list with action containing full params",
			input: types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"type": types.StringType,
					"params": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"severity": types.StringType,
							"title":    types.StringType,
							"type":     types.StringType,
							"endpoint": types.StringType,
							"tags":     types.ListType{ElemType: types.StringType},
							"consolidation": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"amount": types.Int64Type,
									"unit":   types.StringType,
								},
							},
						},
					},
				},
			}, []attr.Value{
				types.ObjectValueMust(
					map[string]attr.Type{
						"type": types.StringType,
						"params": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"severity": types.StringType,
								"title":    types.StringType,
								"type":     types.StringType,
								"endpoint": types.StringType,
								"tags":     types.ListType{ElemType: types.StringType},
								"consolidation": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"amount": types.Int64Type,
										"unit":   types.StringType,
									},
								},
							},
						},
					},
					map[string]attr.Value{
						"type": types.StringValue("alert"),
						"params": types.ObjectValueMust(
							map[string]attr.Type{
								"severity": types.StringType,
								"title":    types.StringType,
								"type":     types.StringType,
								"endpoint": types.StringType,
								"tags":     types.ListType{ElemType: types.StringType},
								"consolidation": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"amount": types.Int64Type,
										"unit":   types.StringType,
									},
								},
							},
							map[string]attr.Value{
								"severity": types.StringValue("high"),
								"title":    types.StringValue("Comprehensive Alert"),
								"type":     types.StringValue("email"),
								"endpoint": types.StringValue("https://example.com/api/alerts"),
								"tags": types.ListValueMust(types.StringType, []attr.Value{
									types.StringValue("security"),
									types.StringValue("critical"),
									types.StringValue("network"),
								}),
								"consolidation": types.ObjectValueMust(
									map[string]attr.Type{
										"amount": types.Int64Type,
										"unit":   types.StringType,
									},
									map[string]attr.Value{
										"amount": types.Int64Value(15),
										"unit":   types.StringValue("minutes"),
									},
								),
							},
						),
					},
				),
			}),
			validate: func(t *testing.T, result []armis.Action, diags diag.Diagnostics) {
				if len(result) != 1 {
					t.Fatalf("Expected 1 action, got %d", len(result))
				}
				action := result[0]
				if action.Type != "alert" {
					t.Errorf("Expected Type 'alert', got '%s'", action.Type)
				}
				if action.Params.Severity != "high" {
					t.Errorf("Expected Severity 'high', got '%s'", action.Params.Severity)
				}
				if action.Params.Title != "Comprehensive Alert" {
					t.Errorf("Expected Title 'Comprehensive Alert', got '%s'", action.Params.Title)
				}
				if action.Params.Type != "email" {
					t.Errorf("Expected Type 'email', got '%s'", action.Params.Type)
				}
				if action.Params.Endpoint != "https://example.com/api/alerts" {
					t.Errorf("Expected Endpoint 'https://example.com/api/alerts', got '%s'", action.Params.Endpoint)
				}
				if len(action.Params.Tags) != 3 {
					t.Errorf("Expected 3 tags, got %d", len(action.Params.Tags))
				} else {
					expectedTags := []string{"security", "critical", "network"}
					for i, expected := range expectedTags {
						if action.Params.Tags[i] != expected {
							t.Errorf("Expected tag[%d] '%s', got '%s'", i, expected, action.Params.Tags[i])
						}
					}
				}
				if action.Params.Consolidation.Amount != 15 {
					t.Errorf("Expected Consolidation.Amount 15, got %d", action.Params.Consolidation.Amount)
				}
				if action.Params.Consolidation.Unit != "minutes" {
					t.Errorf("Expected Consolidation.Unit 'minutes', got '%s'", action.Params.Consolidation.Unit)
				}
				if diags.HasError() {
					t.Error("Expected no diagnostics errors")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, diags := ConvertListToActions(tt.input)
			tt.validate(t, result, diags)
		})
	}
}

// TestConvertModelToRules tests the ConvertModelToRules function.
func TestConvertModelToRules(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    RulesModel
		validate func(t *testing.T, result armis.Rules, diags diag.Diagnostics)
	}{
		{
			name: "null And and Or should return empty rules",
			input: RulesModel{
				And: types.ListNull(types.StringType),
				Or:  types.ListNull(types.StringType),
			},
			validate: func(t *testing.T, result armis.Rules, diags diag.Diagnostics) {
				if result.And != nil {
					t.Error("Expected And to be nil")
				}
				if result.Or != nil {
					t.Error("Expected Or to be nil")
				}
				if diags.HasError() {
					t.Error("Expected no diagnostics errors")
				}
			},
		},
		{
			name: "valid And rules",
			input: RulesModel{
				And: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("rule1"),
					types.StringValue("rule2"),
				}),
				Or: types.ListNull(types.StringType),
			},
			validate: func(t *testing.T, result armis.Rules, diags diag.Diagnostics) {
				if len(result.And) != 2 {
					t.Errorf("Expected And length 2, got %d", len(result.And))
				}
				if result.And[0] != "rule1" {
					t.Errorf("Expected And[0] = 'rule1', got '%v'", result.And[0])
				}
				if result.And[1] != "rule2" {
					t.Errorf("Expected And[1] = 'rule2', got '%v'", result.And[1])
				}
				if result.Or != nil {
					t.Error("Expected Or to be nil")
				}
			},
		},
		{
			name: "valid Or rules",
			input: RulesModel{
				And: types.ListNull(types.StringType),
				Or: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("ruleA"),
					types.StringValue("ruleB"),
					types.StringValue("ruleC"),
				}),
			},
			validate: func(t *testing.T, result armis.Rules, diags diag.Diagnostics) {
				if result.And != nil {
					t.Error("Expected And to be nil")
				}
				if len(result.Or) != 3 {
					t.Errorf("Expected Or length 3, got %d", len(result.Or))
				}
				if result.Or[0] != "ruleA" {
					t.Errorf("Expected Or[0] = 'ruleA', got '%v'", result.Or[0])
				}
			},
		},
		{
			name: "both And and Or rules",
			input: RulesModel{
				And: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("and1"),
					types.StringValue("and2"),
				}),
				Or: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("or1"),
				}),
			},
			validate: func(t *testing.T, result armis.Rules, diags diag.Diagnostics) {
				if len(result.And) != 2 {
					t.Errorf("Expected And length 2, got %d", len(result.And))
				}
				if len(result.Or) != 1 {
					t.Errorf("Expected Or length 1, got %d", len(result.Or))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, diags := ConvertModelToRules(tt.input)
			tt.validate(t, result, diags)
		})
	}
}
