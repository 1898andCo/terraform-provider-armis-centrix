// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"testing"
)

func TestCreatingPolicy(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	payload := PolicySettings{
		Name:        "Test Policy",
		Description: "This is a test policy",
		IsEnabled:   false,
		Labels:      []string{"Security"},
		MitreAttackLabels: []string{
			"Enterprise.TA0009.T1056.001",
			"Enterprise.TA0009.T1056.004",
		},
		RuleType: "ACTIVITY",
		Actions: []Action{
			{
				Type: "alert",
				Params: Params{
					Severity: "high",
					Title:    "Test Security Alert",
					Type:     "Security - Threat",
					Consolidation: Consolidation{
						Amount: 1,
						Unit:   "Days",
					},
				},
			},
		},
		Rules: Rules{
			And: []any{
				"protocol:BMS",
				Rules{Or: []any{"content:(iPhone)", "content:(Android)"}},
			},
		},
	}

	res, err := client.CreatePolicy(context.Background(), payload)
	if err != nil {
		t.Fatalf("create policy: %v", err)
	}
	prettyPrint(res)
}

func TestCreatingTagPolicy(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	payload := PolicySettings{
		Name:        "Test Tag Policy",
		Description: "This is a test tag policy",
		IsEnabled:   false,
		Labels:      []string{"Security"},
		MitreAttackLabels: []string{
			"Enterprise.TA0009.T1056.001",
			"Enterprise.TA0009.T1056.004",
		},
		RuleType: "ACTIVITY",
		Actions: []Action{
			{
				Type: "tag",
				Params: Params{
					Endpoint: "ALL",
					Tags:     []string{"Agent and Scanner Gaps"},
				},
			},
		},
		Rules: Rules{
			And: []any{
				"protocol:BMS",
				Rules{Or: []any{"content:(iPhone)", "content:(Android)"}},
			},
		},
	}

	res, err := client.CreatePolicy(context.Background(), payload)
	if err != nil {
		t.Fatalf("create tag policy: %v", err)
	}
	prettyPrint(res)
}

func TestGettingAllPolicies(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	res, err := client.GetAllPolicies(context.Background())
	if err != nil {
		t.Fatalf("get policies: %v", err)
	}
	prettyPrint(res)
}

func TestGettingPolicy(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	const id = "76884"
	res, err := client.GetPolicy(context.Background(), id)
	if err != nil {
		t.Fatalf("get policy: %v", err)
	}
	prettyPrint(res)
}

func TestUpdatingPolicy(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	payload := PolicySettings{
		Name:        "Test Policy Updated",
		Description: "This is an updated test policy",
		IsEnabled:   true,
		Labels:      []string{"Security"},
		MitreAttackLabels: []string{
			"Enterprise.TA0009.T1056.001",
			"Enterprise.TA0009.T1056.004",
		},
		RuleType: "ACTIVITY",
		Actions: []Action{
			{
				Type: "alert",
				Params: Params{
					Severity: "high",
					Title:    "Test Security Alert",
					Type:     "Security - Threat",
					Consolidation: Consolidation{
						Amount: 1,
						Unit:   "Days",
					},
				},
			},
		},
		Rules: Rules{
			And: []any{
				"protocol:BMS",
				Rules{Or: []any{"content:(iPhone)", "content:(Android)"}},
			},
		},
	}

	const id = "76700"
	res, err := client.UpdatePolicy(context.Background(), payload, id)
	if err != nil {
		t.Fatalf("update policy: %v", err)
	}
	prettyPrint(res)
}

func TestDeletingPolicy(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	const id = "76700"
	ok, err := client.DeletePolicy(context.Background(), id)
	if err != nil {
		t.Fatalf("delete policy: %v", err)
	}
	if !ok {
		t.Fatalf("policy %s not deleted", id)
	}
}
