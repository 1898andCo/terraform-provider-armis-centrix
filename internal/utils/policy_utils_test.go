// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"
	"testing"

	"github.com/1898andCo/armis-sdk-go/armis"
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
		name     string
		input    types.Object
		validate func(t *testing.T, result armis.Consolidation, hasValue bool, diags diag.Diagnostics)
	}{
		{
			name: "null object returns false",
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
			name: "unknown object returns false",
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
			name: "empty list returns empty slice",
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

// TestConvertActionsToList tests the ConvertActionsToList function.
func TestConvertActionsToList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []armis.Action
		validate func(t *testing.T, result types.List)
	}{
		{
			name:  "nil actions should return null list",
			input: nil,
			validate: func(t *testing.T, result types.List) {
				if !result.IsNull() {
					t.Error("Expected null list for nil input")
				}
			},
		},
		{
			name:  "empty actions should return empty list",
			input: []armis.Action{},
			validate: func(t *testing.T, result types.List) {
				if result.IsNull() {
					t.Error("Expected non-null list for empty slice")
				}
				if len(result.Elements()) != 0 {
					t.Errorf("Expected 0 elements, got %d", len(result.Elements()))
				}
			},
		},
		{
			name: "single action with basic fields",
			input: []armis.Action{
				{
					Type: "alert",
					Params: armis.Params{
						Severity: "high",
						Title:    "Test Alert",
					},
				},
			},
			validate: func(t *testing.T, result types.List) {
				elements := result.Elements()
				if len(elements) != 1 {
					t.Errorf("Expected 1 element, got %d", len(elements))
				}
			},
		},
		{
			name: "action with all params fields populated",
			input: []armis.Action{
				{
					Type: "webhook",
					Params: armis.Params{
						Severity: "critical",
						Title:    "Security Alert",
						Type:     "notification",
						Endpoint: "https://example.com/webhook",
						Tags:     []string{"security", "urgent"},
						Consolidation: armis.Consolidation{
							Amount: 5,
							Unit:   "minutes",
						},
					},
				},
			},
			validate: func(t *testing.T, result types.List) {
				elements := result.Elements()
				if len(elements) != 1 {
					t.Errorf("Expected 1 element, got %d", len(elements))
				}
			},
		},
		{
			name: "action with empty type",
			input: []armis.Action{
				{
					Type: "",
					Params: armis.Params{
						Severity: "low",
					},
				},
			},
			validate: func(t *testing.T, result types.List) {
				elements := result.Elements()
				if len(elements) != 1 {
					t.Errorf("Expected 1 element, got %d", len(elements))
				}
			},
		},
		{
			name: "action with empty params",
			input: []armis.Action{
				{
					Type:   "email",
					Params: armis.Params{},
				},
			},
			validate: func(t *testing.T, result types.List) {
				elements := result.Elements()
				if len(elements) != 1 {
					t.Errorf("Expected 1 element, got %d", len(elements))
				}
			},
		},
		{
			name: "action with only consolidation",
			input: []armis.Action{
				{
					Type: "consolidate",
					Params: armis.Params{
						Consolidation: armis.Consolidation{
							Amount: 10,
							Unit:   "seconds",
						},
					},
				},
			},
			validate: func(t *testing.T, result types.List) {
				elements := result.Elements()
				if len(elements) != 1 {
					t.Errorf("Expected 1 element, got %d", len(elements))
				}
			},
		},
		{
			name: "action with nil tags",
			input: []armis.Action{
				{
					Type: "log",
					Params: armis.Params{
						Severity: "info",
						Tags:     nil,
					},
				},
			},
			validate: func(t *testing.T, result types.List) {
				elements := result.Elements()
				if len(elements) != 1 {
					t.Errorf("Expected 1 element, got %d", len(elements))
				}
			},
		},
		{
			name: "action with empty tags slice",
			input: []armis.Action{
				{
					Type: "notify",
					Params: armis.Params{
						Severity: "medium",
						Tags:     []string{},
					},
				},
			},
			validate: func(t *testing.T, result types.List) {
				elements := result.Elements()
				if len(elements) != 1 {
					t.Errorf("Expected 1 element, got %d", len(elements))
				}
			},
		},
		{
			name: "multiple actions with varied params",
			input: []armis.Action{
				{
					Type: "alert",
					Params: armis.Params{
						Severity: "high",
						Title:    "Alert 1",
					},
				},
				{
					Type: "email",
					Params: armis.Params{
						Type:     "notification",
						Endpoint: "admin@example.com",
					},
				},
				{
					Type: "webhook",
					Params: armis.Params{
						Endpoint: "https://api.example.com",
						Tags:     []string{"tag1", "tag2"},
					},
				},
			},
			validate: func(t *testing.T, result types.List) {
				elements := result.Elements()
				if len(elements) != 3 {
					t.Errorf("Expected 3 elements, got %d", len(elements))
				}
			},
		},
		{
			name: "action with zero consolidation amount",
			input: []armis.Action{
				{
					Type: "test",
					Params: armis.Params{
						Consolidation: armis.Consolidation{
							Amount: 0,
							Unit:   "hours",
						},
					},
				},
			},
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
			result := ConvertActionsToList(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestResponseToPolicyFromGet tests the ResponseToPolicyFromGet function.
func TestResponseToPolicyFromGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    armis.GetPolicySettings
		validate func(t *testing.T, result *PolicyResourceModel)
	}{
		{
			name: "minimal policy with required fields only",
			input: armis.GetPolicySettings{
				Name:        "Test Policy",
				Description: "Test Description",
				IsEnabled:   true,
				RuleType:    "aql",
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.Name.ValueString() != "Test Policy" {
					t.Errorf("Expected Name 'Test Policy', got '%s'", result.Name.ValueString())
				}
				if result.Description.ValueString() != "Test Description" {
					t.Errorf("Expected Description 'Test Description', got '%s'", result.Description.ValueString())
				}
				if !result.IsEnabled.ValueBool() {
					t.Error("Expected IsEnabled to be true")
				}
				if result.RuleType.ValueString() != "aql" {
					t.Errorf("Expected RuleType 'aql', got '%s'", result.RuleType.ValueString())
				}
			},
		},
		{
			name: "policy with labels",
			input: armis.GetPolicySettings{
				Name:        "Labeled Policy",
				Description: "Policy with labels",
				IsEnabled:   false,
				RuleType:    "managed",
				Labels:      []string{"security", "compliance", "critical"},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.Labels.IsNull() {
					t.Error("Expected non-null labels list")
				}
				elements := result.Labels.Elements()
				if len(elements) != 3 {
					t.Errorf("Expected 3 labels, got %d", len(elements))
				}
			},
		},
		{
			name: "policy with actions",
			input: armis.GetPolicySettings{
				Name:        "Action Policy",
				Description: "Policy with actions",
				IsEnabled:   true,
				RuleType:    "aql",
				Actions: []armis.Action{
					{
						Type: "alert",
						Params: armis.Params{
							Severity: "high",
							Title:    "Security Alert",
						},
					},
				},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.Actions.IsNull() {
					t.Error("Expected non-null actions list")
				}
				elements := result.Actions.Elements()
				if len(elements) != 1 {
					t.Errorf("Expected 1 action, got %d", len(elements))
				}
			},
		},
		{
			name: "policy with rules (AND)",
			input: armis.GetPolicySettings{
				Name:        "Rules Policy",
				Description: "Policy with AND rules",
				IsEnabled:   true,
				RuleType:    "aql",
				Rules: armis.Rules{
					And: []any{"rule1", "rule2", "rule3"},
				},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.Rules == nil {
					t.Fatal("Expected non-nil rules")
				}
				if result.Rules.And.IsNull() {
					t.Error("Expected non-null AND rules")
				}
				elements := result.Rules.And.Elements()
				if len(elements) != 3 {
					t.Errorf("Expected 3 AND rules, got %d", len(elements))
				}
			},
		},
		{
			name: "policy with rules (OR)",
			input: armis.GetPolicySettings{
				Name:        "OR Rules Policy",
				Description: "Policy with OR rules",
				IsEnabled:   false,
				RuleType:    "managed",
				Rules: armis.Rules{
					Or: []any{"ruleA", "ruleB"},
				},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.Rules == nil {
					t.Fatal("Expected non-nil rules")
				}
				if result.Rules.Or.IsNull() {
					t.Error("Expected non-null OR rules")
				}
				elements := result.Rules.Or.Elements()
				if len(elements) != 2 {
					t.Errorf("Expected 2 OR rules, got %d", len(elements))
				}
			},
		},
		{
			name: "policy with both AND and OR rules",
			input: armis.GetPolicySettings{
				Name:        "Complex Rules Policy",
				Description: "Policy with both rule types",
				IsEnabled:   true,
				RuleType:    "aql",
				Rules: armis.Rules{
					And: []any{"and1", "and2"},
					Or:  []any{"or1"},
				},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.Rules == nil {
					t.Fatal("Expected non-nil rules")
				}
				andElements := result.Rules.And.Elements()
				orElements := result.Rules.Or.Elements()
				if len(andElements) != 2 {
					t.Errorf("Expected 2 AND rules, got %d", len(andElements))
				}
				if len(orElements) != 1 {
					t.Errorf("Expected 1 OR rule, got %d", len(orElements))
				}
			},
		},
		{
			name: "fully populated policy",
			input: armis.GetPolicySettings{
				Name:        "Complete Policy",
				Description: "Fully populated policy",
				IsEnabled:   true,
				RuleType:    "aql",
				Labels:      []string{"prod", "critical"},
				MitreAttackLabels: []armis.MitreAttackLabel{
					{
						Matrix:       "enterprise",
						SubTechnique: "T1234.001",
						Tactic:       "Initial Access",
						Technique:    "Phishing",
					},
				},
				Actions: []armis.Action{
					{
						Type: "alert",
						Params: armis.Params{
							Severity: "critical",
							Title:    "Critical Alert",
							Tags:     []string{"urgent"},
							Consolidation: armis.Consolidation{
								Amount: 15,
								Unit:   "minutes",
							},
						},
					},
				},
				Rules: armis.Rules{
					And: []any{"rule1", "rule2"},
				},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.Name.ValueString() != "Complete Policy" {
					t.Errorf("Expected Name 'Complete Policy', got '%s'", result.Name.ValueString())
				}
				if result.Labels.IsNull() {
					t.Error("Expected non-null labels")
				}
				if result.Actions.IsNull() {
					t.Error("Expected non-null actions")
				}
				if result.Rules == nil {
					t.Fatal("Expected non-nil rules")
				}
			},
		},
		{
			name: "disabled policy",
			input: armis.GetPolicySettings{
				Name:        "Disabled Policy",
				Description: "This policy is disabled",
				IsEnabled:   false,
				RuleType:    "managed",
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.IsEnabled.ValueBool() {
					t.Error("Expected IsEnabled to be false")
				}
			},
		},
		{
			name: "policy with empty labels slice",
			input: armis.GetPolicySettings{
				Name:        "Empty Labels Policy",
				Description: "Policy with empty labels",
				IsEnabled:   true,
				RuleType:    "aql",
				Labels:      []string{},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				// Empty slice should create a non-null list with zero elements
				if !result.Labels.IsNull() {
					elements := result.Labels.Elements()
					if len(elements) != 0 {
						t.Errorf("Expected 0 labels, got %d", len(elements))
					}
				}
			},
		},
		{
			name: "policy with nil labels",
			input: armis.GetPolicySettings{
				Name:        "Nil Labels Policy",
				Description: "Policy with nil labels",
				IsEnabled:   true,
				RuleType:    "aql",
				Labels:      nil,
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if !result.Labels.IsNull() {
					t.Error("Expected null labels list for nil input")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ResponseToPolicyFromGet(context.Background(), tt.input)
			tt.validate(t, result)
		})
	}
}

// TestResponseToPolicyFromUpdate tests the ResponseToPolicyFromUpdate function.
func TestResponseToPolicyFromUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    armis.UpdatePolicySettings
		validate func(t *testing.T, result *PolicyResourceModel)
	}{
		{
			name: "minimal update policy",
			input: armis.UpdatePolicySettings{
				Name:        "Updated Policy",
				Description: "Updated Description",
				IsEnabled:   true,
				RuleType:    "aql",
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.Name.ValueString() != "Updated Policy" {
					t.Errorf("Expected Name 'Updated Policy', got '%s'", result.Name.ValueString())
				}
				if result.Description.ValueString() != "Updated Description" {
					t.Errorf("Expected Description 'Updated Description', got '%s'", result.Description.ValueString())
				}
			},
		},
		{
			name: "update policy with labels",
			input: armis.UpdatePolicySettings{
				Name:        "Labeled Update",
				Description: "Update with labels",
				IsEnabled:   false,
				RuleType:    "managed",
				Labels:      []string{"updated", "modified"},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.Labels.IsNull() {
					t.Error("Expected non-null labels list")
				}
				elements := result.Labels.Elements()
				if len(elements) != 2 {
					t.Errorf("Expected 2 labels, got %d", len(elements))
				}
			},
		},
		{
			name: "update policy with MITRE labels",
			input: armis.UpdatePolicySettings{
				Name:        "MITRE Update",
				Description: "Update with MITRE labels",
				IsEnabled:   true,
				RuleType:    "aql",
				MitreAttackLabels: []armis.MitreAttackLabel{
					{
						Matrix:       "enterprise",
						SubTechnique: "T1566.001",
						Tactic:       "Initial Access",
						Technique:    "Spearphishing Attachment",
					},
					{
						Matrix:       "enterprise",
						SubTechnique: "",
						Tactic:       "Execution",
						Technique:    "Command and Scripting",
					},
				},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				// MITRE labels are converted to strings in BuildPolicySettings, not in ResponseToPolicyFromUpdate
			},
		},
		{
			name: "update policy with actions",
			input: armis.UpdatePolicySettings{
				Name:        "Action Update",
				Description: "Update with actions",
				IsEnabled:   true,
				RuleType:    "aql",
				Actions: []armis.Action{
					{
						Type: "webhook",
						Params: armis.Params{
							Endpoint: "https://new-endpoint.com",
							Severity: "high",
						},
					},
				},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.Actions.IsNull() {
					t.Error("Expected non-null actions list")
				}
				elements := result.Actions.Elements()
				if len(elements) != 1 {
					t.Errorf("Expected 1 action, got %d", len(elements))
				}
			},
		},
		{
			name: "update policy with multiple actions",
			input: armis.UpdatePolicySettings{
				Name:        "Multi-Action Update",
				Description: "Update with multiple actions",
				IsEnabled:   true,
				RuleType:    "aql",
				Actions: []armis.Action{
					{
						Type: "alert",
						Params: armis.Params{
							Severity: "critical",
							Title:    "Alert 1",
						},
					},
					{
						Type: "email",
						Params: armis.Params{
							Type:     "notification",
							Endpoint: "security@example.com",
						},
					},
				},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.Actions.IsNull() {
					t.Error("Expected non-null actions list")
				}
				elements := result.Actions.Elements()
				if len(elements) != 2 {
					t.Errorf("Expected 2 actions, got %d", len(elements))
				}
			},
		},
		{
			name: "update policy with rules",
			input: armis.UpdatePolicySettings{
				Name:        "Rules Update",
				Description: "Update with rules",
				IsEnabled:   true,
				RuleType:    "aql",
				Rules: armis.Rules{
					And: []any{"updated_rule1", "updated_rule2"},
					Or:  []any{"or_rule1"},
				},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.Rules == nil {
					t.Fatal("Expected non-nil rules")
				}
				andElements := result.Rules.And.Elements()
				orElements := result.Rules.Or.Elements()
				if len(andElements) != 2 {
					t.Errorf("Expected 2 AND rules, got %d", len(andElements))
				}
				if len(orElements) != 1 {
					t.Errorf("Expected 1 OR rule, got %d", len(orElements))
				}
			},
		},
		{
			name: "fully populated update policy",
			input: armis.UpdatePolicySettings{
				Name:        "Complete Update",
				Description: "Fully updated policy",
				IsEnabled:   false,
				RuleType:    "managed",
				Labels:      []string{"updated", "v2"},
				MitreAttackLabels: []armis.MitreAttackLabel{
					{
						Matrix:    "enterprise",
						Tactic:    "Defense Evasion",
						Technique: "Obfuscation",
					},
				},
				Actions: []armis.Action{
					{
						Type: "log",
						Params: armis.Params{
							Severity: "info",
							Tags:     []string{"audit"},
						},
					},
				},
				Rules: armis.Rules{
					Or: []any{"rule_a", "rule_b", "rule_c"},
				},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.Name.ValueString() != "Complete Update" {
					t.Errorf("Expected Name 'Complete Update', got '%s'", result.Name.ValueString())
				}
				if result.IsEnabled.ValueBool() {
					t.Error("Expected IsEnabled to be false")
				}
			},
		},
		{
			name: "update policy toggling enabled state",
			input: armis.UpdatePolicySettings{
				Name:        "Toggle Policy",
				Description: "Policy with toggled state",
				IsEnabled:   false,
				RuleType:    "aql",
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if result.IsEnabled.ValueBool() {
					t.Error("Expected IsEnabled to be false")
				}
			},
		},
		{
			name: "update policy with nil labels",
			input: armis.UpdatePolicySettings{
				Name:        "Nil Labels Update",
				Description: "Update with nil labels",
				IsEnabled:   true,
				RuleType:    "aql",
				Labels:      nil,
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if !result.Labels.IsNull() {
					t.Error("Expected null labels list for nil input")
				}
			},
		},
		{
			name: "update policy with empty actions",
			input: armis.UpdatePolicySettings{
				Name:        "Empty Actions Update",
				Description: "Update with empty actions",
				IsEnabled:   true,
				RuleType:    "aql",
				Actions:     []armis.Action{},
			},
			validate: func(t *testing.T, result *PolicyResourceModel) {
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				// Empty slice should create a non-null list with zero elements
				if !result.Actions.IsNull() {
					elements := result.Actions.Elements()
					if len(elements) != 0 {
						t.Errorf("Expected 0 actions, got %d", len(elements))
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ResponseToPolicyFromUpdate(context.Background(), tt.input)
			tt.validate(t, result)
		})
	}
}

// TestBuildPolicySettings tests the BuildPolicySettings function.
func TestBuildPolicySettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *PolicyResourceModel
		validate func(t *testing.T, result armis.PolicySettings, diags diag.Diagnostics)
	}{
		{
			name: "minimal policy model",
			input: &PolicyResourceModel{
				Name:              types.StringValue("Test Policy"),
				Description:       types.StringValue("Test Description"),
				IsEnabled:         types.BoolValue(true),
				RuleType:          types.StringValue("aql"),
				Labels:            types.ListNull(types.StringType),
				MitreAttackLabels: types.ListNull(types.StringType),
				Actions:           types.ListNull(types.ObjectType{}),
				Rules: &RulesModel{
					And: types.ListNull(types.StringType),
					Or:  types.ListNull(types.StringType),
				},
			},
			validate: func(t *testing.T, result armis.PolicySettings, diags diag.Diagnostics) {
				if diags.HasError() {
					t.Errorf("Expected no errors, got %d", len(diags.Errors()))
				}
				if result.Name != "Test Policy" {
					t.Errorf("Expected Name 'Test Policy', got '%s'", result.Name)
				}
				if result.Description != "Test Description" {
					t.Errorf("Expected Description 'Test Description', got '%s'", result.Description)
				}
				if !result.IsEnabled {
					t.Error("Expected IsEnabled to be true")
				}
				if result.RuleType != "aql" {
					t.Errorf("Expected RuleType 'aql', got '%s'", result.RuleType)
				}
			},
		},
		{
			name: "policy model with labels",
			input: &PolicyResourceModel{
				Name:        types.StringValue("Labeled Policy"),
				Description: types.StringValue("Policy with labels"),
				IsEnabled:   types.BoolValue(false),
				RuleType:    types.StringValue("managed"),
				Labels: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("security"),
					types.StringValue("compliance"),
				}),
				MitreAttackLabels: types.ListNull(types.StringType),
				Actions:           types.ListNull(types.ObjectType{}),
				Rules: &RulesModel{
					And: types.ListNull(types.StringType),
					Or:  types.ListNull(types.StringType),
				},
			},
			validate: func(t *testing.T, result armis.PolicySettings, diags diag.Diagnostics) {
				if diags.HasError() {
					t.Errorf("Expected no errors, got %d", len(diags.Errors()))
				}
				if len(result.Labels) != 2 {
					t.Errorf("Expected 2 labels, got %d", len(result.Labels))
				}
				if result.Labels[0] != "security" {
					t.Errorf("Expected first label 'security', got '%s'", result.Labels[0])
				}
				if result.Labels[1] != "compliance" {
					t.Errorf("Expected second label 'compliance', got '%s'", result.Labels[1])
				}
			},
		},
		{
			name: "policy model with MITRE attack labels",
			input: &PolicyResourceModel{
				Name:        types.StringValue("MITRE Policy"),
				Description: types.StringValue("Policy with MITRE labels"),
				IsEnabled:   types.BoolValue(true),
				RuleType:    types.StringValue("aql"),
				Labels:      types.ListNull(types.StringType),
				MitreAttackLabels: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("T1566"),
					types.StringValue("T1059"),
				}),
				Actions: types.ListNull(types.ObjectType{}),
				Rules: &RulesModel{
					And: types.ListNull(types.StringType),
					Or:  types.ListNull(types.StringType),
				},
			},
			validate: func(t *testing.T, result armis.PolicySettings, diags diag.Diagnostics) {
				if diags.HasError() {
					t.Errorf("Expected no errors, got %d", len(diags.Errors()))
				}
				if len(result.MitreAttackLabels) != 2 {
					t.Errorf("Expected 2 MITRE labels, got %d", len(result.MitreAttackLabels))
				}
			},
		},
		{
			name: "policy model with AND rules",
			input: &PolicyResourceModel{
				Name:              types.StringValue("AND Rules Policy"),
				Description:       types.StringValue("Policy with AND rules"),
				IsEnabled:         types.BoolValue(true),
				RuleType:          types.StringValue("aql"),
				Labels:            types.ListNull(types.StringType),
				MitreAttackLabels: types.ListNull(types.StringType),
				Actions:           types.ListNull(types.ObjectType{}),
				Rules: &RulesModel{
					And: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("rule1"),
						types.StringValue("rule2"),
						types.StringValue("rule3"),
					}),
					Or: types.ListNull(types.StringType),
				},
			},
			validate: func(t *testing.T, result armis.PolicySettings, diags diag.Diagnostics) {
				if diags.HasError() {
					t.Errorf("Expected no errors, got %d", len(diags.Errors()))
				}
				if len(result.Rules.And) != 3 {
					t.Errorf("Expected 3 AND rules, got %d", len(result.Rules.And))
				}
				if result.Rules.And[0] != "rule1" {
					t.Errorf("Expected first AND rule 'rule1', got '%v'", result.Rules.And[0])
				}
			},
		},
		{
			name: "policy model with OR rules",
			input: &PolicyResourceModel{
				Name:              types.StringValue("OR Rules Policy"),
				Description:       types.StringValue("Policy with OR rules"),
				IsEnabled:         types.BoolValue(false),
				RuleType:          types.StringValue("managed"),
				Labels:            types.ListNull(types.StringType),
				MitreAttackLabels: types.ListNull(types.StringType),
				Actions:           types.ListNull(types.ObjectType{}),
				Rules: &RulesModel{
					And: types.ListNull(types.StringType),
					Or: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("orRule1"),
						types.StringValue("orRule2"),
					}),
				},
			},
			validate: func(t *testing.T, result armis.PolicySettings, diags diag.Diagnostics) {
				if diags.HasError() {
					t.Errorf("Expected no errors, got %d", len(diags.Errors()))
				}
				if len(result.Rules.Or) != 2 {
					t.Errorf("Expected 2 OR rules, got %d", len(result.Rules.Or))
				}
			},
		},
		{
			name: "policy model with both AND and OR rules",
			input: &PolicyResourceModel{
				Name:              types.StringValue("Complex Rules"),
				Description:       types.StringValue("Policy with AND and OR rules"),
				IsEnabled:         types.BoolValue(true),
				RuleType:          types.StringValue("aql"),
				Labels:            types.ListNull(types.StringType),
				MitreAttackLabels: types.ListNull(types.StringType),
				Actions:           types.ListNull(types.ObjectType{}),
				Rules: &RulesModel{
					And: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("and1"),
					}),
					Or: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("or1"),
						types.StringValue("or2"),
					}),
				},
			},
			validate: func(t *testing.T, result armis.PolicySettings, diags diag.Diagnostics) {
				if diags.HasError() {
					t.Errorf("Expected no errors, got %d", len(diags.Errors()))
				}
				if len(result.Rules.And) != 1 {
					t.Errorf("Expected 1 AND rule, got %d", len(result.Rules.And))
				}
				if len(result.Rules.Or) != 2 {
					t.Errorf("Expected 2 OR rules, got %d", len(result.Rules.Or))
				}
			},
		},
		{
			name: "disabled policy model",
			input: &PolicyResourceModel{
				Name:              types.StringValue("Disabled"),
				Description:       types.StringValue("Disabled policy"),
				IsEnabled:         types.BoolValue(false),
				RuleType:          types.StringValue("aql"),
				Labels:            types.ListNull(types.StringType),
				MitreAttackLabels: types.ListNull(types.StringType),
				Actions:           types.ListNull(types.ObjectType{}),
				Rules: &RulesModel{
					And: types.ListNull(types.StringType),
					Or:  types.ListNull(types.StringType),
				},
			},
			validate: func(t *testing.T, result armis.PolicySettings, diags diag.Diagnostics) {
				if diags.HasError() {
					t.Errorf("Expected no errors, got %d", len(diags.Errors()))
				}
				if result.IsEnabled {
					t.Error("Expected IsEnabled to be false")
				}
			},
		},
		{
			name: "policy model with empty labels list",
			input: &PolicyResourceModel{
				Name:              types.StringValue("Empty Labels"),
				Description:       types.StringValue("Policy with empty labels"),
				IsEnabled:         types.BoolValue(true),
				RuleType:          types.StringValue("aql"),
				Labels:            types.ListValueMust(types.StringType, []attr.Value{}),
				MitreAttackLabels: types.ListNull(types.StringType),
				Actions:           types.ListNull(types.ObjectType{}),
				Rules: &RulesModel{
					And: types.ListNull(types.StringType),
					Or:  types.ListNull(types.StringType),
				},
			},
			validate: func(t *testing.T, result armis.PolicySettings, diags diag.Diagnostics) {
				if diags.HasError() {
					t.Errorf("Expected no errors, got %d", len(diags.Errors()))
				}
				if result.Labels == nil {
					t.Error("Expected non-nil labels slice")
				}
				if len(result.Labels) != 0 {
					t.Errorf("Expected 0 labels, got %d", len(result.Labels))
				}
			},
		},
		{
			name: "policy model with null labels",
			input: &PolicyResourceModel{
				Name:              types.StringValue("Null Labels"),
				Description:       types.StringValue("Policy with null labels"),
				IsEnabled:         types.BoolValue(true),
				RuleType:          types.StringValue("aql"),
				Labels:            types.ListNull(types.StringType),
				MitreAttackLabels: types.ListNull(types.StringType),
				Actions:           types.ListNull(types.ObjectType{}),
				Rules: &RulesModel{
					And: types.ListNull(types.StringType),
					Or:  types.ListNull(types.StringType),
				},
			},
			validate: func(t *testing.T, result armis.PolicySettings, diags diag.Diagnostics) {
				if diags.HasError() {
					t.Errorf("Expected no errors, got %d", len(diags.Errors()))
				}
				if result.Labels != nil {
					t.Error("Expected nil labels for null input")
				}
			},
		},
		{
			name: "policy model with multiple labels and MITRE labels",
			input: &PolicyResourceModel{
				Name:        types.StringValue("Multi-Label Policy"),
				Description: types.StringValue("Policy with multiple label types"),
				IsEnabled:   types.BoolValue(true),
				RuleType:    types.StringValue("aql"),
				Labels: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("label1"),
					types.StringValue("label2"),
					types.StringValue("label3"),
				}),
				MitreAttackLabels: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("T1001"),
					types.StringValue("T1002"),
				}),
				Actions: types.ListNull(types.ObjectType{}),
				Rules: &RulesModel{
					And: types.ListNull(types.StringType),
					Or:  types.ListNull(types.StringType),
				},
			},
			validate: func(t *testing.T, result armis.PolicySettings, diags diag.Diagnostics) {
				if diags.HasError() {
					t.Errorf("Expected no errors, got %d", len(diags.Errors()))
				}
				if len(result.Labels) != 3 {
					t.Errorf("Expected 3 labels, got %d", len(result.Labels))
				}
				if len(result.MitreAttackLabels) != 2 {
					t.Errorf("Expected 2 MITRE labels, got %d", len(result.MitreAttackLabels))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, diags := BuildPolicySettings(tt.input)
			tt.validate(t, result, diags)
		})
	}
}
