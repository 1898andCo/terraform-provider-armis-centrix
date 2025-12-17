// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package sweep

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// SweepArmisRoles will delete all Armis roles with names starting with "tfacc".
func SweepArmisRoles(name string) *resource.Sweeper {
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

			ctx := context.Background()
			roles, err := client.GetRoles(ctx)
			if err != nil {
				return fmt.Errorf("error listing Armis roles: %w", err)
			}

			prefix := "tfacc"
			for _, role := range roles {
				if strings.HasPrefix(role.Name, prefix) {
					log.Printf("[INFO] Deleting Armis role: %s", role.Name)
					_, err := client.DeleteRole(ctx, strconv.Itoa(role.ID))
					if err != nil {
						log.Printf("[ERROR] Failed to delete Armis role %s: %s", role.Name, err)
					}
				}
			}

			return nil
		},
	}
}
