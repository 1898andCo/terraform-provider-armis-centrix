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

// SweepArmisUsers will delete all Armis users with users starting with "test".
func SweepArmisUsers(name string) *resource.Sweeper {
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

			// Get all users
			ctx := context.Background()
			users, err := client.GetUsers(ctx)
			if err != nil {
				return fmt.Errorf("error getting Armis users: %w", err)
			}

			prefix := "test."
			for _, user := range users {
				if strings.HasPrefix(user.Username, prefix) {
					log.Printf("[INFO] Deleting Armis user: %s", user.Username)
					_, err := client.DeleteUser(ctx, strconv.Itoa(user.ID))
					if err != nil {
						log.Printf("[ERROR] Failed to delete Armis user %s: %s", user.Username, err)
					}
				}
			}

			return nil
		},
	}
}
