// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package sweep

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// SweepArmisPolicies returns a sweeper function to clean up Armis policies.
func SweepArmisPolicies(name string) *resource.Sweeper {
	return &resource.Sweeper{
		Name: name,
		F: func(_ string) error {
			client, err := ConfigureSweeperClient(name)
			if err != nil {
				return fmt.Errorf("error configuring Armis client: %w", err)
			}
			if client == nil {
				return nil
			}

			// Get all policies
			ctx := context.Background()
			policies, err := client.GetAllPolicies(ctx)
			if err != nil {
				return fmt.Errorf("error getting Armis policies: %w", err)
			}

			prefix := "tfacc"
			for _, policy := range policies {
				if strings.HasPrefix(policy.Name, prefix) {
					log.Printf("[INFO] Deleting Armis policy: %s", policy.Name)
					_, err := client.DeletePolicy(ctx, policy.ID)
					if err != nil {
						log.Printf("[ERROR] Failed to delete Armis policy %s: %s", policy.Name, err)
					}
				}
			}

			return nil
		},
	}
}
